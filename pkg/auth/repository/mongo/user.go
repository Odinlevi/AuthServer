package mongo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth/credential"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db            *mongo.Collection
	confirmations map[credential.ICredential]uuid.UUID
}

func NewUserRepository(db *mongo.Database, collection string) *UserRepository {
	return &UserRepository{
		db: db.Collection(collection),
	}
}

func (r *UserRepository) Insert(ctx context.Context, user credential.ICredential) error {

	newUser := new(models.User)
	newUser.Username = user.GetLogin()
	newUser.Email = user.GetEmail()
	newUser.Password = user.GetPassword()

	var emailTest *models.User = nil

	err := r.db.FindOne(ctx, bson.M{"email": newUser.Email}).Decode(emailTest)
	if err != mongo.ErrNoDocuments {
		errEmailExist := errors.New("user with this email already exist")
		log.Errorf("%s", errEmailExist.Error())
		return errEmailExist
	}

	_, err = r.db.InsertOne(ctx, newUser)
	if err != nil {
		log.Errorf("error on inserting user: %s", err.Error())
		return auth.ErrUserAlreadyExists
	}

	return nil

	//_, err := r.db.InsertOne(ctx, user)
	//if err != nil {
	//	log.Errorf("error on inserting user: %s", err.Error())
	//	return auth.ErrUserAlreadyExists
	//}
	//
	//return nil
}

func (r *UserRepository) Get(ctx context.Context, username, password string) (*models.User, error) {
	user := new(models.User)

	if err := r.db.FindOne(ctx, bson.M{"_id": username, "password": password}).Decode(user); err != nil {
		log.Errorf("error occured while getting user from db: %s", err.Error())
		if err == mongo.ErrNoDocuments {
			return nil, auth.ErrUserDoesNotExist
		}

		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetConfirmationToken(user credential.ICredential) uuid.UUID {
	token, ok := r.confirmations[user]
	if !ok {
		token = uuid.New()
		r.confirmations[user] = token
	}

	return token
}

func (r *UserRepository) SetConfirmationToken(token string) (credential.ICredential, bool) {
	var user credential.ICredential = nil
	for key, value := range r.confirmations {
		if value.String() == token {
			user = key
		}
	}

	if user == nil {
		return nil, false
	}

	defer delete(r.confirmations, user)

	return user, true
}
