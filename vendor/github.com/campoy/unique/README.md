[![GoDoc](https://godoc.org/github.com/campoy/unique?status.svg)](https://godoc.org/github.com/campoy/unique) [![Build Status](https://travis-ci.org/campoy/unique.svg)](https://travis-ci.org/campoy/unique)


# unique

Package unique provides primitives for sorting slices removing
repeated elements.

## a quick example

Executing this code:

[embedmd]:# (example_test.go /.*s :=/ /Println.*/)
```go
	s := []int{3, 5, 1, 7, 2, 3, 7, 5, 2}
	less := func(i, j int) bool { return s[i] < s[j] }
	unique.Slice(&s, less)
	fmt.Println(s)
```

Will output the following:

```
[1 2 3 5 7]
```
