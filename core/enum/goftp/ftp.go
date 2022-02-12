package goftp

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var RePwdPath = regexp.MustCompile(`\"(.*)\"`)

type FTP struct {
	conn net.Conn

	addr string

	debug     bool
	tlsconfig *tls.Config

	reader *bufio.Reader
	writer *bufio.Writer
}

func (ftp *FTP) Close() error {
	return ftp.conn.Close()
}

type (
	WalkFunc func(path string, info os.FileMode, err error) error

	RetrFunc func(r io.Reader) error
)

func parseLine(line string) (perm string, t string, filename string) {
	for _, v := range strings.Split(line, ";") {
		v2 := strings.Split(v, "=")

		switch v2[0] {
		case "perm":
			perm = v2[1]
		case "type":
			t = v2[1]
		default:
			filename = v[1 : len(v)-2]
		}
	}
	return
}

func (ftp *FTP) Walk(path string, walkFn WalkFunc) (err error) {
	/*
		if err = walkFn(path, os.ModeDir, nil); err != nil {
			if err == filepath.SkipDir {
				return nil
			}
		}
	*/
	if ftp.debug {
		log.Printf("Walking: '%s'\n", path)
	}

	var lines []string

	if lines, err = ftp.List(path); err != nil {
		return
	}

	for _, line := range lines {
		_, t, subpath := parseLine(line)

		switch t {
		case "dir":
			if subpath == "." {
			} else if subpath == ".." {
			} else {
				if err = ftp.Walk(path+subpath+"/", walkFn); err != nil {
					return
				}
			}
		case "file":
			if err = walkFn(path+subpath, os.FileMode(0), nil); err != nil {
				return
			}
		}
	}

	return
}

func (ftp *FTP) Quit() (err error) {
	if _, err := ftp.cmd(StatusConnectionClosing, "QUIT"); err != nil {
		return err
	}

	ftp.conn.Close()
	ftp.conn = nil

	return nil
}

func (ftp *FTP) Noop() (err error) {
	_, err = ftp.cmd(StatusOK, "NOOP")
	return
}

func (ftp *FTP) RawCmd(command string, args ...interface{}) (code int, line string) {
	if ftp.debug {
		log.Printf("Raw-> %s\n", fmt.Sprintf(command, args...))
	}

	code = -1
	var err error
	if err = ftp.send(command, args...); err != nil {
		return code, ""
	}
	if line, err = ftp.receive(); err != nil {
		return code, ""
	}
	code, err = strconv.Atoi(line[:3])
	if ftp.debug {
		log.Printf("Raw<-	<- %d \n", code)
	}
	return code, line
}

func (ftp *FTP) cmd(expects string, command string, args ...interface{}) (line string, err error) {
	if err = ftp.send(command, args...); err != nil {
		return
	}

	if line, err = ftp.receive(); err != nil {
		return
	}

	if !strings.HasPrefix(line, expects) {
		err = errors.New(line)
		return
	}

	return
}

func (ftp *FTP) Rename(from string, to string) (err error) {
	if _, err = ftp.cmd(StatusActionPending, "RNFR %s", from); err != nil {
		return
	}

	if _, err = ftp.cmd(StatusActionOK, "RNTO %s", to); err != nil {
		return
	}

	return
}

func (ftp *FTP) Mkd(path string) error {
	_, err := ftp.cmd(StatusPathCreated, "MKD %s", path)
	return err
}

func (ftp *FTP) Rmd(path string) (err error) {
	_, err = ftp.cmd(StatusActionOK, "RMD %s", path)
	return
}

func (ftp *FTP) Pwd() (path string, err error) {
	var line string
	if line, err = ftp.cmd(StatusPathCreated, "PWD"); err != nil {
		return
	}

	res := RePwdPath.FindAllStringSubmatch(line[4:], -1)

	path = res[0][1]
	return
}

func (ftp *FTP) Cwd(path string) (err error) {
	_, err = ftp.cmd(StatusActionOK, "CWD %s", path)
	return
}

func (ftp *FTP) Dele(path string) (err error) {
	if err = ftp.send("DELE %s", path); err != nil {
		return
	}

	var line string
	if line, err = ftp.receive(); err != nil {
		return
	}

	if !strings.HasPrefix(line, StatusActionOK) {
		return errors.New(line)
	}

	return
}

func (ftp *FTP) AuthTLS(config *tls.Config) error {
	if _, err := ftp.cmd("234", "AUTH TLS"); err != nil {
		return err
	}

	ftp.tlsconfig = config

	ftp.conn = tls.Client(ftp.conn, config)
	ftp.writer = bufio.NewWriter(ftp.conn)
	ftp.reader = bufio.NewReader(ftp.conn)

	if _, err := ftp.cmd(StatusOK, "PBSZ 0"); err != nil {
		return err
	}

	if _, err := ftp.cmd(StatusOK, "PROT P"); err != nil {
		return err
	}

	return nil
}

func (ftp *FTP) ReadAndDiscard() (int, error) {
	var i int
	bufferSize := ftp.reader.Buffered()
	for i = 0; i < bufferSize; i++ {
		if _, err := ftp.reader.ReadByte(); err != nil {
			return i, err
		}
	}
	return i, nil
}

