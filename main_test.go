package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        panic("Error loading .env file")
    }

    // Run the tests
    code := m.Run()

    // Exit with the test code
    os.Exit(code)
}

func Test(t *testing.T) {
    key := make([]byte, 33)
    _, err := rand.Read(key)
    if err != nil {
        log.Fatal("Failed to generate random key:", err)
    }
    fmt.Println(base64.StdEncoding.EncodeToString(key))
}