package blog

import (
	"context"

	authv1 "blog-grpc-microservices/api/protobuf/auth/v1"
	v1 "blog-grpc-microservices/api/protobuf/blog/v1"
	commentv1 "blog-grpc-microservices/api/protobuf/comment/v1"
	postv1 "blog-grpc-microservices/api/protobuf/post/v1"
	userv1 "blog-grpc-microservices/api/protobuf/user/v1"
	"blog-grpc-microservices/internal/pkg/config"
	"blog-grpc-microservices/internal/pkg/exception"
	"blog-grpc-microservices/internal/pkg/interceptor"
	"blog-grpc-microservices/internal/pkg/log"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var prefix = "/" + v1.BlogService_ServiceDesc.ServiceName + "/"

var AuthMethods = map[string]bool{
	prefix + "SignUp":               false, // 不需要验证
	prefix + "SignIn":               false,
	prefix + "CreatePost":           true, // 需要验证
	prefix + "UpdatePost":           true,
	prefix + "GetPost":              false,
	prefix + "ListPosts":            false,
	prefix + "DeletePost":           true,
	prefix + "CreateComment":        true,
	prefix + "UpdateComment":        true,
	prefix + "DeleteComment":        true,
	prefix + "ListCommentsByPostID": false,
}

func NewServer(logger *log.Logger,
	conf *config.Config,
	userClient userv1.UserServiceClient,
	postClient postv1.PostServiceClient,
	commentClient commentv1.CommentServiceClient,
	authClient authv1.AuthServiceClient,
) v1.BlogServiceServer {
	return &Server{
		logger:        logger,
		conf:          conf,
		userClient:    userClient,
		postClient:    postClient,
		commentClient: commentClient,
		authClient:    authClient,
	}
}

type Server struct {
	v1.UnimplementedBlogServiceServer
	logger        *log.Logger
	conf          *config.Config
	userClient    userv1.UserServiceClient
	postClient    postv1.PostServiceClient
	commentClient commentv1.CommentServiceClient
	authClient    authv1.AuthServiceClient
}

func (s *Server) SignUp(ctx context.Context, in *v1.Blog_SignUpRequest) (*v1.Blog_SignUpResponse, error) {
	username := in.GetUsername()
	email := in.GetEmail()
	password := in.GetPassword()

	usernameResp, err := s.userClient.GetUserByUsername(ctx, &userv1.GetUserByUsernameRequest{
		Username: username,
	})
	if err == nil && usernameResp.GetUser().GetId() != 0 {
		return nil, status.Error(codes.AlreadyExists, exception.Msg.Blog.UsernameAlreadyExists)
	}
	emailResp, err := s.userClient.GetUserByEmail(ctx, &userv1.GetUserByEmailRequest{
		Email: email,
	})
	if err == nil && emailResp.GetUser().GetId() != 0 {
		return nil, status.Error(codes.AlreadyExists, exception.Msg.Blog.EmailAlreadyExists)
	}

	userResp, err := s.userClient.CreateUser(ctx, &userv1.CreateUserRequest{User: &userv1.User{
		Username: username,
		Email:    email,
		Password: password,
	}})

	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.Blog.CreateUserFail)
	}

	authResp, err := s.authClient.GenerateToken(ctx, &authv1.GenerateTokenRequest{
		UserId: userResp.GetUser().GetId(),
	})
	if err != nil {
		s.logger.Error(err)
		return nil, status.Error(codes.Internal, exception.Msg.Blog.GenerateTokenFail)
	}

	return &v1.Blog_SignUpResponse{
		Token: authResp.GetToken(),
	}, nil
}
func (s *Server) SignIn(ctx context.Context, in *v1.Blog_SignInRequest) (*v1.Blog_SignInResponse, error) {
	username := in.GetUsername()
	email := in.GetEmail()
	password := in.GetPassword()

	var userID uint64
	if email != "" {
		resp, err := s.userClient.GetUserByEmail(ctx, &userv1.GetUserByEmailRequest{
			Email:    email,
			Password: password,
		})
		if err != nil {
			s.logger.Error(err)
			return nil, status.Errorf(codes.Internal, exception.Msg.Blog.GetUserByEmailFail)
		}

		userID = resp.GetUser().GetId()
	} else {
		resp, err := s.userClient.GetUserByUsername(ctx, &userv1.GetUserByUsernameRequest{
			Username: username,
			Password: password,
		})
		if err != nil {
			s.logger.Error(err)
			return nil, status.Errorf(codes.Internal, exception.Msg.Blog.GetUserByUsernameFail)
		}
		userID = resp.GetUser().GetId()
	}

	authResp, err := s.authClient.GenerateToken(ctx, &authv1.GenerateTokenRequest{
		UserId: userID,
	})
	if err != nil {
		s.logger.Error(err)
		return nil, status.Error(codes.Internal, exception.Msg.Blog.GenerateTokenFail)
	}

	return &v1.Blog_SignInResponse{
		Token: authResp.GetToken(),
	}, nil
}
func (s *Server) CreatePost(ctx context.Context, in *v1.Blog_CreatePostRequest) (*v1.Blog_CreatePostResponse, error) {
	userID, ok := ctx.Value(interceptor.ContextKeyID).(uint64)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, exception.Msg.Blog.UserNotAuthenticated)
	}
	userResp, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: userID})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Blog.GetUserByIDFail)
	}
	postResp, err := s.postClient.CreatePost(ctx, &postv1.CreatePostRequest{Post: &postv1.Post{
		Uuid:    uuid.New().String(),
		Title:   in.GetPost().GetTitle(),
		Content: in.GetPost().GetContent(),
		UserId:  userResp.GetUser().GetId(),
	}})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Blog.CreatePostFail)
	}

	return &v1.Blog_CreatePostResponse{
		Post: &v1.Blog_Post{
			Id:            postResp.GetPost().GetId(),
			Title:         postResp.GetPost().GetTitle(),
			Content:       postResp.GetPost().GetContent(),
			UserId:        postResp.GetPost().GetUserId(),
			CommentsCount: postResp.GetPost().GetCommentsCount(),
			CreatedAt:     postResp.GetPost().GetCreatedAt(),
			UpdatedAt:     postResp.GetPost().GetUpdatedAt(),
			User: &v1.Blog_User{
				Id:       userResp.GetUser().GetId(),
				Username: userResp.GetUser().GetUsername(),
				Avatar:   userResp.GetUser().GetAvatar(),
			},
		},
	}, nil
}
func (s *Server) GetPost(ctx context.Context, in *v1.Blog_GetPostRequest) (*v1.Blog_GetPostResponse, error) {
	postResp, err := s.postClient.GetPost(ctx, &postv1.GetPostRequest{Id: in.GetId()})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Blog.GetPostByIDFail)
	}
	postuserResp, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: postResp.GetPost().UserId})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Blog.GetUserByIDFail)
	}

	return &v1.Blog_GetPostResponse{Post: &v1.Blog_Post{
		Id:            postResp.GetPost().GetId(),
		Title:         postResp.GetPost().GetTitle(),
		Content:       postResp.GetPost().GetContent(),
		UserId:        postResp.GetPost().GetUserId(),
		CommentsCount: postResp.GetPost().GetCommentsCount(),
		CreatedAt:     postResp.GetPost().GetCreatedAt(),
		UpdatedAt:     postResp.GetPost().GetUpdatedAt(),
		User: &v1.Blog_User{
			Id:       postuserResp.GetUser().GetId(),
			Username: postuserResp.GetUser().GetUsername(),
			Avatar:   postuserResp.GetUser().GetAvatar(),
		},
	}}, nil
}
func (s *Server) ListPosts(ctx context.Context, in *v1.Blog_ListPostsRequest) (*v1.Blog_ListPostsResponse, error) {
	postResp, err := s.postClient.ListPosts(ctx, &postv1.ListPostsRequest{
		Limit:  int32(in.GetLimit()),
		Offset: int32(in.GetOffset()),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.Blog.ListPostFail)
	}

	var posts []*v1.Blog_Post
	var postUserIDs []uint64

	for _, post := range postResp.GetPosts() {
		postUserIDs = append(postUserIDs, post.GetUserId())
	}

	postUserResp, err := s.userClient.ListUsersByIDs(ctx, &userv1.ListUsersByIDsRequest{Ids: postUserIDs})
	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.Blog.GetUserByIDFail)
	}

	for _, post := range posts {
		for _, postUser := range postUserResp.GetUsers() {
			if post.GetUserId() == postUser.GetId() {
				posts = append(posts, &v1.Blog_Post{
					Id:            post.GetId(),
					Title:         post.GetTitle(),
					Content:       post.GetTitle(),
					UserId:        post.GetUserId(),
					CommentsCount: post.GetCommentsCount(),
					CreatedAt:     post.GetCreatedAt(),
					UpdatedAt:     post.GetUpdatedAt(),
					User: &v1.Blog_User{
						Id:       postUser.GetId(),
						Username: postUser.GetUsername(),
						Avatar:   postUser.GetAvatar(),
					},
				})
			}
		}
	}

	return &v1.Blog_ListPostsResponse{
		Posts: posts,
		Total: postResp.GetCount(),
	}, nil
}
func (s *Server) UpdatePost(ctx context.Context, in *v1.Blog_UpdatePostRequest) (*v1.Blog_UpdatePostResponse, error) {
	userID, ok := ctx.Value(interceptor.ContextKeyID).(uint64)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, exception.Msg.Blog.UserNotAuthenticated)
	}
	userResp, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: userID})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Blog.GetUserByIDFail)
	}
	postResp, err := s.postClient.GetPost(ctx, &postv1.GetPostRequest{Id: in.GetPost().GetId()})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Blog.GetPostByIDFail)
	}
	// 只有自己可以修改自己的文章
	if userResp.GetUser().GetId() != postResp.GetPost().GetUserId() {
		return nil, status.Errorf(codes.Unauthenticated, exception.Msg.Blog.UserNotAuthenticated)
	}

	// 更新文章
	updatedPost := &postv1.Post{
		Id: in.GetPost().GetId(),
	}
	if in.GetPost().GetTitle() != "" {
		updatedPost.Title = in.GetPost().GetTitle()
	}
	if in.GetPost().GetContent() != "" {
		updatedPost.Content = in.GetPost().GetContent()
	}
	udpatePostResp, err := s.postClient.UpdatePost(ctx, &postv1.UpdatePostRequest{Post: updatedPost})
	if err != nil || udpatePostResp.GetSuccess() {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Blog.UpdatePostFail)
	}

	return &v1.Blog_UpdatePostResponse{
		Success: true,
	}, nil
}
func (s *Server) DeletePost(ctx context.Context, in *v1.Blog_DeletePostRequest) (*v1.Blog_DeletePostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePost not implemented")
}
func (s *Server) CreateComment(ctx context.Context, in *v1.Blog_CreateCommentRequest) (*v1.Blog_CreateCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateComment not implemented")
}
func (s *Server) DeleteComment(ctx context.Context, in *v1.Blog_DeleteCommentRequest) (*v1.Blog_DeleteCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteComment not implemented")
}
func (s *Server) UpdateComment(ctx context.Context, in *v1.Blog_UpdateCommentRequest) (*v1.Blog_UpdateCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateComment not implemented")
}
func (s *Server) ListCommentsByPostID(ctx context.Context, in *v1.Blog_ListCommentsByPostIDRequest) (*v1.Blog_ListCommentsByPostIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCommentsByPostID not implemented")
}
