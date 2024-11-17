# Lesson 6
This time we focus on generic type and function.

## Type parameters
Consider:
```go
func Index[T comparable](s []T, x T) int {
    for id, v : range s {
        if v == x { return id; }
    }

    return -1;
}
```
The declaration means that `s` is a slice type of `T` of any type `T` that fulfills 
built-in constraint `comparable`. `comparable` is a useful constraint that makes `==`
and `!=` possible.

## Generic types
Let's use this functionality on our doubly-lined list:
```go
type list_elem[T any] struct {
    prev *list_elem[T],
    next *list_elem[T],
    val T,
}
```