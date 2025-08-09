package main

import (
	"fmt"
	"sync"
	sets "github.com/srgrr/gombinatorics/sets"
)

type User struct {
	email string
}

type Spam struct {
	message string
}

func sendSpam(user User, spam Spam, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Sent spam\t\"%s\"\t to %s,\thope they will buy my stuff!\n", spam.message, user.email)
}

func main() {
	emails := []User{{email: "sergio@raccoon.me"}, {email: "raquel@cat.me"}}
	spam := []Spam{{message: "Raccoon plushies now 10% discount"}, {message: "Brown hair dye now 5% discount"}}
	var wg sync.WaitGroup
	for emailAndSpam := range sets.CartesianProduct(emails, spam) {
		wg.Add(1)
		go sendSpam(emailAndSpam.First, emailAndSpam.Second, &wg)
	}
	wg.Wait()
}
