# Lesson 1
## Install
First get go compiler installed from [here](https://go.dev/dl/), then use this short 
program to test your installation:

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Welcome to the playground!")
	fmt.Println("The time is", time.Now())
	fmt.Printf("%d", 23)
}
```

Execute `go run src.go` and proceed. 

## Configure
You may use `go env` to display all environment variables, and use `go clean -cache` to clean cached files.

# Lesson 2

## Packages
Here's a short example of importing and using packages:

```go
package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("My favorite number is", rand.Intn(10))
}
```

## Functions
Examples of a function perfrom add on **int**:
```go
package main

func add (a int, b int) int {
  return a + b;
}
```

Also, using a tuple-like feature allows you to return multiple values:
```go
func swap(a string, b string) (string, string) {
  return b, a;
}
func main () {
	a, b := swap ("a nice day.", "Today is");
	fmt.Println(a, b);
}
```

## variables
Pay attention to the difference of "=" and ":=" operator.
```go
/* package level var */
var foo, bar bool = true, false;

func main() {
  /* function level var */
	var baz int;
	baz = 42; /* equiv to baz := 42 */
	fmt.Println(baz, foo, bar);
}
```

## basic types
Go's basic types are:
> bool
> string
> int  int8  int16  int32  int64
> uint uint8 uint16 uint32 uint64 uintptr
> byte(uint8)
> float32, float64
> complex64, complex128

Use `T(v)` to convert variable `v` into type `T`. Use `%T` to tell the type of a variable,
E.g.
```go
func main() {
	i := 42;      /* int */
	f := 2 + 2i;  /* complex128 */
	u := uint(i); /* uint */
	fmt.Println(i, f, u);
	fmt.Printf("%T %T %T\n", i, f, u);
}
```

## constants
Declaration of constants using ":=" is not allowed.
```go
func main() {
	const i int = 42;
	const f complex128 = 2 + 2i;
	u := uint(i);
	fmt.Println(i, f, u);
	fmt.Printf("%T %T %T\n", i, f, u);
}
```
Or to create a set of constants:
```go
const (
	Big = 1 << 16
	Small = Big >> 15
)
```