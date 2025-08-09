# Gombinatorics 🎲

## Basic Types
### Pair
```go
type Pair[P any, Q any] struct {
	First  P
	Second Q
}
```

## Cartesian Product
```go
func CartesianProduct[P any, Q any](A []P, B []Q) <-chan Pair[P, Q]
```

This function *does* compute the cartesian product of two slices of any given types, but it **does not** compute it at once.
Values are computed and served on demand.

For instance, this code
```go
	for pair := range CartesianProduct([]string{"🦝", "🐀"}, []int{0, 1}) {
		fmt.Printf("%v", pair)
	}
```

Prints the following output

```
{🦝 0}
{🦝 1}
{🐀 0}
{🐀 1}
```