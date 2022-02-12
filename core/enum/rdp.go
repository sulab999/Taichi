package enum

import (
	"fmt"
	"sulab/core/utils"

	"github.com/fatih/color"
	"github.com/icodeface/grdp"
	"github.com/icodeface/grdp/glog"
)

func RdpScan(ScanType string, Target []string) {

	ipList := utils.ReadSshIps(ScanType, Target)
	thread := 1000
	userDict, uErr := utils.ReadUserDict("user.txt")
	passDict, pErr := utils.ReadUserDict("pass.txt")
	if utils.UserPassIsExist() {
		userDict, _ := utils.ReadUserDict("userpass.txt")
		scanTasks := GenerateTaskUserPass(ipList, userDict)
		color.Cyan("Number of all task : %d", len(scanTasks))
		RunTask(scanTasks, thread)
	} else {
		if uErr == nil && pErr == nil {
			scanTasks := GenerateTask(ipList, userDict, passDict)
			color.Cyan("Number of all task : %d", len(scanTasks))

			RunTask(scanTasks, thread)
		} else {
			fmt.Println("Read File Err!")
		}
	}

}
func ScanRdp(ip string, port string, username string, password string) (result bool, err error) {
	client := grdp.NewClient(fmt.Sprintf("%s:%d", ip, port), glog.DEBUG)
	err = client.Login(username, password)
	if err == nil {
		result = true

	}
	return result, err
}
