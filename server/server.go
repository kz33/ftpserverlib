package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

type Paradise struct {
	writer        *bufio.Writer
	reader        *bufio.Reader
	theConnection net.Conn
	passiveConn   *net.TCPConn
	waiter        sync.WaitGroup
	user          string
	homeDir       string
	path          string
	ip            string
	command       string
	param         string
	total         int64
	buffer        []byte
}

func NewParadise(connection net.Conn) *Paradise {
	p := Paradise{}

	p.writer = bufio.NewWriter(connection)
	p.reader = bufio.NewReader(connection)
	p.path = "/"
	p.theConnection = connection
	p.ip = connection.RemoteAddr().String()
	return &p
}

func (self *Paradise) HandleCommands() {
	fmt.Println("Got client on: ", self.ip)
	self.writeMessage(220, "Welcome to Paradise")
	for {
		line, err := self.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				continue
			}
			break
		}
		command, param := parseLine(line)
		self.command = command
		self.param = param

		//var t T
		//reflect.ValueOf(&t).MethodByName("Foo").Call([]reflect.Value{})

		if command == "USER" {
			self.handleUser()
		} else if command == "PASS" {
			self.handlePass()
		} else if command == "SYST" {
			self.handleSyst()
		} else if command == "PWD" {
			self.handlePwd()
		} else if command == "TYPE" {
			self.handleType()
		} else if command == "EPSV" || command == "PASV" {
			self.handlePassive()
		} else if command == "LIST" || command == "NLST" {
			self.handleList()
		} else if command == "QUIT" {
			self.handleQuit()
		} else if command == "CWD" {
			self.handleCwd()
		} else if command == "SIZE" {
			self.handleSize()
		} else if command == "RETR" {
			self.handleRetr()
		} else if command == "STAT" {
			self.handleStat()
		} else if command == "STOR" || command == "APPE" {
			self.handleStore()
		} else {
			self.writeMessage(550, "not allowed")
		}

		// close passive connection each time
		self.closePassiveConnection()
	}
}

func (self *Paradise) writeMessage(code int, message string) {
	line := fmt.Sprintf("%d %s\r\n", code, message)
	self.writer.WriteString(line)
	self.writer.Flush()
}

func (self *Paradise) closePassiveConnection() {
	if self.passiveConn != nil {
		self.passiveConn.Close()
	}
}

func parseLine(line string) (string, string) {
	params := strings.SplitN(strings.Trim(line, "\r\n"), " ", 2)
	if len(params) == 1 {
		return params[0], ""
	}
	return params[0], strings.TrimSpace(params[1])
}
