package main
/** ls: list all files in the current directory. */

import (
	"fmt"
	"os"
	"log"
)

func ndigits(num int64) int {
	var ret int = 0;
	for ; num > int64(0); num = num / int64(10) {
		ret += 1;
	}

	return ret;
}

func ls(dir string) {
	dobj, err := os.Open(dir);
	if err != nil {
		log.Fatal(err);
		return;
	}

	// is dobj a directory?
	info, err := dobj.Stat();
	if err != nil {
		log.Fatal(err);
		return;
	}
	if !info.IsDir() {
		fmt.Println("not a directory");
		return;
	}

	/* read all entries */
	entries, err := dobj.Readdir(0);

	/* find out the maximum length of size */
	var maxDigit int = 0;
	var digits []int = make([]int, len(entries));

	for id, entry := range entries {
		digits[id] = ndigits(entry.Size());
		if digits[id] > maxDigit {
			maxDigit = digits[id];
		}
	}

	/** print each entry */
	for id, entry := range entries {
		fmt.Print(entry.Mode());

		// append whitespaces
		fmt.Print(" ");
		for i := 0; i < (maxDigit - digits[id]); i++ {
			fmt.Print(" ");
		}
		// print size
		fmt.Printf("%v", entry.Size());

		// print name
		fmt.Print(" ");
		fmt.Println(entry.Name());
	}

	return;
}

func main() {
	args := os.Args;
	path := "";

	if len(args) >= 2 {
		path = args[1];
	} else {
		path = ".";
	}

	ls(path);
	return;
}
