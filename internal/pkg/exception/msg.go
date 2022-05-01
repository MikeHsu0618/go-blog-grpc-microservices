package exception

type Base struct {
	InvalidToken string
	Auth         AuthMessage
}

type AuthMessage struct {
	GenerateTokenFail string
}

var Msg = Base{
	InvalidToken: "invalid token",
	Auth: AuthMessage{
		GenerateTokenFail: "generate token fail",
	},
}
