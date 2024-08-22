package main

import (
	"fmt"
	"bufio"
	"io"
	"os"
)

type MyStdin struct{
	fobj *os.File;
}

func (reader MyStdin) Read(p []byte) (n int, err error) {
	return reader.fobj.Read(p);
}

var rd io.Reader = MyStdin{ os.Stdin };
var StdinReader *bufio.Reader = bufio.NewReader(rd);

func main() {
	for {
		line, err := StdinReader.ReadString('\n');
		if len(line) > 0 {
			fmt.Printf(line);
		}
		if err != nil {
			break;
		}
	}
}
