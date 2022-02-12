package enum

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sulab/core/model"
	"sulab/core/utils"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"
)

var (
	mutex       sync.Mutex
	successHash map[string]bool
	bruteResult map[string]model.Service2
)

func saveRes(target model.Service2, h string) {
	setTaskHask(h)
	_, ok := bruteResult[h]
	if !ok {
		mutex.Lock()

		color.Cyan("[+] %s %d %s %s \n", target.Ip, target.Port, target.UserName, target.PassWord)
		s := fmt.Sprintf("[+] %s %d %s %s  \n", target.Ip, target.Port, target.UserName, target.PassWord)
		WriteToFile(s, "res.txt")
		bruteResult[h] = model.Service2{Ip: target.Ip, Port: target.Port, Protocol: target.Protocol, UserName: target.UserName, PassWord: target.PassWord}
		mutex.Unlock()
	}
}
func in(target string, str_array []string) bool {

	sort.Strings(str_array)

	index := sort.SearchStrings(str_array, target)

	if index < len(str_array) && str_array[index] == target {

		return true

	}

	return false

}

func runBrute(taskChan chan model.Service2, wg *sync.WaitGroup) {
	for target := range taskChan {

		protocol := strings.ToUpper(target.Protocol)

		var k string
		protocol_list := []string{"RDP", "JAVADEBUG", "REDIS", "FTP", "SNMP", "POSTGRESQL", "SSH", "MONGO", "SMB", "MSSQL", "MYSQL", "ELASTICSEARCH"}
		result := in(protocol, protocol_list)
		if result {
			k = fmt.Sprintf("%v-%v-%v", target.Ip, target.Port, target.Protocol)
		} else {
			k = fmt.Sprintf("%v-%v-%v", target.Ip, target.Port, target.UserName)
		}

		h := utils.MakeTaskHash(k)
		if checkTashHash(h) {
			wg.Done()
			continue
		}

		res, err := ScanFuncMap[protocol](target.Ip, strconv.Itoa(target.Port), target.UserName, target.PassWord)
		if err == nil && res == true {
			saveRes(target, h)
		} else {

		}
		wg.Done()
	}

}

func RunTask(scanTasks []model.Service2, thread int) {

	wg := &sync.WaitGroup{}

	successHash = make(map[string]bool)
	bruteResult = make(map[string]model.Service2)

	taskChan := make(chan model.Service2, thread*2)

	for i := 0; i < thread; i++ {
		go runBrute(taskChan, wg)

	}

	bar := pb.StartNew(len(scanTasks))

	for _, task := range scanTasks {
		wg.Add(1)
		taskChan <- task
		bar.Increment()
	}

	close(taskChan)

	bar.Finish()

	wg.Wait()

	WriteToFile("全部掃描完成\n", "res.txt")

	color.Red("Scan complete. %d vulnerabilities found! \n", len(bruteResult))

}

func WriteToFile(wireteString, filename string) {

	fd, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	buf := []byte(wireteString)
	fd.Write(buf)

}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}

func GenerateTaskUserPass(addr []model.IpAddr, userList []string) (scanTasks []model.Service2) {
	for _, u := range userList {
		uk := strings.Split(u, ":")
		for _, ip := range addr {
			scanTask := model.Service2{Ip: ip.Ip, Port: ip.Port, Protocol: ip.Protocol, UserName: uk[0], PassWord: uk[1]}
			scanTasks = append(scanTasks, scanTask)
		}
	}
	return
}

func GenerateTask(addr []model.IpAddr, userList []string, passList []string) (scanTasks []model.Service2) {

	scanTasks = make([]model.Service2, 0)

	protocol_list := []string{"RDP", "JAVADEBUG", "REDIS", "FTP", "SNMP", "POSTGRESQL", "SSH", "MONGO", "SMB", "MSSQL", "MYSQL", "ELASTICSEARCH"}

	for _, ip := range addr {
		result := in(ip.Protocol, protocol_list)
		if result {
			scanTask := model.Service2{Ip: ip.Ip, Port: ip.Port, Protocol: ip.Protocol, UserName: "", PassWord: ""}
			scanTasks = append(scanTasks, scanTask)
		}
	}

	for _, u := range userList {
		for _, p := range passList {
			for _, ip := range addr {
				scanTask := model.Service2{Ip: ip.Ip, Port: ip.Port, Protocol: ip.Protocol, UserName: u, PassWord: p}
				scanTasks = append(scanTasks, scanTask)
			}
		}
	}

	return
}

func checkTashHash(hash string) bool {
	_, ok := successHash[hash]
	return ok
}

func setTaskHask(hash string) {
	mutex.Lock()
	successHash[hash] = true
	mutex.Unlock()
}
