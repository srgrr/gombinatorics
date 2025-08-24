// Gombinator is a functional-style library designed to play nicely with goroutines and channels
// It allows the user to use functions like map, filter, partition, etc
// Most functions require the user to provide a context for early cancellation
// This library is NOT meant as a replacement for traditional loops: channels add overhead
// Consider using gombinator only when you were intending to distribute load between
// different goroutines in the first place
package gombinator

// Buffer size for channels
var GBufferSize = 64
