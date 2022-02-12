package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sulab/core/enum"
	"sulab/core/model"
	"sulab/core/scan"
	"sulab/core/utils"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/olekukonko/tablewriter"
)

var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}
var Kindstring struct {
	SetKind string
}

func ChangeLivePrefix() (string, bool) {
	return LivePrefixState.LivePrefix, LivePrefixState.IsEnable
}
func Executor(s string) {

	if s == "" {
		LivePrefixState.IsEnable = false
		LivePrefixState.LivePrefix = s
		return
	}

	cmd, args := utils.ParseCmd(s)

	switch cmd {
	case "load":
		cmdLoad(args)

	case "help":
		cmdHelp()
	case "exit", "quit":
		os.Exit(0)
		return
	case "":
	default:
		return
	}

}

func Executor2(s string) {

	cmd, args := utils.ParseCmd(s)

	switch cmd {
	case "load":
		cmdLoad(args)
	case "show":
		cmdShow2(args)
	case "set":
		cmdSet2(args)
	case "help":
		cmdHelp2()
	case "go":
		cmdPortscan2()

	case "exit", "quit":
		os.Exit(0)
		return
	case "":
	default:
		return
	}

}

func Executor3(s string) {

	cmd, args := utils.ParseCmd(s)

	switch cmd {
	case "load":
		cmdLoad(args)
	case "show":
		cmdShow3(args)
	case "set":
		cmdSet3(args)
	case "help":
		cmdHelp2()
	case "go":

		switch Kindstring.SetKind {
		case "ftp":
			cmdFtpburst()
		case "ssh":
			cmdSshburst()
		case "mongodb":
			cmdMongodbburst()
		case "mssql":
			cmdMssqlburst()
		case "mysql":
			cmdMysqlburst()
		case "postgres":
			cmdPostgresburst()
		case "redis":
			cmdRedisburst()
		case "smb":
			cmdSmbburst()
		case "javadebug":
			cmdJavadebugburst()
		case "rdp":
			cmdRdpburst()
		case "snmp":
			cmdSnmpburst()
		}

	case "exit", "quit":
		os.Exit(0)
		return
	case "":
	default:
		return
	}

}

