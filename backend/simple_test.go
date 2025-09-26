package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Testar hash do admin123
	password := "admin123"
	hash := "$2a$10$N9qo8uLOickgx2ZMRZoMye.IjdQvOQ5eqGStBUKx6XgKnrQvp.Fl6"
	
	fmt.Printf("Testando senha: %s\n", password)
	fmt.Printf("Hash do banco: %s\n", hash)
	
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("❌ ERRO: %v\n", err)
	} else {
		fmt.Printf("✅ SUCESSO: Senha confere!\n")
	}
	
	// Gerar novo hash para comparação
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Erro ao gerar novo hash: %v\n", err)
	} else {
		fmt.Printf("Novo hash gerado: %s\n", string(newHash))
	}
}
