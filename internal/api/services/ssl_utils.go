package services

import (
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/files"
	"log"
	"path"
)

func saveCertificateFile(websiteSSL *models.WebsiteSSL, logger *log.Logger) {
	if websiteSSL.PushDir {
		fileOp := files.NewFileOp()
		var (
			pushErr error
			MsgMap  = map[string]interface{}{"path": websiteSSL.Dir, "status": "Success"}
		)
		if pushErr = fileOp.SaveFile(path.Join(websiteSSL.Dir, "privkey.pem"), websiteSSL.PrivateKey, 0666); pushErr != nil {
			MsgMap["status"] = "Failed"
			logger.Println("PushDirLog", MsgMap)
			logger.Println("Push dir failed:" + pushErr.Error())
		}
		if pushErr = fileOp.SaveFile(path.Join(websiteSSL.Dir, "fullchain.pem"), websiteSSL.Pem, 0666); pushErr != nil {
			MsgMap["status"] = "Failed"
			logger.Println("PushDirLog", MsgMap)
			logger.Println("Push dir failed:" + pushErr.Error())
		}
		if pushErr == nil {
			logger.Println("PushDirLog", MsgMap)
		}
	}
}
