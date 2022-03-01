package app

import (
	"context"
	"crypto/subtle"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"github.com/theplant/luhn"
	"path"
	"runtime"
	"strconv"
)

// GetBasePath prepares base path of the project
func GetBasePath() string {
	_, b, _, _ := runtime.Caller(0)
	return path.Dir(path.Dir(path.Dir(b)))
}

// ComparePass compares hashed passwords
func ComparePass(ctx context.Context, expected string, actual string) bool {
	return (subtle.ConstantTimeCompare([]byte(expected), []byte(actual)) == 1)
}

// PrepOrderNumber extracts integer from string and checks if this is valid Luhn number
func PrepOrderNumber(ctx context.Context, n []byte) (bool, int64, error) {
	id, err := strconv.Atoi(string(n))
	if err != nil {
		return false, 0, fmt.Errorf("convert to int fialed: %v", err)
	}

	return luhn.Valid(id), int64(id), nil
}

// UserIDFromContext gets user if from context (it gets there from jwt token parsing)
func UserIDFromContext(ctx context.Context) (int64, error) {
	_, uID, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0, err
	}
	return int64(uID["user_id"].(float64)), nil
}