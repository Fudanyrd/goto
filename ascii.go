package main

/** ascii: a funny game to play with ascii characters!
 * Usage: h = move left, j = move down, k = move up, l = move right;
 * s <new char> to replace current character(colored in red) with new char.
 * <ctrl-z> to exit the game.
 * THIS IS NOT A PRACTICAL TOOL, IT IS BUILD PRIMARILY FOR FUN.
 */

import (
	"strings"
	"fmt"
	"bufio"
	"io"
	"os"
)

const (
	// height of scree.
	height int32 = 15
	width int32 = 16
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

func print_byte(b byte) {
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

type Line [width]byte;
// display a line.
func (l *Line) render(start int32, curse int32) {
	print_hex(start);

	// print hex data
	fmt.Printf(" |");
	for i := int32(0); i < width; i++ {
		b := [width]byte(*l)[i];

		if i == curse {
			fmt.Printf("\033[01;31m");
			print_hex_last(int32(b));
			fmt.Printf("\033[0m");
		} else {
			print_hex_last(int32(b));
		}
		if i == 7 {
		  	fmt.Printf(" |");
		}
	}
	fmt.Printf(" | ");

	// char repr
	fmt.Printf(" | ");
	for i := int32(0); i < width; i++ {
		b := [width]byte(*l)[i];
		if i == curse {
			fmt.Printf("\033[01;31m");
			print_byte(b);
			fmt.Printf("\033[0m");
		} else {
			print_byte(b);
		}
	}

	fmt.Printf(" | ");
}

type Screen struct {
	buf [height]Line;
	h int32;  // vertical position
	w int32;  // horizontal position
	pos int32;  // position in file
}

func (s *Screen) clear() {
	for i := int32(0); i < height; i++ {
		fmt.Printf(Up);
	}

	return;
}
func (s *Screen) render() {
	for i := int32(0); i < height; i++ {
		if i == s.h {
			s.buf[i].render(s.pos + i * width, s.w);
		} else {
			s.buf[i].render(s.pos + i * width, int32(-1));
		}
		fmt.Printf(Down);
	}
	return;
}

func (s *Screen) down() {
	if s.h + 1 < height {
		s.h += 1;
	}
}
func (s *Screen) up() {
	if s.h > 0 {
		s.h -= 1;
	}
}
func (s *Screen) left() {
	if s.w > 0 {
		s.w -= 1;
	}
}

func (s *Screen) right() {
	if s.w + 1 < width {
		s.w += 1;
	}
}

func atoi(word string) int32 {
	var ret int32 = int32(0);
	buf := []rune(word);

	for i := 0; i < len(buf); i++ {
		ret = ret * int32(10);
		ret += int32(buf[i] - rune('0'));
	}
	return ret;
}

func (s *Screen) replace(after string) {
	s.buf[s.h][s.w] = byte(atoi(after));
}

func (s *Screen) cmd(line string) {
	words := strings.Fields(line);
	for i := 0; i < len(line); i++ {
		fmt.Printf("\b");
	}
	fmt.Printf(Up);

	var valid bool = true;
	if len(words) == 0 {
		return;
	}
	switch words[0] {
	case "j": s.down(); valid = true;
	case "k": s.up();valid = true;
	case "h": s.left();valid = true;
	case "l": s.right();valid = true;
	case "s": s.replace(words[1]);valid = true;
	default: valid = false;
	}
	if valid {
		s.clear(); 
		s.render();
	}
}

func (s *Screen) run(b rune) {

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
	s := Screen{};
	s.render();
	s.clear();
	s.render();

	for {
		line, err := StdinReader.ReadString('\n');
		if len(line) > 0 {
			s.cmd(line);
		}
		if err != nil {
			break;
		}
	}
}
