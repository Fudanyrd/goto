package main
/** touch create one or more files. */

import (
	"os"
	"log"
)

func touch(filename string) {
	/** will need to use os.Create See
	<https://cs.opensource.google/go/go/+/go1.23.0:src/os/file.go;l=373>
	*/
	_, err := os.Create(filename);
	if err != nil {
		log.Fatal(err);
	}
	return;
}

func main() {
	var args []string = os.Args;

	for i := 1; i < len(args); i++ {
		touch(args[i]);
	}

	return;
}
