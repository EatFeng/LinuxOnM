package session

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/init/session/psession"
)

func Init() {
	global.SESSION = psession.NewPSession(global.CACHE)
	global.LOG.Info("init session successfully")
}
