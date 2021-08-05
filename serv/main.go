package main

import(
	"bufio"
	"strings"
	"fmt"
	"time"
	"os"
	"log"
	"strconv"
	"net"
)

type member struct {
	conn net.Conn
	nick string
}

type message struct{
	author member
	text string
	timestamp string
}

var setts map[string]string = getsettings("settings.txt")
var mlist []*member

func getln(port string) net.Listener{
	ln, err := net.Listen("tcp",":"+port)
	if err != nil {
		log.Fatal("Error starting the listening process : \n"+err.Error())
	}
	return(ln)
}

func getconn(ln net.Listener, bl []string) net.Conn{
	conn, err := ln.Accept()
	if err != nil {
		fmt.Printf("Error trying to accept a connection, listening will continue : \n"+err.Error()+"\n")
	}
	for _, ip := range bl {
		if ip == strings.Split(conn.RemoteAddr().String(),":")[0] {
			fmt.Println("This IP has been blacklisted : "+strings.Split(conn.RemoteAddr().String(),":")[0])
			conn.Close()
			return(nil)
		}
	} 
	return(conn)
}

func getdir() string {
	dir, err := os.Getwd() 
	if err != nil {
		fmt.Printf("Error trying to import directory : \n"+err.Error())
	}
	return(dir)
}

func getbl(path string) []string {
	var bl []string 
	file, err:= os.Open(path)
	if err != nil {
		log.Fatal("Failed opening the settings file :\n%s",err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if string(scanner.Text()[0]) != "#" {
			bl = append(bl, scanner.Text())		
		}
	}
	return bl
}

func getsettings(path string)map[string]string{
	params := make(map[string]string)
	file, err:= os.Open(path)
	if err != nil {
		log.Fatal("Failed opening the settings file :\n%s",err)
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


func atoi(s string) int{
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("Error trying to convert the number : "+s)
	}
	return(i)
}

func addmember(m member, ml *[]*member) {
	*ml = append(*ml, &m)
}

func getfrom(c chan message, conn net.Conn,dic map[net.Conn]*member) {
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("Error while trying to listen data from "+dic[conn].nick + " ("+conn.RemoteAddr().String()+") :\n"+err.Error()+"\n")
			delconn(dic, &mlist, conn)
			return
		}
		c <- message{author : *(dic[conn]), text : data, timestamp : time.Now().Format("03:04:05")}
	}	
}

func sendmember(m  member, msg message, dic map[net.Conn]*member) {
	_, err := m.conn.Write([]byte(msg.timestamp+"@"+msg.author.nick+" : "+msg.text))
	if err != nil {
		fmt.Println("Unreachable member : ",m.conn.RemoteAddr().String())
		delconn(dic, &mlist, m.conn)
	}
}

func sendlist(dic map[net.Conn]*member, recv member) {
	var string_mlist string = "Here is the list of members (max "+setts["MAX_MEMBER"]+")"+" :"
	for _, m := range mlist {
		if (*m).nick == recv.nick {
			string_mlist += "["+(*m).nick+" (you)]"
		} else {
			string_mlist += "["+(*m).nick+"]"
		}
	}
	fmt.Println(string_mlist)
	sendmember(recv, message{author : member{}, text : string_mlist+"\n", timestamp : time.Now().Format("03:04:05")},dic)
}


func delconn(dic map[net.Conn]*member, l *[]*member, conn net.Conn) {
	delete(dic, conn)
	for i:=0;i<len(*l);i++ {
		if (*(*l)[i]).conn == conn {
			for _, recv := range mlist {
				if conn != recv.conn {
					sendmember(*recv,message{author : member{}, text : "User "+(*(*l)[i]).nick+" has left the room\n", timestamp : time.Now().Format("03:04:05")}, dic)
				}
			}
   			fmt.Println("User "+(*(*l)[i]).nick + " has been deleted")
			(*(*l)[i]) = (*(*l)[len(*l)-1])
    		*l = (*l)[:len(*l)-1]
    		break
		}
	}
}

func newchatter(msg message, m member, ml *[]*member, dic map[net.Conn]*member) {
	for _, member := range *ml {
		if *member != m {
			sendmember(*member,msg,dic)
		}
	}
}


func handlemsg(c chan message, dic map[net.Conn]*member) {
	for {
		msg := <- c
		fmt.Println(msg.text,len(msg.text))
		t := msg.text[:len(msg.text)-2]+"     "
		switch t[:5] {
		case "/list" :
			sendlist(dic,msg.author)
		case "/nick" :
			changenick(&(msg.author), strings.ReplaceAll(t[6:]," ",""), &dic)
		default :
			for _, m := range mlist {
				if *m != msg.author {
					sendmember(*m, msg, dic)
				}
			}
		}
	}
}

func ispseudalready(name string) bool {
	for _, iter := range mlist {
		if (*iter).nick == name {
			return true
		}
	}
	return false
}

func changenick(m *member, name string, dic *map[net.Conn]*member) {
	for i, iter := range mlist {
		if *m == *iter {
			var adding string = ""
			if ispseudalready(name) {
				adding = "bis"
				fmt.Println("User "+(*m).nick+" tried to change his nick for "+name+" but it was already taken, adding a 'bis'")
				sendmember(*m,message{author : member{}, text : "Someone here already has this nickname so we added a 'bis' after. Type /list to see all nicknames in the room\n", timestamp : time.Now().Format("03:04:05")}, *dic)
			}
			fmt.Println("adding is ",adding)
			old := (*m).nick
			name = strings.ReplaceAll(name,"\n","")+adding
			fmt.Printf("User "+old+" changed his nick for ")
			mlist[i] = &member{conn : (*m).conn, nick : name}
			(*dic)[(*m).conn] = &member{conn : (*m).conn, nick : name}
			new := (*dic)[(*m).conn].nick
			fmt.Printf(new+"\n")
			for _, recv := range mlist {
				if *recv != *m {
					sendmember(*recv,message{author : member{}, text : old+" changed his nickname for "+new+"\n", timestamp : time.Now().Format("03:04:05")}, *dic)
				}
			}
		}		  
	}
}

func main() {
	bl := getbl("blacklist.txt")
	commchan := make(chan message)
	conndic := make(map[net.Conn]*member)
	fmt.Println("[*] Server started : "+time.Now().Format("02/01 03:04:05"))
	fmt.Println("________________________________\nSettings successfully imported from "+getdir()+" in settings.txt")
	ln := getln(setts["PORT"])
	fmt.Printf("________________________________\nStarted Listening on port "+setts["PORT"]+"...\n")
	go handlemsg(commchan, conndic)
	var i int 
	for len(mlist) <= atoi(setts["MAX_MEMBER"]) {
		conn := getconn(ln, bl)
		if conn != nil {
			addmember(member{conn:conn,nick:"member_"+strconv.Itoa(i)},&mlist)
			i++
			conndic[conn] = *(&mlist[len(mlist)-1])
			fmt.Printf("New user joined the room ("+conn.RemoteAddr().String()+")\n")
			go newchatter(message{author : member{}, text : "SERVER MESSAGE : a new member has joined the chat : "+mlist[len(mlist)-1].nick+"\n", timestamp : time.Now().Format("03:04:05")},*(mlist[len(mlist)-1]),&mlist,conndic)
			go getfrom(commchan, conn, conndic)
		}
	}
}