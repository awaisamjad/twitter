package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"

	"github.com/gorilla/sessions"
)


func initSessionStore() {
    sessionKey := os.Getenv("SESSION_KEY")
    store = sessions.NewCookieStore([]byte(sessionKey))
}

func rotateSessionKey() {
    newSessionKey := generateRandomKey(32) // Generate a new 32-byte key
    os.Setenv("SESSION_KEY", newSessionKey)
    store = sessions.NewCookieStore([]byte(newSessionKey))
    log.Println("Session key rotated")
}

func generateRandomKey(length int) string {
    key := make([]byte, length)
    _, err := rand.Read(key)
    if err != nil {
        log.Fatal("Failed to generate random key:", err)
    }
    return base64.StdEncoding.EncodeToString(key)
}