func cmdHelp2() {

	data := [][]string{
		[]string{"load moudel", "加載模塊", "load <moudel>"},
		[]string{"set ip", "設置ip", "set ip xxx.xxx.xxx.xxx"},

		[]string{"set file", "設置文件", "set file url.txt"},
		[]string{"show", "顯示參數", "show"},
		[]string{"go", "開始執行", "go"},
		[]string{"exit", "退出", "exit"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"command", "description", "example"})
	table.SetAlignment(3)
	table.SetAutoWrapText(true)
	table.AppendBulk(data)
	table.Render()
}
func cmdHelp() {
	utils.Config.Log.LogInfo("Taiji Penetration Test Framework")
	utils.Config.Log.LogInfo("Available commands:")

	data := [][]string{
		[]string{"Load", "加載模塊", "load <模塊>"},
		[]string{"Set", "設置參數", "set <TYPE> <TARGET>"},

		[]string{"Show", "顯示設置信息", "show"},
		[]string{"Go", "執行", "go"},
		[]string{"Help", "幫助", "help"},

		[]string{"Exit", "退出", "exit"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Area", "Command", "Syntax"})
	table.SetAlignment(3)
	table.SetAutoWrapText(true)
	table.AppendBulk(data)
	table.Render()
}

func cmdLoad(args []string) bool {

	kind, args := utils.ParseNextArg(args)

	switch kind {

	case "portscan":

		p := prompt.New(
			Executor2,
			Completer2,
			prompt.OptionPrefix("[Taiji] > portscan > "),
			prompt.OptionInputTextColor(prompt.White),
		)
		p.Run()
		return true
	case "burst":
		switch args[0] {

		case "ftp":
			Kindstring.SetKind = "ftp"
			p := prompt.New(
				Executor3,
				Completer3,
				prompt.OptionPrefix("[Taiji] > ftpburst > "),
				prompt.OptionInputTextColor(prompt.White),
			)
			p.Run()
			return true
		case "ssh":
			Kindstring.SetKind = "ssh"
			p := prompt.New(
				Executor3,
				Completer3,
				prompt.OptionPrefix("[Taiji] > sshburst > "),
				prompt.OptionInputTextColor(prompt.White),
			)
			p.Run()
			return true
		}

	case "urlscan":

		p := prompt.New(
			Executor2,
			Completer2,
			prompt.OptionPrefix("[Taiji] > urlscan > "),
			prompt.OptionInputTextColor(prompt.White),
		)
		p.Run()
		return true

	case "subscan":

		return true

	}

	return true
}

func cmdPortscan2() {
	filename := "portresult.txt"

	_, err := os.Stat(filename)
	if err == nil {
		os.Remove(filename)
	}

	var ips []string
	hosts := model.GetAllHosts(utils.Config.DB)

	start := time.Now()
	for _, h := range hosts {
		if h.Step == model.NEW.String() || h.Address != "" {

			ips = scan.ScanAllPort(h.Address)
		}
	}

	scan.GetProbes(ips)
	elapsed := time.Since(start)
	fmt.Println("該函數執行完成耗時：", elapsed)
	model.DelHosts(utils.Config.DB)
	model.DelPorts(utils.Config.DB)
}

func cmdFtpburst() {
	filename := "res.txt"

	_, err := os.Stat(filename)
	if err == nil {
		os.Remove(filename)
	}

	hosts := model.GetAllHosts(utils.Config.DB)

	var ips []string
	start := time.Now()
	for _, h := range hosts {
		if h.Step == model.NEW.String() || h.Address != "" {

			ips = append(ips, h.Address)
		}
	}

	enum.FtpScan("FTP", ips)

	elapsed := time.Since(start)
	fmt.Println("該函數執行完成耗時：", elapsed)
	model.DelHosts(utils.Config.DB)
	model.DelPorts(utils.Config.DB)

}

func cmdSshburst() {
	filename := "res.txt"

	_, err := os.Stat(filename)
	if err == nil {
		os.Remove(filename)
	}

	hosts := model.GetAllHosts(utils.Config.DB)
	var ips []string
	start := time.Now()
	for _, h := range hosts {
		if h.Step == model.NEW.String() || h.Address != "" {

			ips = append(ips, h.Address)
		}
	}
	enum.SshScan("SSH", ips)

	elapsed := time.Since(start)
	fmt.Println("該函數執行完成耗時：", elapsed)
	model.DelHosts(utils.Config.DB)
	model.DelPorts(utils.Config.DB)

}

func cmdMongodbburst() {
	filename := "res.txt"

	_, err := os.Stat(filename)
	if err == nil {
		os.Remove(filename)
	}

	hosts := model.GetAllHosts(utils.Config.DB)
	var ips []string
	start := time.Now()
	for _, h := range hosts {
		if h.Step == model.NEW.String() || h.Address != "" {
			ips = append(ips, h.Address)
		}
	}
	enum.MongodbScan("MONGO", ips)

	elapsed := time.Since(start)
	fmt.Println("該函數執行完成耗時：", elapsed)
	model.DelHosts(utils.Config.DB)
	model.DelPorts(utils.Config.DB)

}

func cmdMssqlburst() {
	filename := "res.txt"

	_, err := os.Stat(filename)
	if err == nil {
		os.Remove(filename)
	}

	hosts := model.GetAllHosts(utils.Config.DB)
	var ips []string
	start := time.Now()
	for _, h := range hosts {
		if h.Step == model.NEW.String() || h.Address != "" {
			ips = append(ips, h.Address)
		}
	}
	enum.MssqlScan("MSSQL", ips)

	elapsed := time.Since(start)
	fmt.Println("該函數執行完成耗時：", elapsed)
	model.DelHosts(utils.Config.DB)
	model.DelPorts(utils.Config.DB)

}

func cmdSmbburst() {
	filename := "res.txt"

	_, err := os.Stat(filename)
	if err == nil {
		os.Remove(filename)
	}

	hosts := model.GetAllHosts(utils.Config.DB)
	var ips []string
	start := time.Now()
	for _, h := range hosts {
		if h.Step == model.NEW.String() || h.Address != "" {
			ips = append(ips, h.Address)
		}
	}
	enum.SmbScan("SMB", ips)

	elapsed := time.Since(start)
	fmt.Println("該函數執行完成耗時：", elapsed)
	model.DelHosts(utils.Config.DB)
	model.DelPorts(utils.Config.DB)

}

func cmdMysqlburst() {
	filename := "res.txt"

	_, err := os.Stat(filename)
	if err == nil {
		os.Remove(filename)
	}

	hosts := model.GetAllHosts(utils.Config.DB)
	var ips []string
	start := time.Now()
	for _, h := range hosts {
		if h.Step == model.NEW.String() || h.Address != "" {
			ips = append(ips, h.Address)
		}
	}
	enum.MysqlScan("MYSQL", ips)

	elapsed := time.Since(start)
	fmt.Println("該函數執行完成耗時：", elapsed)
	model.DelHosts(utils.Config.DB)
	model.DelPorts(utils.Config.DB)

}

func cmdPostgresburst() {
	filename := "res.txt"

	_, err := os.Stat(filename)
	if err == nil {
		os.Remove(filename)
	}

	hosts := model.GetAllHosts(utils.Config.DB)
	var ips []string
	start := time.Now()
	for _, h := range hosts {
		if h.Step == model.NEW.String() || h.Address != "" {
			ips = append(ips, h.Address)
		}
	}
	enum.PostgresScan("POSTGRESQL", ips)

	elapsed := time.Since(start)
	fmt.Println("該函數執行完成耗時：", elapsed)
	model.DelHosts(utils.Config.DB)
	model.DelPorts(utils.Config.DB)

}

func cmdRedisburst() {
	filename := "res.txt"

	_, err := os.Stat(filename)
	if err == nil {
		os.Remove(filename)
	}

	hosts := model.GetAllHosts(utils.Config.DB)
	var ips []string
	start := time.Now()
	for _, h := range hosts {
		if h.Step == model.NEW.String() || h.Address != "" {
			ips = append(ips, h.Address)
		}
	}
	enum.RedisScan("REDIS", ips)

	elapsed := time.Since(start)
	fmt.Println("該函數執行完成耗時：", elapsed)
	model.DelHosts(utils.Config.DB)
	model.DelPorts(utils.Config.DB)

}

func cmdJavadebugburst() {
	filename := "res.txt"

	_, err := os.Stat(filename)
	if err == nil {
		os.Remove(filename)
	}

	hosts := model.GetAllHosts(utils.Config.DB)
	var ips []string
	start := time.Now()
	for _, h := range hosts {
		if h.Step == model.NEW.String() || h.Address != "" {
			ips = append(ips, h.Address)
		}
	}
	enum.JavadebugScan("JAVADEBUG", ips)

	elapsed := time.Since(start)
	fmt.Println("該函數執行完成耗時：", elapsed)
	model.DelHosts(utils.Config.DB)
	model.DelPorts(utils.Config.DB)

}

func cmdRdpburst() {
	filename := "res.txt"

	_, err := os.Stat(filename)
	if err == nil {
		os.Remove(filename)
	}

	hosts := model.GetAllHosts(utils.Config.DB)
	var ips []string
	start := time.Now()
	for _, h := range hosts {
		if h.Step == model.NEW.String() || h.Address != "" {
			ips = append(ips, h.Address)
		}
	}
	enum.RdpScan("RDP", ips)

	elapsed := time.Since(start)
	fmt.Println("該函數執行完成耗時：", elapsed)
	model.DelHosts(utils.Config.DB)
	model.DelPorts(utils.Config.DB)

}

func cmdSnmpburst() {
	filename := "res.txt"

	_, err := os.Stat(filename)
	if err == nil {
		os.Remove(filename)
	}

	hosts := model.GetAllHosts(utils.Config.DB)
	var ips []string
	start := time.Now()
	for _, h := range hosts {
		if h.Step == model.NEW.String() || h.Address != "" {
			ips = append(ips, h.Address)
		}
	}
	enum.SnmpScan("SNMP", ips)

	elapsed := time.Since(start)
	fmt.Println("該函數執行完成耗時：", elapsed)
	model.DelHosts(utils.Config.DB)
	model.DelPorts(utils.Config.DB)

}

func cmdShow2(args []string) {

	ShowHosts()

}

func cmdShow3(args []string) {

	ShowHosts()
	ShowPorts2()
}

func ShowHosts() {
	hosts := model.GetAllHosts(utils.Config.DB)
	if len(hosts) == 0 {
		utils.Config.Log.LogError("No hosts are up!")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Address", "Status", "OS", "Info", "Ports"})
	table.SetRowLine(true)
	table.SetAlignment(1)
	table.SetAutoWrapText(true)

	for _, h := range hosts {
		rAddress := h.Address
		rStatus := h.Status
		rOS := h.OS
		rInfo := h.Info
		rPorts := ""
		v := []string{rAddress, rStatus, rOS, rInfo, rPorts}
		table.Append(v)
	}
	table.Render()
}

func ShowPorts2() {
	ports := model.GetAllPorts(utils.Config.DB)

	if len(ports) == 0 {
		utils.Config.Log.LogError("No ports are up!")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Number", "Port", "Status"})
	table.SetRowLine(true)

	table.SetAlignment(3)
	table.SetAutoWrapText(false)

	for _, h := range ports {
		rNumber := strconv.Itoa(h.Number)
		rPort := h.Protocol
		rStatus := h.Status
		v := []string{rNumber, rPort, rStatus}
		table.Append(v)
	}

	table.Render()
}

func SetUrlFile(fname string) {

	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return
	}
	file, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	} else {
		contents := string(file)

		re := regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
		in := re.ReplaceAllString(contents, "")

		in = strings.Replace(in, " ", "", -1)

		lines := strings.Split(in, "\r\n")
		fmt.Println(lines)
		for _, line := range lines {
			model.AddHost(utils.Config.DB, line, "up", model.NEW.String())
		}
	}

}

func cmdSet3(args []string) {

	if len(args) != 2 {
		utils.Config.Log.LogError("Invalid command provided")
		return
	}

	kind, args := utils.ParseNextArg(args)
	src, args := utils.ParseNextArg(args)

	switch kind {
	case "file":

		SetUrlFile(src)
	case "ip":
		ip, parsed := utils.ParseAddress(src)
		if parsed == false {
			utils.Config.Log.LogError("Invalid address provided")
			return
		}
		utils.Config.Log.LogInfo(fmt.Sprintf("Imported target: %s", ip))
		model.AddHost(utils.Config.DB, ip, "up", model.NEW.String())
	case "port":
		port := src
		model.AddPort(utils.Config.DB, 1, port, "up")
	}
}

func cmdSet2(args []string) {

	if len(args) != 2 {
		utils.Config.Log.LogError("Invalid command provided")
		return
	}

	kind, args := utils.ParseNextArg(args)
	src, args := utils.ParseNextArg(args)

	switch kind {
	case "file":

		SetUrlFile(src)
	case "ip":
		ip, parsed := utils.ParseAddress(src)
		if parsed == false {
			utils.Config.Log.LogError("Invalid address provided")
			return
		}
		utils.Config.Log.LogInfo(fmt.Sprintf("Imported target: %s", ip))
		model.AddHost(utils.Config.DB, ip, "up", model.NEW.String())
	case "port":
		port := src
		model.AddPort(utils.Config.DB, 1, port, "up")
	}
}
