package enum

import (
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/mgo.v2"
	"sulab/core/utils"
	"time"
)

func MongodbScan(ScanType string, Target []string) {

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


func ScanMongodb(ip string, port string, username string, password string) (result bool, err error) {
	session, err := mgo.DialWithTimeout("mongodb:
	if err == nil && session.Ping() == nil {
		defer session.Close()
		if err == nil && session.Run("serverStatus", nil) == nil {
			result = true
		}
	}
	return result, err
}

func MongoUnauth(ip string, port string) (err error, result bool) {
	session, err := mgo.Dial(ip + ":" + port)
	if err == nil && session.Run("serverStatus", nil) == nil {
		result = true
	}
	return err, result
}
