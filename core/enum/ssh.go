package enum

import (
	"fmt"
	"net"
	"sulab/core/utils"
	"time"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
)

func ScanSsh(ip string, port string, username string, password string) (result bool, err error) {

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * 1,
	}

	client, err := ssh.Dial("tcp", ip+":"+port, config)
	if err == nil {
		defer client.Close()
		session, err := client.NewSession()
		errRet := session.Run("echo ISOK")
		if err == nil && errRet == nil {
			defer session.Close()
			result = true
		}
	}
	return result, err
}

func SshScan(ScanType string, Target []string) {

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
func SshScan2(ScanType string, Target string) {
Loop:
	for _, u := range utils.UserDic() {
		for _, p := range utils.PassDic() {

			res, err := ScanSsh(Target, "22", u, p)

			if res == true && err == nil {

				fmt.Println(Target + " 22" + " 用户名: " + u + " 密码: " + p)
				break Loop
			}
		}
	}
}
