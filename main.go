package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s HOSTNAME USERNAME PASSWORD\n", os.Args[0])
		os.Exit(1)
	}

	hostname := os.Args[1] + ":5565"
	username := os.Args[2]
	password := os.Args[3]

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
