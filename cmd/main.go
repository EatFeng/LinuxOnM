package main

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/init/db"
	"LinuxOnM/internal/init/log"
	"LinuxOnM/internal/init/router"
	"LinuxOnM/internal/init/viper"
	"net"
	"net/http"
)

func main() {
	viper.Init()
	log.Init()
	db.Init()

	rootRouter := router.Routers()

	tcpItem := "tcp4"

	global.CONF.System.BindAddress = "0.0.0.0"
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

	global.LOG.Infof("listen at http://%s:%s [%s]", global.CONF.System.BindAddress, global.CONF.System.Port, tcpItem)
	if err := server.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)}); err != nil {
		panic(err)
	}
}
