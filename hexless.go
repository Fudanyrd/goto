package main

/**
 * hexless: combination of "hexdump" and "less" program. 
 * Usage: hexless <filename>, and keep press <enter> until EOF is met.
 */

import (
	"fmt"
	"log"
	"os"
	/*
	"strings"
	"bufio"
	"io"
	*/
)

/*-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 *                            Constants
 -+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+*/

const (
	// width and height of screen
	width int = 16
	height int = 16
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
const Down string = "\r\n";
const Up string = "\033[A";

/*-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 *                          Data Structures
 -+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+*/

type Cursor struct {
	row int;
	col int;
}

var cursor Cursor = Cursor{ 0, 0 };

func (cr *Cursor) print_hex(v int32) {
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
	cr.col += 10;  // 0x000 00000
}

func (cr *Cursor) print_hex_last(v int32) {
	str, _ := hmap[(v & 0xf0) >> 4];
	fmt.Printf(" %s", str);
	str, _ = hmap[v & 0x0f];
	fmt.Printf("%s", str);
	cr.col += 3;  // " 0a"
}

func (cr *Cursor) print_byte(b byte) {
	cr.col += 1;
	var str []byte = make([]byte, 1);
	if b <= z && b >= a {
		str[0] = b;
		fmt.Printf("%s", str);
		return;
	}
	if b <= Z && b >= A {
		str[0] = b;
		fmt.Printf("%s", str);
		return;
	}
	if b <= nine && b >= zero {
		str[0] = b;
		fmt.Printf("%s", str);
		return;
	}

	for i:=0; i < len(specials); i++ {
		if b == specials[i] {
			str[0] = b;
			fmt.Printf("%s", str);
			return;
		}
	}
	fmt.Printf(".");
	return;
}

// s may not contain control characters.
func (cr *Cursor) puts(s string) {
	fmt.Printf(s);
	cr.col += len(s);
}

func (cr *Cursor) down() {
	fmt.Printf(Down);
	cr.row += 1;
	cr.col = 0;
}
func (cr *Cursor) up() {
	if cr.row == 0 {
		return;
	}
	fmt.Printf(Up);
	cr.row -= 1;
	cr.col = 0;
}
func (cr *Cursor) backspace() {
	if cr.col == 0 {
		return;
	}
	fmt.Printf("\b");
	cr.col -= 1;
}

// clear the screen and reset the cursor.
func (cr *Cursor) clear() {
	for r := 0; r < cr.row; r++ {
		fmt.Printf(Up);
	}
	cr.row = 0;
	cr.col = 0;
}

func whiteout() {
	cursor.clear();
	for r := 0; r < 17; r++ {
		for c := 0; c < 88; c++ {
			cursor.puts(" ");
		}
		cursor.down();
	}
	cursor.clear();
}

// white out the screen

/* A 'row' in the file(containing 16 bytes) */
type Frow struct {
	prev *Frow;
	next *Frow;
	row  []byte;
	size int;
};

func (fr *Frow) init() {
	fr.row = make([]byte, width);
}

func (fr *Frow) display(offset int) {
	cursor.print_hex(int32(offset));
	cursor.puts(" |");

	// hex segment
	for i := 0; i < fr.size; i++ {
		cursor.print_hex_last(int32(fr.row[i]));
		if i == 7 {
		  	cursor.puts(" |");
		}
	}
	cursor.puts(" | ");

	// char segment
	for i := 0; i < fr.size; i++ {
		cursor.print_byte(fr.row[i]);
	}
	cursor.puts(" |");

	// end
	cursor.down();
	return;
}

/* doubly-linked list repr of a file. */
type Flist struct {
	head Frow;
	tail Frow;
	fobj *os.File;
	row  *Frow;   // current cursor row
	col   int;    // current cursor column
};

func (fl *Flist) push_back(fr *Frow) {
	last := fl.tail.prev;
	fr.prev = last;
	fr.next = &(fl.tail);
	last.next = fr;
	fl.tail.prev = fr;
}

func (fl *Flist) init(fn string) bool {
	// list init
	fl.head.next = &(fl.tail);
	fl.tail.prev = &(fl.head);

	// file init
	file, err := os.Open(fn);
	if err != nil {
		log.Fatal(err);
		return false;
	}
	fl.fobj = file;

	// check directory
	info, err := file.Stat();
	if err != nil || info.IsDir() {
		return false;
	}

	// read file
	for {

		// make and initialize a row.
		tmp := make([]Frow, 1);
		fr := &tmp[0];
		fr.init();
		count, err := file.Read(fr.row);
		fr.size = count;
		fl.push_back(fr);

		if err != nil {
			break;
		}
	}

	fl.row = fl.head.next;

	return true;
}

func (fl *Flist) next() bool {
	if fl.row == &(fl.tail) {
		return false;
	}
	fl.row = fl.row.next;
	return true;
}

func (fl *Flist) display() {
	cursor.clear();
	it := fl.row;
	if it == &fl.tail {
		return;
	}
	pos := 0;
	for i := 0; i < height; i++ {
		it.display(pos);
		pos += it.size;
		it = it.next;
		if it == &fl.tail {
			break;
		}
	}
}

func main() {
	args := os.Args;
	if len(args) != 2 {
		fmt.Printf("usage hexless <filename>\n");
		return;
	}

	fl := Flist{};
	if !fl.init(args[1]) {
		fmt.Printf("fail\n");
		return;
	}

	cursor.row = 0;
	fl.display();
	b := make([]byte, 2);
	for {
		_, err := os.Stdin.Read(b);
		if err != nil {
			break;
		}	
		// erase the enter key pressed
		fmt.Printf(Up);

		ret := fl.next();
		if !ret {
			break;
		}
		fl.display();
	}
}
