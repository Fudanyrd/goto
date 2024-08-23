package main
/* Display current directory(absolute path) and leave. */

import (
	"fmt"
	"os"
	"log"
)

/** Use os.Getwd() to fetch pwd's output. See 
  <https://cs.opensource.google/go/go/+/go1.23.0:src/os/getwd.go;l=22> */
func main() {
	var dir string;
	var err error;
	dir, err = os.Getwd();
	
	if err != nil {
		log.Fatal(err);
	}

	fmt.Println(dir);
	return;
}
