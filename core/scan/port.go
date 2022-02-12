package scan

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func CheckPort(ip net.IP, port int) {
	tcpAddr := net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	conn, err := net.DialTCP("tcp", nil, &tcpAddr)
	if conn != nil {
		fmt.Println(tcpAddr.IP, tcpAddr.Port, "Open")
		conn.Close()
	}
	if err != nil {

	}
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func WriteResult(host []string) {
	filename := "portresult.txt"

	fout, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {

		fmt.Println(filename + " create error")
	}

	defer fout.Close()

	write := bufio.NewWriter(fout)

	for i := 0; i < len(host); i++ {
		write.WriteString(host[i] + "\r\n")
	}

	write.Flush()

}

func PortCheck(host string, port int) (result bool) {
	result = false
	ip := net.ParseIP(host)
	tcpAddr := net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	conn, err := net.DialTCP("tcp", nil, &tcpAddr)
	if conn != nil {
		fmt.Println(tcpAddr.IP, tcpAddr.Port, "Open")

		conn.Close()
		result = true
	}
	if err != nil {

	}
	return result
}

func PortIsOpen(ip net.IP, port int) (result bool, err error) {
	tcpAddr := net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	conn, err := net.DialTCP("tcp", nil, &tcpAddr)
	if conn != nil {

		conn.Close()
		result = true
	}
	if err != nil {

	}
	return result, err
}

type Workdist struct {
	Host string
}

const (
	taskload = 255
	tasknum  = 255
)

var wg sync.WaitGroup

func TaskPort(ip string, debugLog *log.Logger) {
	tasks := make(chan Workdist, taskload)
	wg.Add(tasknum)

	for gr := 1; gr <= tasknum; gr++ {
		go workerPort(tasks, debugLog)
	}

	for i := 1; i < 256; i++ {
		host := fmt.Sprintf("%s.%d", ip, i)
		task := Workdist{
			Host: host,
		}
		tasks <- task
	}
	close(tasks)
	wg.Wait()
}

func workerPort(tasks chan Workdist, debugLog *log.Logger) {
	defer wg.Done()
	task, ok := <-tasks
	if !ok {
		return
	}
	host := task.Host

	ScanPort2(host)

}

var DefaultPorts = []int{21, 22, 23, 25, 80, 443, 8080, 110, 135, 139, 445, 389, 489, 587, 1433, 1434, 1521, 1522, 1723, 2121, 3000, 3306, 3389, 4899, 5631, 5632, 5800, 5900, 7071, 43958, 65500, 4444, 8888, 6789, 4848, 5985, 5986, 8081, 8089, 8443, 10000, 6379, 7001, 7002}

func ScanPort2(host string) {

	var mutex sync.Mutex

	finish := make(chan int)

	channel := make(chan int, 100)
	var openPorts []int
	var timeoutPorts []int

	addOpenPorts := func(port int) {
		mutex.Lock()
		defer mutex.Unlock()
		openPorts = append(openPorts, port)
	}

	scan := func(ip string, port int) {

		address := ip + ":" + strconv.Itoa(port)
		_, err := net.DialTimeout("tcp", address, time.Second*2)
		if err != nil {

			if strings.Contains(err.Error(), "timeout") {

				timeoutPorts = append(timeoutPorts, port)
			}
		} else {
			fmt.Println(address + " open")
			addOpenPorts(port)
		}
		i := <-channel
		if i == 1 {
			finish <- 0
		}
	}
	num := len(DefaultPorts)
	for i := 0; i < num; i++ {
		if i == num-1 {
			channel <- 1
		} else {
			channel <- 0
		}

		go scan(host, DefaultPorts[i])

	}
	<-finish

}

func ScanAllPort(host string) []string {

	runtime.GOMAXPROCS(runtime.NumCPU())
	var mutex sync.Mutex

	finish := make(chan int)

	channel := make(chan int, 70000)
	var openPorts []string
	var openPorts2 []string
	var timeoutPorts []int
	var ports []int

	for i := 1; i <= 65535; i++ {
		ports = append(ports, i)
	}

	addOpenPorts := func(port string) {
		mutex.Lock()
		defer mutex.Unlock()
		openPorts = append(openPorts, port)
	}

	scan := func(ip string, port int) {

		address := ip + ":" + strconv.Itoa(port)
		conn, err := net.DialTimeout("tcp", address, time.Second*2)
		if err != nil {

			if strings.Contains(err.Error(), "timeout") {

				timeoutPorts = append(timeoutPorts, port)
			}
		} else {
			defer conn.Close()
			out := address + " open "
			fmt.Println(out)
			addOpenPorts(out)
			openPorts2 = append(openPorts2, address)
		}
		i := <-channel
		if i == 1 {
			finish <- 0
		}

	}
	num := len(ports)

	for i := 0; i < num; i++ {
		if i == num-1 {
			channel <- 1
		} else {
			channel <- 0
		}
		go scan(host, ports[i])

	}
	<-finish

	return openPorts2
}

func Worker(tasksCh <-chan int, wg *sync.WaitGroup, ips string) {
	defer wg.Done()
	for {
		task, ok := <-tasksCh
		if !ok {
			return
		}
		ScanAllPort(ips)

		fmt.Println("processing task", task)
	}
}

func TcpPort(host string, port int) bool {
	p := strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", host+":"+p, time.Second*2)
	if err != nil {

		return false
	} else {

		fmt.Println(host, p, "Open")
		conn.Close()
		return true
	}
}

func Scan(host []string) {
	ips := host
	var ports []int

	for i := 1; i <= 65535; i++ {
		ports = append(ports, i)
	}

	wg := NewSizeWG(70000)
	if len(ips) != 0 && len(ports) != 0 {
		for _, ip := range ips {
			for _, port := range ports {
				wg.Add()
				go func(ip string, port int) {
					defer wg.Done()
					TcpPort(ip, port)

				}(ip, port)
			}
		}
	}
	wg.Wait()
}
