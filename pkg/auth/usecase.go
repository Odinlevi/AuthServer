package auth

import (
	"context"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/models"
)

type UseCase interface {
	SignUp(ctx context.Context, user *models.User) error
	SignIn(ctx context.Context, user *models.User) (string, error)

	Send(ctx context.Context, message *models.Message) error
	Get(ctx context.Context, userID string, senderID string) ([]models.Message, error)
	AddCompanion(ctx context.Context, userID string, companionID string) error
	RemoveCompanion(ctx context.Context, userID string, companionID string) error
	GetCompanions(ctx context.Context, userID string) ([]string, error)
}
