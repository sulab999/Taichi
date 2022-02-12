package enum

import (
	"fmt"
	"strconv"
	"sulab/core/utils"
	"time"

	"github.com/fatih/color"
	"github.com/gosnmp/gosnmp"
)

func SnmpScan(ScanType string, Target []string) {

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

func ScanSnmp(ip string, port string, username string, password string) (result bool, err error) {

	p, err := strconv.Atoi(port)
	gosnmp.Default.Target = ip
	gosnmp.Default.Port = uint16(p)
	gosnmp.Default.Community = "public"
	gosnmp.Default.Timeout = 3 * time.Second

	err = gosnmp.Default.Connect()
	if err == nil {
		oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0"}
		_, err := gosnmp.Default.Get(oids)
		if err == nil {
			result = true
		}
	}

	return result, err
}
