package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const version = "check_platypus 0.2"

func usage() {
	fmt.Fprintln(os.Stderr, version)
	fmt.Fprintf(os.Stderr, "usage: %s [args] HOSTNAME USERNAME PASSWORD\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(64)
}

func main() {
	port := 5565
	debug := false

	flag.Usage = usage
	flag.BoolVar(&debug, "debug", false, "enable debug output")
	flag.IntVar(&port, "port", 5565, "Platypus API service port")
	flag.Parse()

	if len(flag.Args()) != 3 {
		usage()
	}

	hostname := fmt.Sprintf("%s:%d", flag.Arg(0), port)
	if debug {
		fmt.Fprintf(os.Stderr, "Host: %s\n", hostname)
	}

	username := flag.Arg(1)
	if debug {
		fmt.Fprintf(os.Stderr, "Username: %s\n", username)
	}

	password := flag.Arg(2)
	if debug {
		fmt.Fprintf(os.Stderr, "Password: %s\n", password)
	}

	addr, err := net.ResolveTCPAddr("tcp", hostname)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(2)
	}

	var platcmd = PLATXML{}
	platcmd.Header = ""
	platcmd.Body.Data.Protocol = "Plat"
	platcmd.Body.Data.Object = "addusr"
	platcmd.Body.Data.Action = "Login"
	platcmd.Body.Data.Username = username
	platcmd.Body.Data.Password = password
	platcmd.Body.Data.Logintype = "staff"
	platcmd.Body.Data.Parameters.Logintype = "Staff"
	platcmd.Body.Data.Parameters.Username = username
	platcmd.Body.Data.Parameters.Password = password
	platcmd.Body.Data.Parameters.Datatype = "XML"

	xmlcmd, err := xml.Marshal(platcmd)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(2)
	}

	xmlcmd = append([]byte(xml.Header), xmlcmd...)

	prefix := []byte("Content-Length:" + strconv.Itoa(len(xmlcmd)) + "\r\n\r\n")

	rawout := append(prefix, xmlcmd...)

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(2)
	}

	_, err = conn.Write(rawout)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(2)
	}

	header, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(2)
	}

	connlen, err := strconv.Atoi(strings.TrimSpace(header[15:]))
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(2)
	}

	buf := make([]byte, connlen)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(2)
	}

	conn.Close()

	var platresp = PLATXML{}
	err = xml.Unmarshal(buf, &platresp)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(2)
	}

	fmt.Println(platresp.Body.Data.ResponseText)

	if platresp.Body.Data.Success == 1 {
		os.Exit(0)
	} else if platresp.Body.Data.Success == 0 {
		os.Exit(0)
	}

	os.Exit(2)
}
