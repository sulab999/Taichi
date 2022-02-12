package enum

import (
	"fmt"
	"net"
	"sulab/core/utils"

	"github.com/fatih/color"
)

func JavadebugScan(ScanType string, Target []string) {

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

func JavaDebug(ip string, port string, username string, password string) (result bool, err error) {
	defer func() {
		if err := recover(); err != nil {

			return
		}
	}()

	conn, _ := net.Dial("tcp", ip+":"+port)

	conn.Write([]byte{0x4a, 0x44, 0x57, 0x50, 0x2d, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65})
	defer conn.Close()
	buffer := make([]byte, 32)

	res, _ := conn.Read(buffer)

	if res == 14 {
		result = true
	}
	return result, err
}