func (ftp *FTP) Type(t TypeCode) error {
	_, err := ftp.cmd(StatusOK, "TYPE %s", t)
	return err
}

type TypeCode string

const (
	TypeASCII = "A"

	TypeEBCDIC = "E"

	TypeImage = "I"

	TypeLocal = "L"
)

func (ftp *FTP) receiveLine() (string, error) {
	line, err := ftp.reader.ReadString('\n')

	if ftp.debug {
		log.Printf("< %s", line)
	}

	return line, err
}

func (ftp *FTP) receive() (string, error) {
	line, err := ftp.receiveLine()

	if err != nil {
		return line, err
	}

	if (len(line) >= 4) && (line[3] == '-') {

		closingCode := line[:3] + " "
		for {
			str, err := ftp.receiveLine()
			line = line + str
			if err != nil {
				return line, err
			}
			if len(str) < 4 {
				if ftp.debug {
					log.Println("Uncorrectly terminated response")
				}
				break
			} else {
				if str[:4] == closingCode {
					break
				}
			}
		}
	}
	ftp.ReadAndDiscard()

	return line, err
}

func (ftp *FTP) receiveNoDiscard() (string, error) {
	line, err := ftp.receiveLine()

	if err != nil {
		return line, err
	}

	if (len(line) >= 4) && (line[3] == '-') {

		closingCode := line[:3] + " "
		for {
			str, err := ftp.receiveLine()
			line = line + str
			if err != nil {
				return line, err
			}
			if len(str) < 4 {
				if ftp.debug {
					log.Println("Uncorrectly terminated response")
				}
				break
			} else {
				if str[:4] == closingCode {
					break
				}
			}
		}
	}

	return line, err
}

func (ftp *FTP) send(command string, arguments ...interface{}) error {
	if ftp.debug {
		log.Printf("> %s", fmt.Sprintf(command, arguments...))
	}

	command = fmt.Sprintf(command, arguments...)
	command += "\r\n"

	if _, err := ftp.writer.WriteString(command); err != nil {
		return err
	}

	if err := ftp.writer.Flush(); err != nil {
		return err
	}

	return nil
}

func (ftp *FTP) Pasv() (port int, err error) {
	doneChan := make(chan int, 1)
	go func() {
		defer func() {
			doneChan <- 1
		}()
		var line string
		if line, err = ftp.cmd("227", "PASV"); err != nil {
			return
		}
		re := regexp.MustCompile(`\((.*)\)`)
		res := re.FindAllStringSubmatch(line, -1)
		if len(res) == 0 || len(res[0]) < 2 {
			err = errors.New("PasvBadAnswer")
			return
		}
		s := strings.Split(res[0][1], ",")
		if len(s) < 2 {
			err = errors.New("PasvBadAnswer")
			return
		}
		l1, _ := strconv.Atoi(s[len(s)-2])
		l2, _ := strconv.Atoi(s[len(s)-1])

		port = l1<<8 + l2

		return
	}()

	select {
	case _ = <-doneChan:

	case <-time.After(time.Second * 10):
		err = errors.New("PasvTimeout")
		ftp.Close()
	}

	return
}

func (ftp *FTP) newConnection(port int) (conn net.Conn, err error) {
	addr := fmt.Sprintf("%s:%d", strings.Split(ftp.addr, ":")[0], port)

	if ftp.debug {
		log.Printf("Connecting to %s\n", addr)
	}

	if conn, err = net.Dial("tcp", addr); err != nil {
		return
	}

	if ftp.tlsconfig != nil {
		conn = tls.Client(conn, ftp.tlsconfig)
	}

	return
}

func (ftp *FTP) Stor(path string, r io.Reader) (err error) {
	if err = ftp.Type(TypeImage); err != nil {
		return
	}

	var port int
	if port, err = ftp.Pasv(); err != nil {
		return
	}

	if err = ftp.send("STOR %s", path); err != nil {
		return
	}

	var pconn net.Conn
	if pconn, err = ftp.newConnection(port); err != nil {
		return
	}
	defer pconn.Close()

	var line string
	if line, err = ftp.receive(); err != nil {
		return
	}

	if !strings.HasPrefix(line, StatusFileOK) {
		err = errors.New(line)
		return
	}

	if _, err = io.Copy(pconn, r); err != nil {
		return
	}
	pconn.Close()

	if line, err = ftp.receive(); err != nil {
		return
	}

	if !strings.HasPrefix(line, StatusClosingDataConnection) {
		err = errors.New(line)
		return
	}

	return

}

func (ftp *FTP) Syst() (line string, err error) {
	if err := ftp.send("SYST"); err != nil {
		return "", err
	}
	if line, err = ftp.receive(); err != nil {
		return
	}
	if !strings.HasPrefix(line, StatusSystemType) {
		err = errors.New(line)
		return
	}

	return strings.SplitN(strings.TrimSpace(line), " ", 2)[1], nil
}

