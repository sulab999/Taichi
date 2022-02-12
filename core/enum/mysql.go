package enum

import (
	"database/sql"
	"fmt"
	"sulab/core/utils"

	"github.com/fatih/color"

	"time"
)

func MysqlScan(ScanType string, Target []string) {

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
func ScanMysql(ip string, port string, username string, password string) (result bool, err error) {
	result = false
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/?timeout=%ds", username, password, ip+":"+port, time.Second*3)
	db, err := sql.Open("mysql", connStr)
	if err == nil {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			defer db.Close()
			result = true
		}
	}
	return result, err
}
