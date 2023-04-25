package utils

import (
	"math/rand"
	"strings"
	"time"
)

const (
	ALPHABET = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"
)

// Encode63 Функция для генерации случайного набора символов из хранящихся в переменной ALPHABET
func Encode63() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune(ALPHABET)
	length := 10
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()
	return str
}
