package utils

import (
    "crypto/rand"
    "math/big"
)

const (
    charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    codeLength = 6
)

func GenerateShortCode() (string, error) {
    code := make([]byte, codeLength)
    for i := range code {
        num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
        if err != nil {
            return "", err
        }
        code[i] = charset[num.Int64()]
    }
    return string(code), nil
}