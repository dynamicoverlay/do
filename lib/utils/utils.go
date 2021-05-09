package utils

import "math/rand"

//StrValue converts a string pointer to a string
func StrValue(ptr *string) string {
	if ptr == nil {
		return ""
	}

	return *ptr
}

//IntValue converts a int pointer to an int
func IntValue(ptr *int) int {
	if ptr == nil {
		return -1
	}

	return *ptr
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

//GenerateRandomString generates an alphernumic random string that is the length of the argument provided
func GenerateRandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
