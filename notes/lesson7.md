# Lesson 7
This time we study golang's concurrency support.

## Threads
In golang, goroutine(i.e. threads in go) can be created via **go** keyword:
```go
var sum int = 0;
// a worker thread
func Twork() {
	for i := 0; i < 10000; i++ {
		sum += 1;
	}
}
func main() {
	// start two worker thread.
	go Twork();
	go Twork();
	// print 0.
	fmt.Println(sum);
}
```

## Channels
Channels provide a way to send and receive values. 
```go
func main() {
    // make a channel, type int
    var ch chan int = make(chan int);
    // send 0x123 to the channel
    ch <- 0x123;
    // receive the value
    var val int = <- ch;
    fmt.Printf("%x\n", val);
}
```

Channels can also be buffered, sends to be buffered chan will be blocked when buffer is full:
```go
func main() {
    var ch chan int = make(chan int, 2);
    ch <- 0x123;
    ch <- 0x456;
    x := <- ch;
    y := <- ch;
    fmt.Printf("%x %x\n", x, y);
}
```

## Range and close
`v, ok := <- ch` ok is false if no more value can be received; the loop `for i:= range c` receives value
until `close` is called on c.

Example:
```go
func count(n int, ch chan int) {
	for i := 0; i < n; i++ {
		ch <- i;
	}
	ch <- cap(ch); // cap works on ch!
	close(ch);
}
func main() {
	var n int = 10;
	ch := make(chan int, n);
	go count(n, ch);

	for i := range ch {
		fmt.Println(i);
	}
	/* or equivalently: 
	for i,ok := <-ch; ok; i,ok := <-ch {}
	*/
}
```

## select
`select` statement lets a goroutine wait on multiple communication operations.
A select blocks if no case is ready; if multiple are ready, randomly select one:
```go
func count(c, quit chan int) {
	for i := 0; i < 8; i++ {
		select {
		case c <- i: 
		case <- quit:
			fmt.Println("done.");
			return;
		}
	}
}
func main() {
	ch := make(chan int);
	quit := make(chan int);
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<- ch);
		}
		quit <- 1;
	}()
    // 0, 1, ..., 7
	count(ch, quit);
}
```