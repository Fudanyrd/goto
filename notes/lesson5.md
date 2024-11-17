# Lesson 5
Most examples in this class will be long. So **be patient**.

## <a href="https://go.dev/tour/methods/2">method</a>
Consider the following doubly-linked list declaration:
```go
type list_elem struct {
	prev *list_elem;
	next *list_elem;
};
type list struct {
	head list_elem;
	tail list_elem;
};
```

Suppose we want to do:
```go
func main() {
    var lst list = list{};
    lst.init(); // initialize the list
}
```

Then `init` should be a method of type `list`, i.e. `init` should have a **receiver argument**:
```go
func (lst list*) init() {
	lst.head.next = &lst.tail;
	lst.tail.prev = &lst.head;
}
```

> This receiver argument is a pointer. Why?

Else, the list object will be **copied** and destroyed after execution of `init`!

## <a href="https://go.dev/tour/methods/9">interface</a>
An **interface type** is defined as a set of signatures. Use animal as an example:
```go
type Animal interface {
	MakeNoise()
};
```
And we have two kinds of "animals":
```go
type Number int;
/* Number implements make noise */
func (num Number) MakeNoise() {
	fmt.Println("HAHAHAHA");
}

type Dog struct {
	age int;
};
/* *Dog implement make noise */
func (dog *Dog) MakeNoise() {
	if dog.age <= 4 {
		fmt.Println("Waowaoaoaoaoao");
	} else {
		fmt.Println("Bark!");
	}
}
```

Interestingly, we can use an object of `Animal` to store a receiver of the signature:
```go
func main() {
	var obj Animal = nil;
	d := Dog{6};
	obj = &d; // why pointer? for the receiver is a pointer!
	obj.MakeNoise();

	var n Number = 6; // can be cast to 'int'.
	obj = n;  // why not a pointer? same reason.
	obj.MakeNoise();
}
```

> Also, you can use fmt.Printf("%v %T", obj, obj) to print the value and type of an interface.

You may use empty interface to hold any kinds of value:
```go
func main() {
	var i interface{} = 3;
	fmt.Printf("%v %T\n", i, i);  // 3 int
	i = Dog{23};
	fmt.Printf("%v %T\n", i, i);  // 23 main.dog
}
```

## <a href="https://go.dev/tour/methods/15">Type Assertion</a>
Grammar: `t, ok := i.(T)`, If interface `i` has type `T`, then `t` stores value of `i`, 
and `ok` is true, otherwise ok is false:
```go
func addFn(l, r int) int { return l + r; }
func main() {
	var i interface{} = addFn;
	t, ok := i.(func (l int, r int) int);	
	fmt.Println(t, ok);

	s, ok := i.(string);
	fmt.Println(s, ok);
}
```

It is also possible to do type in switches:
```go
switch v := i.(type) {
case T1: // here v has type T1;
case T2: // here v has type T2;
default: // v has same type as I.
}
```

## <a href="https://go.dev/tour/methods/17">Stringers</a>
This is defined in **fmt** package:
```go
type Stringer interface {
    String()
}
```

Example:
```go
func (d Dog) String() string {
    return fmt.Sprintf("%v years", d.age);
}
func main() {
	fmt.Println(Dog{23}); // 23 years
}
```

## <a href="https://go.dev/tour/methods/19">errors</a>
The error interface:
```go
type error interface {
    // built-in
    Error() string
};
```

Function often returns an error value:
```go
type MyError struct {
	what string;
};
var DivByZero MyError = MyError{"division by zero"};
func (err *MyError) Error() string {
	return err.what;
}

func divFn(l, r int) (int, error) {
	if r == 0 {
		return -1, &DivByZero;
	}
	return l / r, nil;
}
func main() {
	l, r := 2, 0;
	fmt.Println(divFn(l, r)); // -1 division by zero
	fmt.Println(divFn(637, 231)); // 2 nil
}
```

## <a href="https://go.dev/tour/methods/21">Readers</a>
The `io` package specifies the `io.Reader` interface:
```go
func (rd T) Read(b []byte) (int, error) {...}
```
It returns an `io.EOF` error when the stream ends.

## Resources
Learn more about go packages, go <a href="https://pkg.go.dev/std">here</a>.