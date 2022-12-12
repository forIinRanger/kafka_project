package validator

import "context"

type Validator interface {
	Validate(ctx context.Context, s string) (bool, error)
}
