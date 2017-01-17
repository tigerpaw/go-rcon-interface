package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	steam "github.com/kidoman/go-steam"
)

const version string = "0.1.1"

func getRconCredentials(debug bool) (string, string) {
	var addr, pass string
	fmt.Print("Enter RCON Address & Port (if not 27015): ")
	fmt.Scanln(&addr)
	if addr != "" && strings.Count(addr, ":") <= 1 {
		if strings.ContainsAny(addr, ":") {
			input := strings.Split(addr, ":")
			fmt.Println(input[0] + ":" + input[1])
		} else {
			addr += ":27015"
		}
	} else {
		fmt.Println("Incorrect address format, Example: 1.3.3.7:27017")
	}
	fmt.Print("Enter RCON Password: ")
	fmt.Scanln(&pass)

	return addr, pass
}

func sendRconCommand(addr string, pass string, cmd string) (string, error) {
	opt := &steam.ConnectOptions{RCONPassword: pass}
	rcon, err := steam.Connect(addr, opt)
	if err != nil {
		return "", err
	} else if cmd == "test" {
		rcon.Close()
		return "good", nil
	}

	defer rcon.Close()
	resp, err := rcon.Send(string(cmd))
	if err != nil {
		return "", err
	}

	return resp, nil
}

// add function for processing input & intercepting commands

func main() {
	fmt.Println("Source Dedicated Server RCON Interface (" + version + ")")
	debug := flag.Bool("debug", false, "debug")
	flag.Parse()
	if *debug {
		steam.SetLog(log.New())
	}

	addr, pass := getRconCredentials(*debug)
	if addr == "" {
		fmt.Println("What am I supposed to do with empty strings?")
		return
	}

	response, err := sendRconCommand(addr, pass, "test")
	if err != nil {
		fmt.Println(err)
		fmt.Println(response)
		return
	} else if response == "good" {
		for {
			fmt.Printf("[%s]RCON> ", addr)
			var cmd string
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
			}
			// sanitize input
			cmd = input
			if *debug {
				fmt.Println("[Debug] Command: " + cmd)
			}
			response, err := sendRconCommand(addr, pass, cmd)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(response)
		}
	}
}
