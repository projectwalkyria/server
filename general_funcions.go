package main

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

func tokenToSha256(token string) (string) {
    hash := sha256.Sum256([]byte(data))
    hashStr := hex.EncodeToString(hash[:])
    fmt.Println(hashStr)
}
