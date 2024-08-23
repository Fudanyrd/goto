package main
// citation: <https://go.dev/doc/articles/wiki/>

// A simple file server that answers request for files 
// in the directories where the server program resides.

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "io"
)

func handler(w http.ResponseWriter, r *http.Request) {
    filename := r.URL.Path[1:];
    if len(filename) == 0 {
        // raise 404 Not Found.
        http.Error(w, "empty file name", http.StatusNotFound);
        return;
    }

    fobj, err := os.Open(filename);
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound);
        return;
    }

	// is fobj a directory?
	info, err := fobj.Stat();
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError);
		return;
	}
	if info.IsDir() {
        http.Error(w, "requested file is a directory", http.StatusInternalServerError);
		return;
	}

    // make buffer and read the file
	buf := make([]byte, 1024);
	for {
		count, err := fobj.Read(buf);

		if count > 0 {
			fmt.Fprintf(w, "%s", buf[:count]);
		}
		if err != nil {
			if err != io.EOF {
				log.Fatal(err);
			}
			break;
		}
	}
    return;
}

func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
