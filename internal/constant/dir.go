package constant

import (
	"LinuxOnM/internal/global"
	"path"
)

var (
	DataDir       = global.CONF.System.DataDir
	RecycleBinDir = "/opt/.LinuxOnM_clash"
	SSLLogDir     = path.Join(global.CONF.System.DataDir, "log", "ssl")
)
