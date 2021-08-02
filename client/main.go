package main

import(
	"os"
	"strconv"
	"strings"
	"fmt"
	"net"
	"log"
	"bufio"
) 

var setts map[string]string =  getsettings("settings.txt")

//removing an element from a slice
func remove(slice []string, s int) []string {
    return append(slice[:s], slice[s+1:]...)
}


//a func to check if a pseudo is right or not
func checknick(nick string) {
	nick = nick[:(len(nick)-2)]
	authorized := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
	if len(nick) <= 4 || len(nick) >= 16 {
		log.Fatal("Your nickname might contain between 4 and 16 characters. To get help about how the client works, type client.exe /help\n")
	}
	for i:=0;i<(len(nick));i++ {
		if strings.Contains(authorized,string(nick[i])) == false {
			log.Fatal("Something happens with '"+string(nick[i])+"'\nYour nickname might contain only characters "+authorized)
		}
	}
}


//a func to check if a string is a valid IP or not
func checkip(ip string) {
	if net.ParseIP(ip) == nil {
		log.Fatal("The ip you entered is not valid. Use the IPV4 ip address only. To get help about how the client works, type client.exe /help\n")
	}
}

//a func to check if the port is valid or not
func checkport(port string) {
	i, err := strconv.Atoi(port)
	if err != nil {
		fmt.Printf("An error has been raised checking the port. The port might be a number between 1 and 65 534 \nTo get help about how the client works, type client.exe /help\n %s",err)
	}
	if i < 0 || i > 65534 {
		fmt.Printf("An error has been raised checking the port. The port might be a number between 1 and 65 534 \nTo get help about how the client works, type client.exe /help\n%s",err)
	}
}


//get the settings in the settings.txt file
func getsettings(path string)map[string]string{
	params := make(map[string]string)
	file, err:= os.Open(path)
	if err != nil {
		log.Fatal("Failed opening the "+path+" file :\n%s",err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if string(scanner.Text()[0]) != "#" {
			params[strings.Split(scanner.Text(),":")[0]] = strings.Split(scanner.Text(),":")[1]
		}
	}
	return(params)
}

//a basic func to check if a slice contains or not the value we are looking for 
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func helpfunc() {
	log.Fatal("Welcome to the client of olc tcp chatroom.\nTo use this client is pretty easy just type :\n'client.exe <IP> <PORT>'\nWhere IP is an IPV4 adress and PORT is a valid port. Olc servers usually run on port 9000, at least it is the default port.\nOnce you are connected to a room, a nickname like member_12 is given to you.\nTo change your nickanme in the room just type (once you are connected to it) :\n/nick <nickname>\nWhere nickname is between 4 and 16 characters and only contains 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_'\nIf you want to nick automatically just type :\n/nick\nAnd the room will give you your default nickname (which one you can modify in the settings.txt). By default, your default nickname might be John_Galt.\nTo get some help type client.exe help.\nFor more information feel free to contact me. ")
}

//func to get a connection with the server
func getconn(IP,port string) net.Conn {
	conn, err := net.Dial("tcp",IP+":"+port)
	if err != nil {
		fmt.Println("An error has been raised while trying to reach the server")
		log.Fatal(err)
	}
	return(conn)
}

func displaymsg(c chan string) {
	for {
		fmt.Print(<-c)
	}
}

func getmsg(conn net.Conn, c chan string) {
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error trying to read data from server")
			log.Fatal(err)
		}
		c <- data
	}
}

func getaddr() (string,string) {
	if contains(os.Args, "help") {
		helpfunc()
		return "",""
	} else {
		return os.Args[1],os.Args[2]
	}
}

func getinpt(c chan string) {
	var text string
	for {
		scanner := bufio.NewReader(os.Stdin)
		text, _ = scanner.ReadString('\n')
		if len(text) >= 5 {
			if text[:5] == "/nick" {
				if len(text) > 7 {
					fmt.Printf("Checking for nickname "+text[6:(len(text)-2)]+"...")
					checknick(text[6:])
					fmt.Printf(" Valid name\n")
				} else {
					fmt.Printf("Checking for nickname "+setts["defaultnick"]+"...")
					text = "/nick " + setts["defaultnick"]
					fmt.Printf(" Valid name\n")
				}
			}
		}
		c <- text
	}
}


func main() {
	comming := make(chan string)
	outing := make(chan string)
	fmt.Println("chat has been successfully launched !\nYour default nickname is "+setts["defaultnick"])
	IP, port := getaddr()
	checkip(IP)
	checkport(port)
	conn := getconn(IP,port)
	go getmsg(conn, comming)
	go displaymsg(comming)
	go getinpt(outing)
	var msg string
	for {
		msg = <- outing
		_, err := conn.Write([]byte(msg+"    \n"))
		if err != nil {
			log.Fatal("Cannot reach the server ("+conn.RemoteAddr().String()+")")
		}
	}
}