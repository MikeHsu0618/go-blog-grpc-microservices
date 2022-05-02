package exception

type Base struct {
	InvalidToken string
	Auth         AuthMessage
	User         UserMessage
}

type AuthMessage struct {
	GenerateTokenFail string
}

type UserMessage struct {
	GetUserListByIDsFail  string
	GetUserByIDFail       string
	GetUserByEmailFail    string
	IncorrectPassword     string
	GetUserByUsernameFail string
	GeneratePasswordFail  string
	CreateUserFail        string
	UpdateUserFail        string
	DeleteUserFail        string
}

var Msg = Base{
	InvalidToken: "invalid token",
	Auth: AuthMessage{
		GenerateTokenFail: "generate token fail",
	},
	User: UserMessage{
		GetUserListByIDsFail:  "get user list by ids fail: %v",
		GetUserByIDFail:       "get user by id fail: %v",
		GetUserByEmailFail:    "get user by email fail: %v",
		GetUserByUsernameFail: "get user by username fail: %v",
		CreateUserFail:        "create user fail: &v",
		UpdateUserFail:        "update user fail: &v",
		DeleteUserFail:        "delete user fail: &v",
		IncorrectPassword:     "incorrect password",
		GeneratePasswordFail:  "generate bcrypt password fail: %v",
	},
}
