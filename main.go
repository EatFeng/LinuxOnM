// @title LinuxOnM API
// @version 1.0.0
// @description This is the API documentation for LinuxOnM.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@yourwebapp.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8848
// @BasePath /api/handler

package main

import (
	"LinuxOnM/internal/cron"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/init/app"
	"LinuxOnM/internal/init/cache"
	"LinuxOnM/internal/init/db"
	"LinuxOnM/internal/init/hook"
	"LinuxOnM/internal/init/log"
	"LinuxOnM/internal/init/migration"
	"LinuxOnM/internal/init/router"
	"LinuxOnM/internal/init/session"
	"LinuxOnM/internal/init/session/psession"
	"LinuxOnM/internal/init/validator"
	"LinuxOnM/internal/init/viper"
	"crypto/tls"
	"encoding/gob"
	"net"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

func main() {
	viper.Init()
	log.Init()
	db.Init()
	migration.Init()
	app.Init()
	validator.Init()
	gob.Register(psession.SessionUser{})
	cache.Init()
	session.Init()
	gin.SetMode("debug")
	cron.Run()
	hook.Init()

	rootRouter := router.Routers()

	tcpItem := "tcp4"

	server := &http.Server{
		Addr:    global.CONF.System.BindAddress + ":" + global.CONF.System.Port,
		Handler: rootRouter,
	}

	ln, err := net.Listen(tcpItem, server.Addr)
	if err != nil {
		panic(err)
	}

	type tcpKeepAliveListener struct {
		*net.TCPListener
	}

	if global.CONF.System.SSL == "enable" {
		certPath := path.Join(global.CONF.System.BaseDir, "LinuxOnM/secret/server.crt")
		keyPath := path.Join(global.CONF.System.BaseDir, "LinuxOnM/secret/server.key")
		certificate, err := os.ReadFile(certPath)
		if err != nil {
			panic(err)
		}
		key, err := os.ReadFile(keyPath)
		if err != nil {
			panic(err)
		}
		cert, err := tls.X509KeyPair(certificate, key)
		if err != nil {
			panic(err)
		}
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		global.LOG.Infof("listen at https://%s:%s [%s]", global.CONF.System.BindAddress, global.CONF.System.Port, tcpItem)

		if err := server.ServeTLS(tcpKeepAliveListener{ln.(*net.TCPListener)}, certPath, keyPath); err != nil {
			panic(err)
		}
	} else {
		global.LOG.Infof("listen at http://%s:%s [%s]", global.CONF.System.BindAddress, global.CONF.System.Port, tcpItem)
		if err := server.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)}); err != nil {
			panic(err)
		}
	}
}
