package main
import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	res, err := http.Get(os.Args[1]);
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
	fmt.Printf("%s", body);
}
