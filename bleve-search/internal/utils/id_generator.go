package utils

import (
	"github.com/google/uuid"
)

func GenerateUniqueID() string {
	// Implement your unique ID generation logic here (e.g., UUIDs, timestamp-based, etc.)
	// Example: Using UUID
	id := uuid.New().String()
	return id

}
