package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type User struct {
	Username string
	Email    string
	Password string
}

func (u *User) SetPassword(password string) {
	hash := sha256.Sum256([]byte(password))
	u.Password = hex.EncodeToString(hash[:])
}

func (u *User) VerifyPassword(password string) bool {
	hash := sha256.Sum256([]byte(password))
	return u.Password == hex.EncodeToString(hash[:])
}

func main() {
	user := User{
		Username: "Ilya",
		Email:    "Murometc@mail.ru",
	}

	user.SetPassword("123456789")
	fmt.Printf("Пользователь: %s, Email: %s\n", user.Username, user.Email)
	fmt.Printf("Хэш пароля: %s\n", user.Password)

	if user.VerifyPassword("123456789") {
		fmt.Println("Пароль верный!")
	} else {
		fmt.Println("Пароль неверный!")
	}

	if user.VerifyPassword("12345678") {
		fmt.Println("Пароль верный!")
	} else {
		fmt.Println("Пароль неверный!")
	}
}
