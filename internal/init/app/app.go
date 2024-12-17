package app

import (
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/utils/files"
	"LinuxOnM/internal/utils/firewall"
	"path"
)

func Init() {
	constant.DataDir = global.CONF.System.DataDir
	constant.SSLLogDir = path.Join(global.CONF.System.DataDir, "log", "ssl")

	dirs := []string{constant.DataDir, constant.SSLLogDir}

	fileOp := files.NewFileOp()
	for _, dir := range dirs {
		createDir(fileOp, dir)
	}

	if f, err := firewall.NewFirewallClient(); err == nil {
		_ = f.EnableForward()
	}

}

func createDir(fileOp files.FileOp, dirPath string) {
	if !fileOp.Stat(dirPath) {
		_ = fileOp.CreateDir(dirPath, 0755)
	}
}
