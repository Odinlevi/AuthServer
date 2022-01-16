package credential

type ICredential interface {
	GetLogin() string
	GetEmail() string
	GetPassword() string
	SetLogin(string)
	SetEmail(string)
	SetPassword(string)
}
