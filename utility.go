package main

import (
    "regexp"
    "unicode"
	// "errors"
)

// ? Password Rules
// Minimum length : 8 Characters
// Must include at least one `UPPERCASE` letter
// Must include at least one `lowercase` letter
// Must include at least one digit
// Must include at least one special character (e.g., !@#$%^&*...)
func IsPasswordValid(password string) bool {
    var (
        hasMinLen  = false
        hasUpper   = false
        hasLower   = false
        hasNumber  = false
        hasSpecial = false
    )
    if len(password) >= 8 {
        hasMinLen = true
    }
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsDigit(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }
    return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

// Must be between 3 and 20 characters long
// No spaces or special characters
func IsUsernameValid(username string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
    return re.MatchString(username)
}

// TODO send email and get user to verify
func IsEmailValid(email string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(email)
}

// Only alphabetic characters
// No spaces
// Between 2 and 50 Characters
func IsNameValid(name string) bool {
    re := regexp.MustCompile(`^[a-zA-Z]{2,50}$`)
    return re.MatchString(name)
}