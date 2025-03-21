package handlers

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/utils/cmd"
	"LinuxOnM/internal/utils/copier"
	"LinuxOnM/internal/utils/ssh"
	"LinuxOnM/internal/utils/terminal"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (b *BaseApi) WsSsh(c *gin.Context) {
	// Upgrade the HTTP connection to a WebSocket connection
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.LOG.Errorf("gin context http handler failed, err: %v", err)
		return
	}
	defer wsConn.Close()

	// Retrieve and validate the 'id' parameter from the request query
	id, err := strconv.Atoi(c.Query("id"))
	if wshandleError(wsConn, errors.WithMessage(err, "invalid param id in request")) {
		return
	}

	// Retrieve and validate the 'cols' parameter from the request query, with a default value of 80 if not provided
	cols, err := strconv.Atoi(c.DefaultQuery("cols", "80"))
	if wshandleError(wsConn, errors.WithMessage(err, "invalid param cols in request")) {
		return
	}

	// Retrieve and validate the 'rows' parameter from the request query, with a default value of 40 if not provided
	rows, err := strconv.Atoi(c.DefaultQuery("rows", "40"))
	if wshandleError(wsConn, errors.WithMessage(err, "invalid param rows in request")) {
		return
	}

	// Retrieve host information based on the provided 'id'
	host, err := hostService.GetHostInfo(uint(id))
	if wshandleError(wsConn, errors.WithMessage(err, "load host info by id failed")) {
		return
	}

	// Set up the SSH connection information based on the retrieved host information
	var connInfo ssh.ConnInfo
	_ = copier.Copy(&connInfo, &host)
	connInfo.PrivateKey = []byte(host.PrivateKey)
	if len(host.PassPhrase) != 0 {
		connInfo.PassPhrase = []byte(host.PassPhrase)
	}

	// Establish the SSH client connection using the connection information
	client, err := connInfo.NewClient()
	if wshandleError(wsConn, errors.WithMessage(err, "failed to set up the connection. Please check the host information")) {
		return
	}
	defer client.Close()

	// Create an SSH WebSocket session with the specified parameters and the established SSH client and WebSocket connection
	sws, err := terminal.NewLogicSshWsSession(cols, rows, true, connInfo.Client, wsConn)
	if wshandleError(wsConn, err) {
		return
	}
	defer sws.Close()

	// Create a channel to signal the end of the SSH WebSocket session
	quitChan := make(chan bool, 3)

	// Start the SSH WebSocket session
	sws.Start(quitChan)

	// Wait for the SSH WebSocket session to complete in a separate goroutine
	go sws.Wait(quitChan)

	// Wait for the signal from the quitChan, indicating the end of the session
	<-quitChan

	// Check for any errors after the session has completed
	if wshandleError(wsConn, err) {
		return
	}
}

func wshandleError(ws *websocket.Conn, err error) bool {
	if err != nil {
		global.LOG.Errorf("handler ws faled:, err: %v", err)
		dt := time.Now().Add(time.Second)
		if ctlerr := ws.WriteControl(websocket.CloseMessage, []byte(err.Error()), dt); ctlerr != nil {
			wsData, err := json.Marshal(terminal.WsMsg{
				Type: terminal.WsMsgCmd,
				Data: base64.StdEncoding.EncodeToString([]byte(err.Error())),
			})
			if err != nil {
				_ = ws.WriteMessage(websocket.TextMessage, []byte("{\"type\":\"cmd\",\"data\":\"failed to encoding to json\"}"))
			} else {
				_ = ws.WriteMessage(websocket.TextMessage, wsData)
			}
		}
		return true
	}
	return false
}

func (b *BaseApi) ContainerWsSsh(c *gin.Context) {
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.LOG.Errorf("gin context http handler failed, err: %v", err)
		return
	}
	defer wsConn.Close()

	containerID := c.Query("containerid")
	command := c.Query("command")
	user := c.Query("user")
	if len(command) == 0 || len(containerID) == 0 {
		if wshandleError(wsConn, errors.New("error param of command or containerID")) {
			return
		}
	}
	cols, err := strconv.Atoi(c.DefaultQuery("cols", "80"))
	if wshandleError(wsConn, errors.WithMessage(err, "invalid param cols in request")) {
		return
	}
	rows, err := strconv.Atoi(c.DefaultQuery("rows", "40"))
	if wshandleError(wsConn, errors.WithMessage(err, "invalid param rows in request")) {
		return
	}

	cmds := []string{"exec", containerID, command}
	if len(user) != 0 {
		cmds = []string{"exec", "-u", user, containerID, command}
	}
	if cmd.CheckIllegal(user, containerID, command) {
		if wshandleError(wsConn, errors.New("  The command contains illegal characters.")) {
			return
		}
	}
	stdout, err := cmd.ExecWithCheck("docker", cmds...)
	if wshandleError(wsConn, errors.WithMessage(err, stdout)) {
		return
	}

	commands := []string{"exec", "-it", containerID, command}
	if len(user) != 0 {
		commands = []string{"exec", "-it", "-u", user, containerID, command}
	}
	pidMap := loadMapFromDockerTop(containerID)
	slave, err := terminal.NewCommand(commands)
	if wshandleError(wsConn, err) {
		return
	}
	defer killBash(containerID, command, pidMap)
	defer slave.Close()

	tty, err := terminal.NewLocalWsSession(cols, rows, wsConn, slave, true)
	if wshandleError(wsConn, err) {
		return
	}

	quitChan := make(chan bool, 3)
	tty.Start(quitChan)
	go slave.Wait(quitChan)

	<-quitChan

	global.LOG.Info("websocket finished")
	if wshandleError(wsConn, err) {
		return
	}
}

func loadMapFromDockerTop(containerID string) map[string]string {
	pidMap := make(map[string]string)
	sudo := cmd.SudoHandleCmd()

	stdout, err := cmd.Execf("%s docker top %s -eo pid,command ", sudo, containerID)
	if err != nil {
		return pidMap
	}
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		pidMap[parts[0]] = strings.Join(parts[1:], " ")
	}
	return pidMap
}

func killBash(containerID, comm string, pidMap map[string]string) {
	sudo := cmd.SudoHandleCmd()
	newPidMap := loadMapFromDockerTop(containerID)
	for pid, command := range newPidMap {
		isOld := false
		for pid2 := range pidMap {
			if pid == pid2 {
				isOld = true
				break
			}
		}
		if !isOld && command == comm {
			_, _ = cmd.Execf("%s kill -9 %s", sudo, pid)
		}
	}
}
