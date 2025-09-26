package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := flag.String("password", "", "senha em texto para gerar hash bcrypt")
	cost := flag.Int("cost", bcrypt.DefaultCost, "custo do bcrypt (padrão 10)")
	flag.Parse()

	if *password == "" {
		fmt.Println("uso: go run ./cmd/generate_hash.go -password <senha> [-cost 10]")
		os.Exit(2)
	}

	h, err := bcrypt.GenerateFromPassword([]byte(*password), *cost)
	if err != nil {
		fmt.Println("erro:", err)
		os.Exit(1)
	}

	fmt.Println(string(h))
}



