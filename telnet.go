package main

/** telnet example usage:
 * 1. have the server started;
 * 2. execute telnet localhost 8080;
 * 3. type the following request header:
```
GET /telnet.go HTTP/1.1
Host: localhost

```
 */

import (
	"log"
	"os"
	"fmt"
	"net"
	"bufio"
	"io"
)

type MyStdin struct{
	fobj *os.File;
}

func (reader MyStdin) Read(p []byte) (n int, err error) {
	return reader.fobj.Read(p);
}

var rd io.Reader = MyStdin{ os.Stdin };
var StdinReader *bufio.Reader = bufio.NewReader(rd);

func receive(conn net.Conn) int {
	reader := bufio.NewReader(conn);
	buf := make([]byte, 256);
	
	rb := 0;
	for {
		count, err := reader.Read(buf);

		if count > 0 {
			fmt.Printf("%s", buf[:count]);
			rb += count;
		}
		if count < len(buf) || err != nil {
			break;
		}
	}

	fmt.Print("\n");
	fmt.Println("connection closed by foreign host.");
	return rb;
}

func connect(ip string) bool {
	conn, err := net.Dial("tcp", ip);
	if err != nil {
		return false;
	}

	for {
		line, err := StdinReader.ReadString('\n');
		if len(line) > 2 {
			fmt.Fprintf(conn, line);
		} else {
			fmt.Fprintf(conn, line);
			receive(conn);
			break;
		}
		if err != nil {
			break;
		}
	}

	return true;
}

func main() {
	args := os.Args;
	if len(args) != 3 {
		fmt.Printf("Usage telnet <host-name> <host-port>\n");
		return;
	}

	// nslookup
	host := args[1];
	ips, err := net.LookupHost(host);
	if err != nil {
		log.Fatal(err);
		return;
	}

	port := args[2];
	for _, ip := range ips {
		str := ip;
		fmt.Println("Trying", str);
		if connect(str + ":" + port) {
			return;
		}
	}

	fmt.Println("Sorry, all connects failed");
	return;
}
