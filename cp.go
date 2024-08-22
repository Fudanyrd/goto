package main

import (
	"os"
	"fmt"
	"io"
	"log"
)

func main() {
	args := os.Args;
	if len(args) != 3 {
		fmt.Println("usage cp <src> <dest>");
		return;
	}

	buf := make([]byte, 1024);	

	// opening source file
	file, err := os.Open(args[1]);
	if err != nil {
		log.Fatal(err);
		return;
	}

	for {
		count, err := file.Read(buf);
		if count > 0 {
			os.WriteFile(args[2], buf[:count], 0666);
		}

		if err != nil {
			if err != io.EOF {
				log.Fatal(err);
			}

			break;
		}
	}
}
