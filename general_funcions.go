package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"
)

func tokenToSha256(token string) string {
	hash := sha256.Sum256([]byte(token))
	hashStr := hex.EncodeToString(hash[:])
	return hashStr
}

func logStuff(content string) {
	// Get the current date and time
	currentTime := time.Now()

	// Log the custom formatted date and content
	log.Println(currentTime.Format("") + " " + content)
}
