package password

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// VerifyPassword compares password and the hashed password
func VerifyPassword(passwordHash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}

// HashPassword creates a bcrypt password hash
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 11)
}

// VerifyPasswordSha1 compares the password and hashed password using SHA-1, and generates a new bcrypt hash replacement
func VerifyPasswordSha1(passwordHash, passwordSalt, password string) error {
	if HashPasswordSha1(password, passwordSalt) != passwordHash {
		return errors.New("hashedPassword is not the hash of the given password")
	}
	return nil
}

// UCS2Encode Encodes string to UCS2 encoding.
func UCS2Encode(s []byte) []byte {
	e := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	es, _, err := transform.Bytes(e.NewEncoder(), s)
	if err != nil {
		return s
	}
	return es
}

// HashPasswordSha1 creates a SHA-1 password hash
func HashPasswordSha1(password, salt string) string {
	saltBuffer, _ := base64.StdEncoding.DecodeString(salt)
	passwordBuffer := UCS2Encode([]byte(password))
	passwordWithSalt := append(saltBuffer[:], passwordBuffer[:]...)
	hash := sha1.Sum(passwordWithSalt)
	return base64.StdEncoding.EncodeToString(hash[:])
}
