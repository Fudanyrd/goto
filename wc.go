package main
/** wc: word count, prints bytes, lines, words of a file. */

import (
	"os"
	"strings"
	"fmt"
	"io"
	"log"
	"bufio"
)

type LineReader struct {
	fobj *os.File;
}
func (reader LineReader) Read(p []byte) (int, error) {
	return reader.fobj.Read(p);
}

func main() {
	args := os.Args;
	if len(args) != 2 {
		fmt.Println("usage wc <filename>");
	}

	// open the file
	fobj, err := os.Open(args[1]);
	if err != nil {
		log.Fatal(err);
		return;
	}

	b := 0;  // bytes
	l := 0;  // lines
	w := 0;  // words

	// create file reader.
	rd := LineReader{ fobj };
	var reader *bufio.Reader = bufio.NewReader(rd);

	for {
		line, err := reader.ReadString('\n');
		if len(line) > 2 {
			l += 1;
			b += len(line) - 2;
			w += len(strings.Fields(line));
		}
		if err != nil {
			if err != io.EOF {
				log.Fatal(err);
			}
			break;
		}
	}

	fmt.Println("byte:", b, "line:", l, "word:", w);
	return;
}
