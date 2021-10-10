package service

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

/* The PBKDF2_* constants may be changed without breaking existing stored hashes. */
const (
	// pbkdf2Iterations sets the amount of iterations used by the PBKDF2 hashing algorithm
	pbkdf2Iterations int = 4096

	// slatBytes sets the amount of bytes for the salt used in the PBKDF2 / scrypt hashing algorithm
	slatBytes int = 32
	// hashBytes sets the amount of bytes for the hash output from the PBKDF2 / scrypt hashing algorithm
	hashBytes int = 64
)

/* altering the HASH_* constants breaks existing stored hashes */
const (
	// hashSections identifies the expected amount of parameters encoded in a hash generated and/or tested in this package
	hashSections int = 4
	// hashIterationIndex identifies the position of the iteration count used by PBKDF2 in a hash generated and/or tested in this package
	hashIterationIndex int = 1
	// hashSaltIndex identifies the position of the used salt in a hash generated and/or tested in this package
	hashSaltIndex int = 2
	// hashPBKDF2Index identifies the position of the actual password hash in a hash generated and/or tested in this package
	hashPBKDF2Index int = 3
)

// CreateHash creates a salted cryptographic hash with key stretching (PBKDF2),
// suitable for storage and usage in password authentication mechanisms.
func createHash(password string) (string, error) {
	salt := make([]byte, slatBytes)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := pbkdf2.Key([]byte(password), salt, pbkdf2Iterations, hashBytes, sha512.New)

	/* format: algorithm:iterations:salt:hash */
	return fmt.Sprintf(
		"%s:%d:%s:%s", "sha512", pbkdf2Iterations,
		base64.StdEncoding.EncodeToString(salt), base64.StdEncoding.EncodeToString(hash),
	), err
}

// ValidatePassword hashes a password according to the setup found in the correct hash string
// and does a constant time compare on the correct hash and calculated hash.
func validatePassword(password string, correctHash string) bool {
	params := strings.Split(correctHash, ":")
	if len(params) < hashSections {
		return false
	}
	it, err := strconv.Atoi(params[hashIterationIndex])
	if err != nil {
		return false
	}
	salt, err := base64.StdEncoding.DecodeString(params[hashSaltIndex])
	if err != nil {
		return false
	}
	hash, err := base64.StdEncoding.DecodeString(params[hashPBKDF2Index])
	if err != nil {
		return false
	}

	testHash := pbkdf2.Key([]byte(password), salt, it, len(hash), sha512.New)

	return subtle.ConstantTimeCompare(hash, testHash) == 1
}
