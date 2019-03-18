package main

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func getValidUsername(u string) (uuid.UUID, error) {
	// We're dealing with a DB that might be case sensitive. Don't parse
	// uppercase as if it is valid
	if strings.ToLower(u) != u {
		return uuid.UUID{}, errors.New("need lowercase username")
	}
	uname, err := uuid.Parse(u)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uname, nil
}

func validKey(k string) bool {
	kn := sanitizeString(k)
	if utf8.RuneCountInString(k) == 40 && utf8.RuneCountInString(kn) == 40 {
		// Correct length and all chars valid
		return true
	}
	return false
}

func validSubdomain(s string) bool {
	// validate sane domain name part, like "com" or "a-b-c.com"
	if regexp.MustCompile(`^([a-z0-9]([a-z0-9-]*[a-z0-9])?)([.][a-z0-9]([a-z0-9-]*[a-z0-9])?)*$`).
		MatchString(s) {
		return true
	}
	// The following is not needed, as the regex above matches as well
	// validate uuid, like "d25989a6-c59a-4670-b294-a8cb0c5ad8d2"
	// _, err := uuid.Parse(s)
	// if err == nil {
	// 	return true
	// }
	return false
}

func validTXT(s string) bool {
	sn := sanitizeString(s)
	if utf8.RuneCountInString(s) == 43 && utf8.RuneCountInString(sn) == 43 {
		// 43 chars is the current LE auth key size, but not limited / defined by ACME
		return true
	}
	return false
}

func correctPassword(pw string, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)); err == nil {
		return true
	}
	return false
}
