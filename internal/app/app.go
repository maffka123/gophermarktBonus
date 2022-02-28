package app

import (
	"context"
	"crypto/subtle"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"github.com/theplant/luhn"
	"strconv"
)

func ComparePass(ctx context.Context, expected string, actual string) bool {
	return (subtle.ConstantTimeCompare([]byte(expected), []byte(actual)) == 1)

}

func PrepOrderNumber(ctx context.Context, n []byte) (bool, int64, error) {
	id, err := strconv.Atoi(string(n))
	if err != nil {
		return false, 0, fmt.Errorf("convert to int fialed: %v", err)
	}

	return luhn.Valid(id), int64(id), nil
}

func UserIDFromContext(ctx context.Context) (int64, error) {
	_, uID, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0, err
	}
	return int64(uID["user_id"].(float64)), nil
}
