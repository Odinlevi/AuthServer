package models

type User struct {
	Username string `json:"username" bson:"_id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func (u *User) GetLogin() string {
	return u.Username
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetPassword() string {
	return u.Password
}

func (u *User) SetLogin(login string) {
	u.Username = login
}

func (u *User) SetEmail(email string) {
	u.Email = email
}

func (u *User) SetPassword(password string) {
	u.Password = password
}
