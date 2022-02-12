package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sulab/core/model"

	"github.com/jinzhu/gorm"
)

var Config config

var Const_notification_delay_unit = 10
var Const_example_target_cidr = "127.0.0.1/32"
var Const_example_target_desc = "Target CIDR or /32 for single target"

var Const_UDP_PORTS = "19,53,69,79,111,123,135,137,138,161,177,445,500,514,520,1434,1900,5353"
var Const_NMAP_SWEEP = "-n -sn -PE -PP"
var Const_NMAP_TCP_FULL = "--randomize-hosts -Pn -sS -sC -A -T4 -g53 -p-"
var Const_NMAP_TCP_STANDARD = "--randomize-hosts -Pn -sS -A -T4 -g53 --top-ports 1000"
var Const_NMAP_TCP_PROD = "--randomize-hosts -Pn -sT -sV -T3 -p-"
var Const_NMAP_TCP_VULN = "--randomize-hosts -Pn -sT -sV -p- --script=vulscan/vulscan.nse"
var Const_NMAP_UDP_STANDARD = fmt.Sprintf("--randomize-hosts -Pn -sU -sC -A -T4 -p%s", Const_UDP_PORTS)
var Const_NMAP_UDP_PROD = fmt.Sprintf("--randomize-hosts -Pn -sU -sC -sV -T3 -p%s", Const_UDP_PORTS)

var WORDLIST_FUZZ_NAMELIST = "/usr/share/wfuzz/wordlist/fuzzdb/wordlists-user-passwd/names/namelist.txt"
var WORDLIST_MSF_PWDS = "/usr/share/wordlists/metasploit/unix_passwords.txt"
var WORDLIST_FINGER_USER = WORDLIST_FUZZ_NAMELIST
var WORDLIST_FTP_USER = WORDLIST_FUZZ_NAMELIST
var WORDLIST_SMTP = WORDLIST_FUZZ_NAMELIST
var WORDLIST_SNMP = "/usr/share/doc/onesixtyone/dict.txt"
var WORDLIST_DNS_BRUTEFORCE = WORDLIST_FUZZ_NAMELIST
var WORDLIST_HYDRA_SSH_USER = WORDLIST_FUZZ_NAMELIST
var WORDLIST_HYDRA_SSH_PWD = WORDLIST_MSF_PWDS
var WORDLIST_HYDRA_FTP_USER = WORDLIST_FUZZ_NAMELIST
var WORDLIST_HYDRA_FTP_PWD = WORDLIST_MSF_PWDS

type config struct {
	Outfolder string
	Log       *Logger
	DB        *gorm.DB
	DBPath    string
}
type Bar struct {
	percent int64
	cur     int64
	total   int64
	rate    string
	graph   string
}

func (bar *Bar) NewOption(start, total int64) {
	bar.cur = start
	bar.total = total
	if bar.graph == "" {
		bar.graph = "█"
	}
	bar.percent = bar.getPercent()
	for i := 0; i < int(bar.percent); i += 2 {
		bar.rate += bar.graph
	}
}
func (bar *Bar) getPercent() int64 {
	return int64(float32(bar.cur) / float32(bar.total) * 100)
}
func (bar *Bar) NewOptionWithGraph(start, total int64, graph string) {
	bar.graph = graph
	bar.NewOption(start, total)
}
func (bar *Bar) Play(cur int64) {
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last && bar.percent%2 == 0 {
		bar.rate += bar.graph
	}
	fmt.Printf("\r[%-50s]%3d%%  %8d/%d", bar.rate, bar.percent, bar.cur, bar.total)
}
func (bar *Bar) Finish() {
	fmt.Println()
}

func InitConfig() {
	Config = config{}

	Config.Log = InitLogger()

	if os.Getenv("OUT_FOLDER") != "" {
		Config.Outfolder = filepath.Join(os.Getenv("OUT_FOLDER"), "goscan")
	} else {
		usr, _ := user.Current()
		Config.Outfolder = filepath.Join(usr.HomeDir, ".goscan")
	}
	EnsureDir(Config.Outfolder)

	if os.Getenv("GOSCAN_DB_PATH") != "" {
		Config.DBPath = os.Getenv("GOSCAN_DB_PATH")
	} else {
		Config.DBPath = filepath.Join(Config.Outfolder, "goscan.db")
		fmt.Println(Config.DBPath)
	}
	Config.DB = model.InitDB(Config.DBPath)

	Config.Log.LogDebug("Connected to DB")
}

func ChangeOutFolder(path string) {

	Config.Outfolder = path
	EnsureDir(Config.Outfolder)

	Config.DBPath = filepath.Join(Config.Outfolder, "goscan.db")
	fmt.Println(Config.DBPath)
	Config.DB = model.InitDB(Config.DBPath)
	Config.Log.LogDebug("Connected to DB")
}

func ParseCmd(s string) (string, []string) {

	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return "", make([]string, 0)
	}

	tokens := strings.Fields(s)

	cmd, args := tokens[0], tokens[1:]
	return cmd, args
}

func ParseNextArg(args []string) (string, []string) {
	if len(args) < 2 {
		return args[0], make([]string, 0)
	}
	return args[0], args[1:]
}

func ParseAllArgs(args []string) string {
	return strings.Join(args, " ")
}

func ShellCmd(cmd string) (string, error) {
	Config.Log.LogDebug(fmt.Sprintf("Executing command: %s", cmd))
	output, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		if !strings.Contains(err.Error(), "exit status 1") {
			Config.Log.LogError(fmt.Sprintf("Error while executing command: %s", err.Error()))
		}
		return string(output), err
	}
	return string(output), err
}

func EnsureDir(dir string) {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
		Config.Log.LogDebug(fmt.Sprintf("Created directory: %s", dir))
	}
}

func RemoveDir(dir string) {
	os.RemoveAll(dir)
	Config.Log.LogDebug(fmt.Sprintf("Deleted directory: %s", dir))
}

func CleanPath(s string) string {
	return strings.Replace(s, "/", "_", -1)
}

func WriteArrayToFile(path string, s []string) {
	Config.Log.LogDebug(fmt.Sprintf("Writing output to file: %s", path))
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		Config.Log.LogError("Cannot create file")
	}
	defer f.Close()

	sep := "\n"
	for _, line := range s {
		if _, err = f.WriteString(line + sep); err != nil {
			Config.Log.LogError(fmt.Sprintf("Error while writing to file: %s", err))
		}
	}
}

func MD5(s string) (m string) {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func MakeTaskHash(k string) string {
	hash := MD5(k)
	return hash
}
