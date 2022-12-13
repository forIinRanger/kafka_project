package adapters

import (
	"context"
	"errors"
	"strings"
)

type ValidatorImpl struct{}

func NewValidator() ValidatorImpl {
	return ValidatorImpl{}
}
func (v ValidatorImpl) Validate(ctx context.Context, s string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}
	if strings.ContainsAny(s, "0123456789") {
		return false, errors.New("error: msg cannot contains numerics")
	}
	return true, nil
}
