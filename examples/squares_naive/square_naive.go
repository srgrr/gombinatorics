package main

func main() {
	// 1. Create a slice of numbers from 1 to 10
	numbers := make([]int, 10)
	for i := 1; i <= 10; i++ {
		numbers[i-1] = i
	}
	// 2. Filter even numbers
	evenNumbers := make([]int, 0)
	for _, n := range numbers {
		if n%2 == 0 {
			evenNumbers = append(evenNumbers, n)
		}
	}
	// 3. Map the even numbers to their squares
	squaredEvenNumbers := make([]int, len(evenNumbers))
	for i, n := range evenNumbers {
		squaredEvenNumbers[i] = n * n
	}
	for _, n := range squaredEvenNumbers {
		println(n)
	}
}
