package main

import (
	"./platypus"
	"flag"
	"fmt"
	"os"
	"time"
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
	lastEventRun := 0

	flag.Usage = usage
	flag.BoolVar(&debug, "debug", false, "enable debug output")
	flag.IntVar(&port, "port", 5565, "Platypus API service port")
	flag.IntVar(&lastEventRun, "lasteventrun", 0, "check if event scheduler was triggered in past N minutes")
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

	if debug {
		plat.Debug = true
	}

	if lastEventRun > 0 {
		lastRun, err := plat.LastRun()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		var deltaTime time.Duration = time.Duration(lastEventRun) * time.Minute
		deadline := lastRun.Add(deltaTime)

		if time.Now().Unix() > deadline.Unix() {
			fmt.Printf("Deadline of %s passed\n", deadline)
			os.Exit(2)
		} else {
			os.Exit(0)
		}

	} else {
		err = plat.Login(username, password)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}

}
