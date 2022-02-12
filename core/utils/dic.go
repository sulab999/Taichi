package utils

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"sulab/core/model"

	"os"
)

func UserPassIsExist() bool {
	if IsExist("userpass.txt") {
		return true
	}
	return false
}

func PwdIsExist() bool {
	if IsExist("userpass.txt") {
		return true
	}
	if IsExist("user.txt") {
		return true
	}
	if IsExist("pass.txt") {
		return true
	}
	return false
}

func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func TxtRead(filename string) (lines []string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Open ", filename, "error, %v", err)
	}
	fi, _ := os.Stat(filename)
	if fi.Size() == 0 {
		fmt.Println("Error: " + filename + " is null!")
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ip := strings.TrimSpace(scanner.Text())
		if ip != "" {
			lines = append(lines, ip)
		}
	}
	return lines
}
func UserDic() (users []string) {
	dicname := "user.txt"
	file, err := os.Open(dicname)
	if err != nil {
		fmt.Println("Open "+dicname+" error, %v", err)
	}
	fi, _ := os.Stat(dicname)
	if fi.Size() == 0 {
		fmt.Println("Error: " + dicname + " is null!")
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		user := strings.TrimSpace(scanner.Text())
		if user != "" {
			users = append(users, user)
		}
	}
	return users
}

func PassDic() (password []string) {
	dicname := "pass.txt"
	file, err := os.Open(dicname)
	if err != nil {
		fmt.Println("Open "+dicname+" error, %v", err)
	}
	fi, _ := os.Stat(dicname)
	if fi.Size() == 0 {
		fmt.Println("Error: " + dicname + " is null!")
		os.Exit(1)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		passwd := strings.TrimSpace(scanner.Text())
		if passwd != "" {
			password = append(password, passwd)
		}
	}
	return password
}

func UserPassDic() (userpass []string) {
	dicname := "userpass.txt"
	file, err := os.Open(dicname)
	if err != nil {
		fmt.Println("Open "+dicname+" error, %v", err)
	}
	fi, _ := os.Stat(dicname)
	if fi.Size() == 0 {
		fmt.Println("Error: " + dicname + " is null!")
		os.Exit(1)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		passwd := strings.TrimSpace(scanner.Text())
		if passwd != "" {
			userpass = append(userpass, passwd)
		}
	}
	return userpass
}
func ReadUserDict(userDict string) (users []string, err error) {
	file, err := os.Open(userDict)
	if err != nil {
		fmt.Println("Open user dict file err, %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		user := strings.TrimSpace(scanner.Text())
		if user != "" {
			users = append(users, user)
		}
	}
	return users, err
}

func ReadIpList(fileName string) (ipList []model.IpAddr) {
	ipListFile, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Open ip List file err, %v", err)
	}

	defer ipListFile.Close()

	scanner := bufio.NewScanner(ipListFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		ipPort := strings.TrimSpace(line)
		t := strings.Split(ipPort, ":")
		ip := t[0]
		portProtocol := t[1]
		tmpPort := strings.Split(portProtocol, "|")

		if len(tmpPort) == 2 {
			port, _ := strconv.Atoi(tmpPort[0])
			protocol := strings.ToUpper(tmpPort[1])
			if SupportProtocols[protocol] {
				addr := model.IpAddr{Ip: ip, Port: port, Protocol: protocol}
				ipList = append(ipList, addr)
			} else {
				fmt.Printf("Not support %v, ignore: %v:%v", protocol, ip, port)
			}
		} else {

			port, err := strconv.Atoi(tmpPort[0])
			if err == nil {
				protocol, ok := PortNames[port]
				if ok && SupportProtocols[protocol] {
					addr := model.IpAddr{Ip: ip, Port: port, Protocol: protocol}
					ipList = append(ipList, addr)
				}
			}
		}

	}

	return ipList
}

func ReadIps(ScanType string, Target []string) (ipList []model.IpAddr) {

	port := model.GetPorts(Config.DB, "up")

	for _, p := range port {
		if p.Protocol != "" {
			protocol, _ := strconv.Atoi(p.Protocol)
			for _, i := range Target {
				ip := i
				port := protocol
				protocol := ScanType
				addr := model.IpAddr{Ip: ip, Port: port, Protocol: protocol}
				ipList = append(ipList, addr)
			}
		} else {
			for _, i := range Target {
				ip := i
				port := 21
				protocol := ScanType
				addr := model.IpAddr{Ip: ip, Port: port, Protocol: protocol}
				ipList = append(ipList, addr)
			}
		}
	}

	return ipList
}

func ReadSshIps(ScanType string, Target []string) (ipList []model.IpAddr) {

	port := model.GetPorts(Config.DB, "up")
	if port == nil {
		for _, i := range Target {
			ip := i
			port := 22
			protocol := ScanType
			addr := model.IpAddr{Ip: ip, Port: port, Protocol: protocol}
			ipList = append(ipList, addr)
		}

	} else {
		for _, i := range Target {
			ip := i
			port := 22
			protocol := ScanType
			addr := model.IpAddr{Ip: ip, Port: port, Protocol: protocol}
			ipList = append(ipList, addr)
		}

	}

	return ipList
}
