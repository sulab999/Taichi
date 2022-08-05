package enum

import (
	"github.com/fatih/color"
	_ "github.com/lib/pq"
	"sulab/core/utils"

	"database/sql"
	"fmt"
)

func PostgresScan(ScanType string, Target []string) {

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
func ScanPostgres(ip string, port string, username string, password string) (result bool, err error) {
	// pgurl
	pgurl := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable", username, password, ip, port)
	db, err := sql.Open("postgres", pgurl)
	if err != nil {
		return false, err
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return false, err
	} else {
		return true, err
	}
}
