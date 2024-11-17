# Lesson 3

## for loop
Example:
```go
/* computes sum of 1,2,...,n */
func sumtop(n int) int {
	sum := 0;
	for i := 1; i <= n; i++ {
		sum += i;
	}
	return sum;
}
```

Similar to C, the following is an infinite loop:
```go
func hang() {
  for {}
}
```

## if-else statement
Example:
```go
func tellday (day int) {
	if (day % 7 + 1) <= 5 {
		fmt.Println("Weekdays are bad. Emm...");
	} else {
		fmt.Println("Weekends are good! Yay!");
	}
	return;
}
```

## switch statement
Example:
```go
	switch day := (date % 7 + 1); day {
	case 1:
		fmt.Println("Monday is worst. Hmm...");
	case 6: 
		fmt.Println("Weekends are best");
	case 7:
		fmt.Println("Weekends are best");
	default:
		fmt.Println("Weekdays are so-so.");
	}
```

## defer statement
Deferred statement will be run at function return.
```go
/* 9, 8, ..., 0, Wow! */
func main() {
	defer fmt.Println("Wow!");
	for i := 0; i < 10; i+=1 {
		defer fmt.Printf("%d, ", i);
	}
}
```

# Lesson 4

## Structs
Suppose you're coding a doubly-linked list:
```go
type list_elem struct {
	prev *list_elem;
	next *list_elem;
}
type list struct {
	head list_elem;
	tail list_elem;
}; /* semicolon is optional */
```

## Pointers
Keypoints: 
<ul>
 <li>`*T` is pointer type; `&` to get address; `*` is deference operator. </li>
 <li> There's NO pointer arithmetic; </li>
 <li> There's no pointer-integer conversion; </li>
 <li> There's no different type pointer conversion; </li> 
 <li> null pointer in golang is "nil"(Think nil as a keyword) </li>
</ul>

```go
func list_end (lptr *list) *list_elem {
	/* explicit dereference */
	return &((*lptr).tail);
}
func list_start (lptr *list) *list_elem {
	/* implicit dereference */
	return lptr.head.next;
}
```

## fixed-size array
`[n]T` declares an array of length n(const), type T. Example:
```go
const len int = 4 * 2;
func main() {
	var arr [len] int = [len] int {1, 2, 3, 4};
	// ...
}
```

## Array slice
Slice are like reference to array members.

```go
func printSlice(s []int) {
	fmt.Printf("len = %d, cap = %d\n", len(s), cap(s));
}
func main() {
	var arr [8] int = [8] int {1, 2, 3, 4};
	var slice = arr[:2];
	printSlice(slice); /* len = 2, cap = 8 */
}
```

Default value to slice is **nil**:
```go
func main() {
	var slice []int = nil;
	printSlice(slice); /* len = 0, cap = 0 */
}
```

## Dynamically-sized array
Use the **built-in** make and append function:
```go
func main() {
	var slice []int = make([]int, 5);
	slice = append(slice, 3);
	fmt.Println(slice); /* [0 0 0 0 0 3] */
	printSlice(slice);  /* len = 6, cap = 6 */
}
```

## range
The range form of **for** loop iterates over a slice. The first is the **index**, 
the second is **a copy of that slice element**.
```go
var primes []int = []int {2, 3, 5, 7, 11};
func main() {
	for idx, prime := range primes {
		fmt.Printf("%d prime is %d.\n", idx, prime);
	}
}
```

## map
`map[K]V` declares a map of key type `K` and value type `V`. Map literal:
```go
var mi map[int]int  = map[int]int{
	1: 2,
	2: 3,
	3: 4,
};
```

Can also use **make** to create a map:
```go
func main () {
	mi = make(map[int]int);
	mi[2] = -3;
	fmt.Println(mi); /* map[2:-3] */
}
```

To delete a key(say, 2), do:
```go
delete(mi, 2);
```

To test the existence of a key(say, 2), do:
```go
val, ok = mi[2];
```

## function(revisited)
Functions are like other objects: they can be assigned!
```go
var addFn func (l int, r int) int = func (l int, r int) int {
	return l + r;
}
var multFn func (l int, r int) int = func (l int, r int) int {
	return l * r;
}

func operator(operation string) func (l int, r int) int {
	switch operation {
	case "+": return addFn;
	case "*": return multFn;
	default: return func (l int, r int) int { return 0; }
	}
}
```
