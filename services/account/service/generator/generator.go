package generator

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func GenerateAccountName(ID uuid.UUID) string {
	rand.Seed(time.Now().UnixNano())

	animal := animals[rand.Intn(len(animals))]

	return "Anonymous " + animal + "-" + fmt.Sprintf("%x", md5.Sum([]byte(animal+ID.String())))[:6]
}
