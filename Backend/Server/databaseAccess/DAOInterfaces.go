package databaseAccess

type UsersDAO interface {
	InsertUser(email, salt, hash string) error
	IsEmailInUse(email string) (bool, error)
	GetAuthenticationInformation(email string) (int, string, string, error)
}