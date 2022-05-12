package enum

import (
	"context"
	"fmt"
	"sulab/core/utils"
	"time"

	"github.com/fatih/color"
	redis "github.com/go-redis/redis/v8"
)

func RedisScan(ScanType string, Target []string) {

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
func ScanRedis(ip string, port string, username string, password string) (result bool, err error) {
	client := redis.NewClient(&redis.Options{Addr: ip + ":" + port, Password: password, DB: 0, DialTimeout: time.Second * 3})
	var ctx = context.Background()
	defer client.Close()
	//_, err = client.Ping().Result()
	_, err = client.Ping(ctx).Result()
	if err == nil {
		result = true
	}
	return result, err
}
