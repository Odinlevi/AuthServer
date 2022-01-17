package auth

import (
	"context"
	"github.com/google/uuid"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth/credential"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/models"
)

type Repository interface {
	Insert(ctx context.Context, user credential.ICredential) error
	Get(ctx context.Context, username, password string) (*models.User, error)

	GetConfirmationToken(user credential.ICredential) uuid.UUID
	SetConfirmationToken(token string) (credential.ICredential, bool)
}
