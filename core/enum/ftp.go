package enum

import (
	"fmt"
	"sulab/core/utils"
	"time"

	"github.com/fatih/color"
	"github.com/jlaffaye/ftp"
)

func FtpScan(ScanType string, Target []string) {

	ipList := utils.ReadIps(ScanType, Target)

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

func ScanFtp(ip string, port string, username string, password string) (result bool, err error) {
	conn, err := ftp.DialTimeout(ip+":"+port, time.Second*1)

	if err == nil {
		err = conn.Login(username, password)
		if err == nil {
			result = true
			conn.Logout()
		}
	}
	return result, err
}
