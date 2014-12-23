package main

import (
	"./platypus"
	"flag"
	"fmt"
	"os"
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

	plat, err := platypus.New(hostname, username, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not prep connection: %s", err)
		os.Exit(1)
	}

	err = plat.Login(username, password)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
