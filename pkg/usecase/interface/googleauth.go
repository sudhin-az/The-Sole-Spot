package interfaces

import (
	"context"
	"ecommerce_clean_architecture/pkg/domain"
)

type AuthUseCase interface {
	HandleGoogleCallback(ctx context.Context, code string) (domain.Users, string, error)
}
