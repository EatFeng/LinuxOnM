package handlers

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/utils/copier"
	"LinuxOnM/internal/utils/ssh"
	"LinuxOnM/internal/utils/terminal"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsSsh is a route handler function that establishes an SSH connection over WebSocket.
// It upgrades the HTTP connection to a WebSocket connection, validates and retrieves necessary parameters from the request,
// sets up the SSH connection details, creates an SSH WebSocket session, starts the session, and waits for it to complete.
//
// @Tags SSH, WebSocket
// @Summary Establish SSH connection over WebSocket
// @Description This endpoint upgrades the incoming HTTP connection to a WebSocket connection and then
// establishes an SSH connection to a specified host using the provided parameters. It allows interaction with the
// remote host's shell via the WebSocket connection.
// @Accept websocket
// @Produce json
// @Param id query int true "The ID of the host to connect to"
// @Param cols query int false "The number of columns for the terminal session (default: 80)"
// @Param rows query int false "The number of rows for the terminal session (default: 40)"
// @Success 200 {string} string "WebSocket connection established successfully and SSH session started."
// @Failure 400 {string} string "Invalid parameters in the request"
// @Failure 500 {string} string "Failed to establish SSH connection or other internal errors occurred"
// @Router /ws-ssh [get]
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
