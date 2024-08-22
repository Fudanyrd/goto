package main

import (
	"os"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func extract_file(url string) string {
	idx := strings.LastIndex(url, "/");
	if idx == -1 {
		return url;
	}
	return url[idx + 1 :];
}

func download_save(url, fn string) {
	res, err := http.Get(url);
	if err != nil {
		log.Fatal(err);
	}

	body, err := io.ReadAll(res.Body);
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", 
					res.StatusCode, body);
	}
	if err != nil {
		log.Fatal(err);
	}

	os.WriteFile(fn, body, 0666);
}

func main() {
	args := os.Args;

	/* request and download each file. */
	for i := 1; i < len(args); i++ {

		// get filename and print
		fn := extract_file(args[i]);
		fmt.Println("Downloading", fn, "...");

		// request the file specified and save the file
		download_save(args[i], fn);
	}

	/* done. */
	return;
}
