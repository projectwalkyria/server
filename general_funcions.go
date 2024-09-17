package main

import (
    "crypto/sha256"
    "encoding/hex"
)

func tokenToSha256(token string) (string) {
    hash := sha256.Sum256([]byte(token))
    hashStr := hex.EncodeToString(hash[:])
    return hashStr
}
