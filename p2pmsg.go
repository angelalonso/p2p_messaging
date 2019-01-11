package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type message struct {
	UserName string `json:"username"`
	Text     string `json:"text"`
}

type user struct {
	UserName string
	IpPort   string
}

func main() {
	self := user{}
	remote := []user{}
	// check parameters received
	if len(os.Args) < 3 {
		help()
		os.Exit(2)
	} else {
		self.UserName = os.Args[1]
		self.IpPort = "127.0.0.1:" + os.Args[2]
		// TODO: take more than one remote address
		if len(os.Args) > 3 {
			remote = append(remote, user{
				IpPort: os.Args[3],
			})
		} else {
			remote = append(remote, user{})
		}
	}
	certSender, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Fatalf("client: loadkeys: %s", err)
	}
	configSender := tls.Config{Certificates: []tls.Certificate{certSender}, InsecureSkipVerify: true}

	certListener, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	configListener := tls.Config{Certificates: []tls.Certificate{certListener}}

	var wg sync.WaitGroup
	m := make(chan message)
	c := make(chan bool)

	fmt.Println(len(remote))

	fmt.Printf("Welcome, %s\nREMOTE: %s\n", self.UserName, remote[0].IpPort)

	go listener(strings.Split(self.IpPort, `:`)[1], &configListener, &wg, c)
	go sender(remote, &configSender, m, c)
	go read(m, self.UserName)
	wg.Add(3)

	wg.Wait()

	// create list of users, including oneself

	// send one's data to other user

	// receive other users' data and add to our list

	// send message to all in list

	// receive message

}

func listener(localPort string, config *tls.Config, wg *sync.WaitGroup, c chan bool) {
	defer wg.Done()
	ln, err := tls.Listen("tcp", ":"+localPort, config)
	if err != nil {
		log.Println("Starting listener error!", err)
		return
	}
	defer ln.Close()
	for {
		fmt.Print("\n")
		log.Println("Listening...")
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		defer conn.Close()
		var msg = message{}
		for {
			err := json.NewDecoder(conn).Decode(&msg)

			if err != nil {
				c <- false
				log.Println("Disconnected!")
				break
			}
			if msg.Text != "" {
				log.Printf("%s: %s", msg.UserName, msg.Text)
				fmt.Print(">>")
			}
		}
	}
}

func sender(remote []user, config *tls.Config, m chan message, c chan bool) {
CONNECTION:
	for {
		conn, err := tls.Dial("tcp", remote[0].IpPort, config)
		if err != nil {
			log.Println("Connecting...")
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println("Connected!")
		fmt.Print(">>")

		for {
			select {
			case msg := <-m:
				err := json.NewEncoder(conn).Encode(msg)
				if err != nil {
					continue CONNECTION
				}

			case alive := <-c:
				if !alive {
					continue CONNECTION
				}
			}
		}
	}
}

func read(m chan message, userName string) {
	var reader = bufio.NewReader(os.Stdin)
	var msg = message{UserName: userName}
	for {
		fmt.Print(">>")
		if text, _ := reader.ReadString('\n'); text != "\n" {

			msg.Text = strings.TrimSpace(text)
			m <- msg
		}
	}
}

func help() {
	fmt.Println(`ERROR:
You need at least two parameters

SYNTAX:
p2pmsg <your user name> <your Port> <remote IP:remotePort>`)
}
