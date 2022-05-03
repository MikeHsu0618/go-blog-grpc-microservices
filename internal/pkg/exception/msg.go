package exception

type Base struct {
	InvalidToken string
	Auth         AuthMessage
	User         UserMessage
	Post         PostMessage
	Comment      CommentMessage
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

type PostMessage struct {
	GetPostByIDFail string
	CreatePostFail  string
	UpdatePostFail  string
	DeletePostFail  string
	ListPostFail    string
}

type CommentMessage struct {
	GetCommentByIDFail    string
	ListCommentFail       string
	DeleteCommentFail     string
	UpdateCommentFail     string
	CreateCommentFail     string
	GetCommentByPostDFail string
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
	Post: PostMessage{
		GetPostByIDFail: "get post by id fail: %v",
		CreatePostFail:  "create post fail: &v",
		UpdatePostFail:  "update post fail: &v",
		DeletePostFail:  "delete post fail: &v",
		ListPostFail:    "list post fail: &v",
	},
	Comment: CommentMessage{
		GetCommentByIDFail:    "get comment by id fail: %v",
		GetCommentByPostDFail: "get comment by post id fail: %v",
		CreateCommentFail:     "create comment fail: &v",
		UpdateCommentFail:     "update comment fail: &v",
		DeleteCommentFail:     "delete comment fail: &v",
		ListCommentFail:       "list comment fail: &v",
	},
}
