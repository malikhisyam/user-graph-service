package util

import (
	"context"
	"errors"

	"github.com/malikhisyam/user-graph-service/domains/users/models/dto"
)

func GetAuthUser(ctx context.Context) (*dto.AuthUserDto, error) {
	userRaw := ctx.Value("user")
	user, ok := userRaw.(*dto.AuthUserDto)
	if !ok || user == nil {
		return nil, errors.New("unauthorized: user not found in context")
	}
	return user, nil
}