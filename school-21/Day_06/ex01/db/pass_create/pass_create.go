package main

// Этот файл поможет хешировать пароль с использованием библиотеки golang.org/x/crypto/bcrypt.
// Он генерирует безопасный хеш для заданного пароля, который можно использовать для безопасного
// хранения паролей пользователей.

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Определите пароль для хеширования
	password := "pass123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hash))
}