var (
	SystemTypeUnixL8    = "UNIX Type: L8"
	SystemTypeWindowsNT = "Windows_NT"
)

var reSystStatus = map[string]*regexp.Regexp{
	SystemTypeUnixL8:    regexp.MustCompile(""),
	SystemTypeWindowsNT: regexp.MustCompile(""),
}

func (ftp *FTP) Stat(path string) ([]string, error) {
	if err := ftp.send("STAT %s", path); err != nil {
		return nil, err
	}

	stat, err := ftp.receive()
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(stat, StatusFileStatus) &&
		!strings.HasPrefix(stat, StatusDirectoryStatus) &&
		!strings.HasPrefix(stat, StatusSystemStatus) {
		return nil, errors.New(stat)
	}
	if strings.HasPrefix(stat, StatusSystemStatus) {
		return strings.Split(stat, "\n"), nil
	}
	lines := []string{}
	for _, line := range strings.Split(stat, "\n") {
		if strings.HasPrefix(line, StatusFileStatus) {
			continue
		}

		lines = append(lines, strings.TrimSpace(line))

	}

	return lines, nil
}

func (ftp *FTP) Retr(path string, retrFn RetrFunc) (s string, err error) {
	if err = ftp.Type(TypeImage); err != nil {
		return
	}

	var port int
	if port, err = ftp.Pasv(); err != nil {
		return
	}

	if err = ftp.send("RETR %s", path); err != nil {
		return
	}

	var pconn net.Conn
	if pconn, err = ftp.newConnection(port); err != nil {
		return
	}
	defer pconn.Close()

	var line string
	if line, err = ftp.receiveNoDiscard(); err != nil {
		return
	}

	if !strings.HasPrefix(line, StatusFileOK) {
		err = errors.New(line)
		return
	}

	if err = retrFn(pconn); err != nil {
		return
	}

	pconn.Close()

	if line, err = ftp.receive(); err != nil {
		return
	}

	if !strings.HasPrefix(line, StatusClosingDataConnection) {
		err = errors.New(line)
		return
	}

	return
}

/*func GetFilesList(path string) (files []string, err error) {

}*/

func (ftp *FTP) List(path string) (files []string, err error) {
	if err = ftp.Type(TypeASCII); err != nil {
		return
	}

	var port int
	if port, err = ftp.Pasv(); err != nil {
		return
	}

	if err = ftp.send("MLSD %s", path); err != nil {
	}

	var pconn net.Conn
	if pconn, err = ftp.newConnection(port); err != nil {
		return
	}
	defer pconn.Close()

	var line string
	if line, err = ftp.receiveNoDiscard(); err != nil {
		return
	}

	if !strings.HasPrefix(line, StatusFileOK) {

		if err = ftp.send("LIST %s", path); err != nil {
			return
		}

		if line, err = ftp.receiveNoDiscard(); err != nil {
			return
		}

		if !strings.HasPrefix(line, StatusFileOK) {

			err = errors.New(line)
			return
		}
	}

	reader := bufio.NewReader(pconn)

	for {
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return
		}

		files = append(files, string(line))
	}

	pconn.Close()

	if line, err = ftp.receive(); err != nil {
		return
	}

	if !strings.HasPrefix(line, StatusClosingDataConnection) {
		err = errors.New(line)
		return
	}

	return
}

/*



func (ftp *FTP) SmartLogin(username string, password string) (err error) {
	var code int

	code, _ = ftp.RawCmd("NOOP")

	if code == 220 || code == 530 {

		code, _ = ftp.RawCmd("NOOP")
		if code == 530 {

			code, _ = ftp.RawCmd("USER %s", username)
			code, _ = ftp.RawCmd("NOOP")
			if code == 331 {

				code, _ = ftp.RawCmd("PASS %s", password)
				code, _ = ftp.RawCmd("PASS %s", password)
				if code == 230 {
					code, _ = ftp.RawCmd("NOOP")
					return
				}
			}
		}

	}

	return ftp.Login(username, password)
}

*/

func (ftp *FTP) Login(username string, password string) (err error) {
	if _, err = ftp.cmd("331", "USER %s", username); err != nil {
		if strings.HasPrefix(err.Error(), "230") {

			err = nil
		} else {
			return
		}
	}

	if _, err = ftp.cmd("230", "PASS %s", password); err != nil {
		return
	}

	return
}

func Connect(addr string) (*FTP, error) {
	var err error
	var conn net.Conn

	if conn, err = net.Dial("tcp", addr); err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	object := &FTP{conn: conn, addr: addr, reader: reader, writer: writer, debug: false}
	object.receive()

	return object, nil
}

func ConnectDbg(addr string) (*FTP, error) {
	var err error
	var conn net.Conn

	if conn, err = net.Dial("tcp", addr); err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	var line string

	object := &FTP{conn: conn, addr: addr, reader: reader, writer: writer, debug: true}
	line, _ = object.receive()

	log.Print(line)

	return object, nil
}

func (ftp *FTP) Size(path string) (size int, err error) {
	line, err := ftp.cmd("213", "SIZE %s", path)

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(line[4 : len(line)-2])
}
