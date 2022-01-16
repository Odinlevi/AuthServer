package usecase

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/spf13/viper"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth/repository/mongo"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/models"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/rsamail"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/sendsmtp"
	"time"
)

type Authorizer struct {
	userRepo    auth.Repository
	chatRepo    *mongo.ChatRepository
	messageRepo *mongo.MessageRepository

	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
	sender         sendsmtp.ISender
}

func NewAuthorizer(userRepo auth.Repository, chatRepo *mongo.ChatRepository, messageRepo *mongo.MessageRepository,
	hashSalt string, signingKey []byte, expireDuration time.Duration) *Authorizer {
	sender := sendsmtp.NewSender(
		viper.GetString("email.from"),
		viper.GetString("email.password"),
		viper.GetString("email.host"),
		viper.GetString("email.port"))

	return &Authorizer{
		userRepo:       userRepo,
		chatRepo:       chatRepo,
		messageRepo:    messageRepo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: expireDuration,
		sender:         sender,
	}
}

func (a *Authorizer) SignUp(ctx context.Context, user *models.User) error {
	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(a.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	result := a.userRepo.Insert(ctx, user)
	if result != nil {
		return result
	}

	publicKeyPath, privateKeyPath := viper.GetString("rsamail.public_key"), viper.GetString("rsamail.private_key")
	privateKey := rsamail.GenerateKeyPair(publicKeyPath, privateKeyPath)

	user.SetEmail(rsamail.Encrypt(privateKey, user.GetEmail()))

	id := a.userRepo.GetConfirmationToken(user)
	if user.Email != "" {
		a.sender.Send(rsamail.Decrypt(privateKey, user.GetEmail()),
			"Sign Up",
			"Hello, "+user.GetLogin()+"!"+"\n"+
				"Thank you for registering with our service.\n"+
				"Confirmation: $URL/auth/confirm/"+id.String())
	}

	return result
}

func (a *Authorizer) SignIn(ctx context.Context, user *models.User) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(a.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := a.userRepo.Get(ctx, user.Username, user.Password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireDuration)),
			IssuedAt:  jwt.At(time.Now()),
		},
		Username: user.Username,
	})

	if user.Email != "" {
		a.sender.Send(user.Email,
			"Sign In",
			"Whats up fucker?!")
	}

	return token.SignedString(a.signingKey)
}

func (a *Authorizer) Confirm(ctx context.Context, token string) error {
	user, ok := a.userRepo.SetConfirmationToken(token)
	if !ok {
		return errors.New("Token does not exist.")
	}

	return a.userRepo.Insert(ctx, user)
}

func (a *Authorizer) Send(ctx context.Context, message *models.Message) error {
	return a.messageRepo.Send(ctx, message)
}

func (a *Authorizer) Get(ctx context.Context, userID string, senderID string) ([]models.Message, error) {
	result, err := a.messageRepo.Get(ctx, userID, senderID)
	return result, err
}

func (a *Authorizer) AddCompanion(ctx context.Context, userID string, companionID string) error {
	return a.chatRepo.AddCompanion(ctx, userID, companionID)
}

func (a *Authorizer) RemoveCompanion(ctx context.Context, userID string, companionID string) error {
	return a.chatRepo.RemoveCompanion(ctx, userID, companionID)
}

func (a *Authorizer) GetCompanions(ctx context.Context, userID string) ([]string, error) {
	return a.chatRepo.GetCompanions(ctx, userID)
}
