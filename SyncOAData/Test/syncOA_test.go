package Test

import (
	"goStudy/SyncOAData"
	"testing"
)

func TestSyncOA(t *testing.T) {
	SyncOAData.SyncOneAgentData()
}

func TestSyncOAI(t *testing.T) {
	SyncOAData.SyncAgentIntegrationData()
}

//
//func Md5Test(t *testing.T) {
//	utils.MD5String(version.Config)
//}
