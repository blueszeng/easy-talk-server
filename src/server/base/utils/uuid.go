package utils

import (
	"github.com/satori/go.uuid"
)

func GetUUID() string {
	return uuid.NewV1().String()
}
