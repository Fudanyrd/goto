package main

/**
 * bide: my tiny "binary IDE", support the following vim commands:
 * "hjkl" move the cursor; "w" save changes; "i" insert a byte before cursor;
 * "x" delete a byte; "s" replace current byte.
 */

import (
	"fmt"
	"log"
	"os"
	"strings"
	"bufio"
	"io"
)

const o rune = rune('0');
const nin rune = rune('9');

func atoi(word string) int32 {
	var first rune = []rune(word)[0];
	if first < o || first > nin {
		// not a number, return the ascii of first character
		return int32(first);
	}

	var ret int32 = int32(0);
	buf := []rune(word);

	for i := 0; i < len(buf); i++ {
		ret = ret * int32(10);
		ret += int32(buf[i] - rune('0'));
	}
	return ret;
}

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
const Cs string = "\033[01;31m";
const Ce string = "\033[0m";

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

func (fr *Frow) removeAll() {
	p, n := fr.prev, fr.next;
	p.next = n;
	n.prev = p;
}

func (fr *Frow) removeAt(off int) {
	if off < 0 || off >= fr.size {
		return;
	}

	for i := off + 1; i < fr.size; i++ {
		fr.row[i - 1] = fr.row[i];
	}
	fr.size--;
}

func (fr *Frow) insertAt(b int32, off int) {
	// move elems in bulks.
	for i := fr.size; i > off; i-- {
		fr.row[i] = fr.row[i - 1];
	}
	fr.row[off] = byte(b);
	fr.size ++;
}

func (fr *Frow) display(offset int, curse int) {
	cursor.print_hex(int32(offset));
	cursor.puts(" |");

	// hex segment
	for i := 0; i < fr.size; i++ {
		if i == curse {
			cursor.puts(Cs);
		}
		cursor.print_hex_last(int32(fr.row[i]));
		if i == curse {
			cursor.puts(Ce);
		}
		if i == 7 {
		  	cursor.puts(" |");
		}
	}
	cursor.puts(" | ");

	// char segment
	for i := 0; i < fr.size; i++ {
		if i == curse {
			cursor.puts(Cs);
		}
		cursor.print_byte(fr.row[i]);
		if i == curse {
			cursor.puts(Ce);
		}
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
	file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0644);
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
	whiteout();
	it := fl.row;
	if it == &fl.tail {
		return;
	}
	pos := 0;
	for i := 0; i < height; i++ {
		if i == 0 {
			it.display(pos, fl.col);
		} else {
			it.display(pos, -1);
		}
		pos += it.size;
		it = it.next;
		if it == &fl.tail {
			break;
		}
	}
}

// move the cursor up
func (fl *Flist) up() {
	fl.row = fl.row.prev;
	if fl.row == &(fl.head) {
		fl.row = fl.row.next;
	}
}

// move the cursor down
func (fl *Flist) down() {
	fl.row = fl.row.next;
	if fl.row == &(fl.tail) {
		fl.row = fl.row.prev;
	}
}

// move the cursor left
func (fl *Flist) left() {
	if fl.col > 0 {
		fl.col--;
	}
}

// move the cursor right
func (fl *Flist) right() {
	fl.col++;
	if fl.col >= fl.row.size {
		fl.col = fl.row.size - 1;
		if fl.col < 0 {
			fl.col = 0;
		}
	}
}

// remove current byte
func (fl *Flist) remove() {
	fl.row.removeAt(fl.col);

	if fl.row.size == 0 {
		n := fl.row.next;
		fl.row.removeAll();
		fl.row = n;		
	}

	if fl.row.size <= fl.col {
		// fl.row.size != 0
		fl.col = fl.row.size - 1;
	}
}

// insert a byte before cursor
func (fl *Flist) insertBefore(b int32) {
	if fl.row.size == width {
		left := width / 2;
		// too full, split the row first.

		tmp := make([]Frow, 1);
		new := &(tmp[0]);
		new.init();
		new.size = width - left;

		for i := 0; i < new.size; i++ {
			new.row[i] = fl.row.row[i + left];
		}
		fl.row.size = left;

		// insert the row into the back of fl.row.
		nxt := fl.row.next;
		nxt.prev = new;
		fl.row.next = new;
		new.prev = fl.row;
		new.next = nxt;

		if fl.col >= left {
			// make fl.row point to next row.
			fl.col -= left;
			fl.row = new;
		} else {
		}
	}

	fl.row.insertAt(b, fl.col);
}

// insert a byte after cursor
func (fl *Flist) insertAfter(b int32) {
	if fl.row.size == width {
		// split the row
		left := width / 2;
		// too full, split the row first.

		tmp := make([]Frow, 1);
		new := &(tmp[0]);
		new.init();
		new.size = width - left;

		for i := 0; i < new.size; i++ {
			new.row[i] = fl.row.row[i + left];
		}
		fl.row.size = left;

		// insert the row into the back of fl.row.
		nxt := fl.row.next;
		nxt.prev = new;
		fl.row.next = new;
		new.prev = fl.row;
		new.next = nxt;

		if fl.col >= left {
			// make fl.row point to next row.
			fl.col -= left;
			fl.row = new;
		} else {
		}
	}

	if fl.row.size == 0 {
		fl.row.insertAt(b, fl.col);
	} else {
		fl.row.insertAt(b, fl.col + 1);
		fl.col++;
	}
}

// execute commands
func (fl *Flist) exec(args []string, disp bool) {
	if len(args) == 0 {
		// empty commands
		return;
	}

	switch args[0] {
	case "h": fl.left();
	case "j": fl.down();
	case "k": fl.up();
	case "l": fl.right();
	case "w": fl.save();
	case "x": fl.remove();
	case "i": if len(args) > 1 { fl.insertBefore(atoi(args[1])); }
	case "s": if len(args) > 1 { fl.replace(atoi(args[1])); }
	case "a": if len(args) > 1 { fl.insertAfter(atoi(args[1])); }
	default: {  // may want to execute a command multiple times.
		count := atoi(args[0]);
		for i := int32(0); i < count; i++ {
			// execute, do not display.
			fl.exec(args[1:], false);
		}
	}
	}

	// clear and show contents.
	if disp {
		fl.display();
	}
}

// save the file.
func (fl *Flist) save() {
	fl.fobj.Seek(0, 0);
	bwrt := 0;
	for r := fl.head.next; r != &fl.tail; r = r.next {
		if r.size > 0 {
			b, err := fl.fobj.Write(r.row[:r.size]);
			if err != nil {
				break;
			}
			bwrt += b;
		}
	}

	fl.fobj.Truncate(int64(bwrt));
}

// replace current character
func (fl *Flist) replace(b int32) {
	fl.row.row[fl.col] = byte(b);
}

type MyStdin struct{
	fobj *os.File;
}
func (reader MyStdin) Read(p []byte) (n int, err error) {
	return reader.fobj.Read(p);
}
var rd io.Reader = MyStdin{ os.Stdin };
var StdinReader *bufio.Reader = bufio.NewReader(rd);


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
	for {
		msg, err := StdinReader.ReadString('\n');
		if err != nil {
			break;
		}	
		// erase the enter key pressed
		fmt.Printf(Up);

		fl.exec(strings.Fields(msg), true);
	}

	// finally save all changes.
	fl.save();
	fl.fobj.Close();
}
