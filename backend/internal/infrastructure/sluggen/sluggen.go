package sluggen

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/sqids/sqids-go"
)

// Generator generates human-readable slugs using Sqids
type Generator struct {
	sqids *sqids.Sqids
}

// NewGenerator creates a new slug generator
func NewGenerator() (*Generator, error) {
	// Initialize Sqids with a custom alphabet (URL-safe characters)
	s, err := sqids.New(sqids.Options{
		MinLength: 8, // Minimum slug length
	})
	if err != nil {
		return nil, err
	}

	return &Generator{
		sqids: s,
	}, nil
}

// Generate creates a unique, human-readable slug
// It uses the current timestamp and a random number to ensure uniqueness
func (g *Generator) Generate() (string, error) {
	// Use current timestamp in milliseconds
	now := time.Now().UnixMilli()

	// Add a random number for additional uniqueness
	randomNum, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}

	// Encode both numbers into a slug
	slug, err := g.sqids.Encode([]uint64{uint64(now), randomNum.Uint64()})
	if err != nil {
		return "", err
	}

	return slug, nil
}
