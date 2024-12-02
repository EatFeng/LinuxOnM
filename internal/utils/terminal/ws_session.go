package terminal

import (
	"LinuxOnM/internal/global"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"io"
	"sync"
	"time"
)

const (
	WsMsgCmd       = "cmd"
	WsMsgResize    = "resize"
	WsMsgHeartbeat = "heartbeat"
)

type safeBuffer struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

func (w *safeBuffer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}

func (w *safeBuffer) Bytes() []byte {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Bytes()
}

type LogicSshWsSession struct {
	stdinPipe       io.WriteCloser
	comboOutput     *safeBuffer
	logBuff         *safeBuffer
	inputFilterBuff *safeBuffer
	session         *ssh.Session
	wsConn          *websocket.Conn
	isAdmin         bool
	IsFlagged       bool
}

type WsMsg struct {
	Type      string `json:"type"`
	Data      string `json:"data,omitempty"`      // WsMsgCmd
	Cols      int    `json:"cols,omitempty"`      // WsMsgResize
	Rows      int    `json:"rows,omitempty"`      // WsMsgResize
	Timestamp int    `json:"timestamp,omitempty"` // WsMsgHeartbeat
}

func NewLogicSshWsSession(cols, rows int, isAdmin bool, sshClient *ssh.Client, wsConn *websocket.Conn) (*LogicSshWsSession, error) {
	sshSession, err := sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	stdinP, err := sshSession.StdinPipe()
	if err != nil {
		return nil, err
	}

	comboWriter := new(safeBuffer)
	logBuf := new(safeBuffer)
	inputBuf := new(safeBuffer)
	sshSession.Stdout = comboWriter
	sshSession.Stderr = comboWriter

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := sshSession.RequestPty("xterm", rows, cols, modes); err != nil {
		return nil, err
	}
	if err := sshSession.Shell(); err != nil {
		return nil, err
	}
	return &LogicSshWsSession{
		stdinPipe:       stdinP,
		comboOutput:     comboWriter,
		logBuff:         logBuf,
		inputFilterBuff: inputBuf,
		session:         sshSession,
		wsConn:          wsConn,
		isAdmin:         isAdmin,
		IsFlagged:       false,
	}, nil
}

func (sws *LogicSshWsSession) Close() {
	if sws.session != nil {
		sws.session.Close()
	}
	if sws.logBuff != nil {
		sws.logBuff = nil
	}
	if sws.comboOutput != nil {
		sws.comboOutput = nil
	}
}

func (sws *LogicSshWsSession) Start(quitChan chan bool) {
	go sws.receiveWsMsg(quitChan)
	go sws.sendComboOutput(quitChan)
}

// receiveWsMsg is a method of the LogicSshWsSession structure, primarily used for handling the logic of receiving messages from a WebSocket connection.
// It operates in a separate coroutine, continuously monitoring messages on the WebSocket connection, and performs
// corresponding processing operations based on the different types of messages.
// It also handles potential exceptions and exits the coroutine under specific conditions.
func (sws *LogicSshWsSession) receiveWsMsg(exitCh chan bool) {
	defer func() {
		if r := recover(); r != nil {
			global.LOG.Errorf("[A panic occurred during receive ws message, error message: %v", r)
		}
	}()
	wsConn := sws.wsConn
	defer setQuit(exitCh)
	for {
		select {
		case <-exitCh:
			return
		default:
			_, wsData, err := wsConn.ReadMessage()
			if err != nil {
				return
			}
			msgObj := WsMsg{}
			_ = json.Unmarshal(wsData, &msgObj)
			switch msgObj.Type {
			case WsMsgResize:
				if msgObj.Cols > 0 && msgObj.Rows > 0 {
					if err := sws.session.WindowChange(msgObj.Rows, msgObj.Cols); err != nil {
						global.LOG.Errorf("ssh pty change windows size failed, err: %v", err)
					}
				}
			case WsMsgCmd:
				decodeBytes, err := base64.StdEncoding.DecodeString(msgObj.Data)
				if err != nil {
					global.LOG.Errorf("websock cmd string base64 decoding failed, err: %v", err)
				}
				sws.sendWebsocketInputCommandToSshSessionStdinPipe(decodeBytes)
			case WsMsgHeartbeat:
				err = wsConn.WriteMessage(websocket.TextMessage, wsData)
				if err != nil {
					global.LOG.Errorf("ssh sending heartbeat to webSocket failed, err: %v", err)
				}
			}
		}
	}
}

func (sws *LogicSshWsSession) sendWebsocketInputCommandToSshSessionStdinPipe(cmdBytes []byte) {
	if _, err := sws.stdinPipe.Write(cmdBytes); err != nil {
		global.LOG.Errorf("ws cmd bytes write to ssh.stdin pipe failed, err: %v", err)
	}
}

// sendComboOutput is a method of the LogicSshWsSession structure, mainly responsible for regularly sending
// the standard output and standard error output of the SSH session (merged in the comboOutput buffer) to
// the WebSocket connection, so that the client can receive the execution results and other information of
// the SSH session. At the same time, it also handles potential error situations and ends its own execution
// based on specific conditions or external exit signals.
func (sws *LogicSshWsSession) sendComboOutput(exitCh chan bool) {
	wsConn := sws.wsConn
	defer setQuit(exitCh)

	tick := time.NewTicker(time.Millisecond * time.Duration(60))
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if sws.comboOutput == nil {
				return
			}
			bs := sws.comboOutput.Bytes()
			if len(bs) > 0 {
				wsData, err := json.Marshal(WsMsg{
					Type: WsMsgCmd,
					Data: base64.StdEncoding.EncodeToString(bs),
				})
				if err != nil {
					global.LOG.Errorf("encoding combo output to json failed, err: %v", err)
					continue
				}
				err = wsConn.WriteMessage(websocket.TextMessage, wsData)
				if err != nil {
					global.LOG.Errorf("ssh sending combo output to webSocket failed, err: %v", err)
				}
				_, err = sws.logBuff.Write(bs)
				if err != nil {
					global.LOG.Errorf("combo output to log buffer failed, err: %v", err)
				}
				sws.comboOutput.buffer.Reset()
			}
			if string(bs) == string([]byte{13, 10, 108, 111, 103, 111, 117, 116, 13, 10}) {
				sws.Close()
				return
			}

		case <-exitCh:
			return
		}
	}
}

func (sws *LogicSshWsSession) Wait(quitChan chan bool) {
	if err := sws.session.Wait(); err != nil {
		setQuit(quitChan)
	}
}

func setQuit(ch chan bool) {
	ch <- true
}
