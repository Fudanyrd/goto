package main

import (
	"fmt"
	"os"
	"log"
	"io"
)

var hmap map[int32]string = map[int32]string{
	0: "0",
	1: "1",
	2: "2",
	3: "3",
	4: "4",
	5: "5",
	6: "6",
	7: "7",
	8: "8",
	9: "9",
	10: "a",
	11: "b",
	12: "c",
	13: "d",
	14: "e",
	15: "f",
};
var specials []byte = []byte("~`!@#$%^&*()_+-=[]{}\\|:;'\"< >,/");
const zero byte = byte('0');
const nine byte = byte('9');
const a byte = byte('a');
const A byte = byte('A');
const z byte = byte('z');
const Z byte = byte('Z');

func print_hex(v int32) {
	var arr [8]int32;
	var i int32;

	fmt.Printf("0x");
	for i = 7; i >= 0; i -= 1 {
		arr[i] = v & 0xf;
		v = v >> 4;
	}

	for i = 0; i < 8; i++ {
		str, _ := hmap[arr[i]];
		fmt.Printf("%s", str);
	}
}

func print_hex_last(v int32) {
	str, _ := hmap[(v & 0xf0) >> 4];
	fmt.Printf(" %s", str);
	str, _ = hmap[v & 0x0f];
	fmt.Printf("%s", str);
}

func print_byte(b byte) byte {
	if b <= z && b >= a {
		return b;
	}
	if b <= Z && b >= A {
		return b;
	}
	if b <= nine && b >= zero {
		return b;
	}

	for i:=0; i < len(specials); i++ {
		if b == specials[i] {
			return b;
		}
	}
	return byte('.');
}

func main() {
	args := os.Args;
	if (len(args) != 2) {
		fmt.Printf("usage hexdump <filename>\n");
		return;
	}

	var pos int32 = 0;

	/* Open targeted file */
	file, err := os.Open(args[1]);
	if err != nil {
		log.Fatal(err);
		return;
	}

	buf := make([]byte, 16);

	for {
		count, err := file.Read(buf);
		if count > 0 {
			// start position
			print_hex(pos);
			pos += int32(count);

			// hex repr
			fmt.Printf(" |");
			for i := 0; i < count; i++ {
				print_hex_last(int32(buf[i]));
				buf[i] = print_byte(buf[i]);
				if i == 7 {
				  	fmt.Printf(" |");
				}
			}
			fmt.Printf(" | ");

			// char repr
			fmt.Printf(" | ");
			fmt.Printf("%s", buf[:count]);
			fmt.Printf(" | ");

			// clear
			fmt.Printf("\n");
		}

		if err != nil {
			if err != io.EOF {
				log.Fatal(err);
			}
			break;
		}
	}
}
