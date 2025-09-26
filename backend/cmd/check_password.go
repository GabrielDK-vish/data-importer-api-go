package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := flag.String("password", "", "senha em texto para validar")
	hash := flag.String("hash", "", "hash bcrypt para comparar")
	flag.Parse()

	if *password == "" || *hash == "" {
		fmt.Println("uso: go run ./cmd/check_password.go -password <senha> -hash <hash_bcrypt>")
		os.Exit(2)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*hash), []byte(*password)); err != nil {
		fmt.Println("MATCH=FALSE")
		os.Exit(1)
	}
	fmt.Println("MATCH=TRUE")
}



