package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	steam "github.com/kidoman/go-steam"
)

const version string = "0.2.0"

func getRconCredentials(debug bool) (string, string) {
	var addr, pass, input string
	if debug {
		addr = "###"
		pass = "###"
		return addr, pass
	}
	fmt.Print("Enter RCON Address & Port (if not 27015): ")
	fmt.Scanln(&input)
	if input != "" && strings.Count(input, ":") <= 1 {
		chk, err := regexp.MatchString("^[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}:[0-9]{1,5}$", input)
		if err != nil {
			fmt.Println(err)
		}
		if strings.ContainsAny(input, ":") && chk {
			addr = input
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

func getRconCommand() string {
	var cmd string
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}

	if string([]rune(input)[0]) == "!" {
		switch string([]rune(input)) {
		case "!exit":
			// placeholder
			os.Exit(0)
		default:
			fmt.Println("Intercepted command: " + input)
			cmd = ""
		}
	} else {
		cmd = input
	}

	return cmd
}

func checkResponse(response string) string {
	return response
}

func main() {
	fmt.Println("\nSource Dedicated Server RCON Interface (" + version + ")\n")
	debug := flag.Bool("debug", false, "debug")
	flag.Parse()
	if *debug {
		steam.SetLog(log.New())
	}

	addr, pass := getRconCredentials(*debug)
	if addr == "" {
		fmt.Println("What am I supposed to do with an empty string?")
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
			cmd := getRconCommand()
			if cmd == "" {
				continue
			}
			if *debug {
				fmt.Println("[Debug] User Input: " + cmd)
			}
			response, err := sendRconCommand(addr, pass, cmd)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(checkResponse(response))
		}
	}
}
