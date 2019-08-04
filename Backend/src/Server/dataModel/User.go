package dataModel

import "github.com/badoux/checkmail"

type UserDAO interface {
	CreateUser(email, salt, hash string) (int, error)
	IsEmailInUse(email string) (bool, error)
	GetAuthenticationInformation(email string) (*UserAuthenticationDTO, error)
	GetUserProfile(userId int) (*User, error)
	GetUserWithPermissions(leagueId, userId int) (*UserWithPermissions, error)
}

type UserCreationInformation struct {
	Email    string
	Password string
}

type UserAuthenticationDTO struct {
	UserId int    `json:"userId"`
	Salt   string `json:"salt"`
	Hash   string `json:"hash"`
}

type User struct {
	UserId int    `json:"userId"`
	Email  string `json:"email"`
}

type UserWithPermissions struct {
	UserId            int                    `json:"userId"`
	Email             string                 `json:"email"`
	LeaguePermissions *LeaguePermissionsCore `json:"leaguePermissions"`
	TeamPermissions   []*TeamPermissions     `json:"teamPermissions"`
}

type TeamManager struct {
	UserId        int    `json:"userId"`
	Email         string `json:"email"`
	Administrator bool   `json:"administrator"`
	Information   bool   `json:"information"`
	Games         bool   `json:"games"`
}

// UserCreationInformation
func (user *UserCreationInformation) Validate(userDao UserDAO) (bool, string, error) {
	return validate(
		validateEmailFormat(user.Email),
		user.uniqueness(userDao),
		validatePasswordLength(user.Password))
}

func (user *UserCreationInformation) email() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return false
		} else {
			return true
		}
	}
}

func (user *UserCreationInformation) uniqueness(userDao UserDAO) ValidateFunc {
	return func(problemDest *string, errorDest *error) bool {
		valid := false
		inUse, err := userDao.IsEmailInUse(user.Email)
		if err != nil {
			*errorDest = err
		} else if inUse {
			*problemDest = EmailInUse
		} else {
			valid = true
		}
		return valid
	}
}

func (user *UserCreationInformation) password() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		valid := false
		if len(user.Password) > MaxPasswordLength {
			*problemDest = PasswordTooLong
		} else if len(user.Password) < MinPasswordLength {
			*problemDest = PasswordTooShort
		} else {
			valid = true
		}
		return valid
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req *LoginRequest) Validate() (bool, string, error) {
	return validate(
		validateEmailFormat(req.Email),
		validatePasswordLength(req.Password))
}
