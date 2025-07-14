package main

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
)

func generateSecretKey(byteLength int) (string, error) {
    key := make([]byte, byteLength)
    _, err := rand.Read(key)
    if err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(key), nil
}

func main() {
    secretKey, err := generateSecretKey(32)
    if err != nil {
        fmt.Println("Error generating secret key:", err)
        return
    }
    fmt.Println("JWT Secret Key:", secretKey)
}
