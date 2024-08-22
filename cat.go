package main

import (
	"fmt"
	"os"
	"log"
	"io"
)

func main() {
	args := os.Args;
	if (len(args) != 2) {
		fmt.Printf("usage cat <filename>\n");
		return;
	}

	file, err := os.Open(args[1]);
	if err != nil {
		log.Fatal(err);
		return;
	}

	buf := make([]byte, 256);
	for {
		count, err := file.Read(buf);

		if count > 0 {
			fmt.Printf("%s", buf[:count]);
		}
		if err != nil {
			if err != io.EOF {
				log.Fatal(err);
			}
			break;
		}
	}

	fmt.Println("");
	return;
}
