package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

// Header represents the header part of a JWT.
type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// Payload represents the payload part of a JWT.
type Payload struct {
	Id  int   `json:"id"`
	Iat int64 `json:"iat"`
	Exp int64 `json:"exp"`
}

// EncodeSegment encodes the header or payload into a URL-safe base64-encoded string.
func EncodeSegment(segment interface{}) (string, error) {
	bytes, err := json.Marshal(segment)
	if err != nil {
		return "", err
	}

	encoded := base64.URLEncoding.EncodeToString(bytes)
	encoded = strings.TrimRight(encoded, "=") // Remove any trailing '='
	return encoded, nil
}

// CreateToken creates a JWT token with a given secret.
func CreateToken(id int) (string, error) {
	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	payload := Payload{
		Id:  id,
		Iat: time.Now().UTC().Unix(),
		Exp: time.Now().UTC().Add(time.Hour * 24).Unix(),
	}

	encodedHeader, err := EncodeSegment(header)
	if err != nil {
		return "", err
	}

	encodedPayload, err := EncodeSegment(payload)
	if err != nil {
		return "", err
	}

	signingInput := encodedHeader + "." + encodedPayload
	signature := strings.TrimRight(computeHMAC256(signingInput, "your-256-bit-secret"), "=")

	return signingInput + "." + signature, nil
}

// computeHMAC256 computes HMAC SHA-256 signature.
func computeHMAC256(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
