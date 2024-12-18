package main

import (
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

func TestPort(t *testing.T) {
    port := os.Getenv("PORT")
    if port == "" {
        t.Error("PORT environment variable is not set")
    } else {
        t.Log("PORT:", port)
    }
}