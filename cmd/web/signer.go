package main

import (
	"fmt"
	"strings"
	"time"

	goalone "github.com/bwmarrin/go-alone"
)

// my Key : Should not store in code but I am storing as it is just side Project
const secret = "sainipankaj"

var secretKey []byte

// NewURLSigner creates a new signer
func NewURLSigner() {
	secretKey = []byte(secret)
}

// GenerateTokenFromString generates a signed token
func GenerateTokenFromString(data string) string {
	var urlToSign string

	s := goalone.New(secretKey, goalone.Timestamp)
	if strings.Contains(data, "?") {
		urlToSign = fmt.Sprintf("%s&hash=", data)
	} else {
		urlToSign = fmt.Sprintf("%s?hash=", data)
	}

	tokenBytes := s.Sign([]byte(urlToSign))
	token := string(tokenBytes)

	return token
}

// VerifyToken verifies a signed token
func VerifyToken(token string) bool {
	s := goalone.New(secretKey, goalone.Timestamp)
	_, err := s.Unsign([]byte(token))

	if err != nil {
		// signature is not valid. Token was tampered with, forged, or maybe it's
		// not even a token at all! Either way, it's not safe to use it.
		fmt.Println("This is not from our database, Something is wrong")
		return false
	}
	// valid hash
	return true

}

// Expired checks to see if a token has expired
func Expired(token string, minutesUntilExpire int) bool {
	s := goalone.New(secretKey, goalone.Timestamp)
	ts := s.Parse([]byte(token))

	// time.Duration(seconds)*time.Second
	return time.Since(ts.Timestamp) > time.Duration(minutesUntilExpire)*time.Minute
}
