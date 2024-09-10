package helper

import (
	"crypto/md5"
	"encoding/base64"
	"regexp"
)

func GenerateSecurityToken(input string) string {
	// Convert input to bytes
	temp := []byte(input)

	// Generate MD5 hash
	hash := md5.New()
	hash.Write(temp)
	tokenBytes := hash.Sum(nil)

	// Encode hash using Base64
	tokenBase64 := base64.StdEncoding.EncodeToString(tokenBytes)

	// Remove the last two characters
	if len(tokenBase64) > 2 {
		tokenBase64 = tokenBase64[:len(tokenBase64)-2]
	}

	// Remove non-alphanumeric characters
	re := regexp.MustCompile("[^A-Za-z0-9]+")
	token := re.ReplaceAllString(tokenBase64, "")

	return token
}
