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
	"github.com/dtm-labs/dtm/dtmgrpc"
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
		return nil, status.Error(codes.Internal, err.Error())
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
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		userID = resp.GetUser().GetId()
	} else {
		resp, err := s.userClient.GetUserByUsername(ctx, &userv1.GetUserByUsernameRequest{
			Username: username,
			Password: password,
		})
		if err != nil {
			s.logger.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		userID = resp.GetUser().GetId()
	}

	authResp, err := s.authClient.GenerateToken(ctx, &authv1.GenerateTokenRequest{
		UserId: userID,
	})
	if err != nil {
		s.logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
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
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	postResp, err := s.postClient.CreatePost(ctx, &postv1.CreatePostRequest{Post: &postv1.Post{
		Uuid:    uuid.New().String(),
		Title:   in.GetPost().GetTitle(),
		Content: in.GetPost().GetContent(),
		UserId:  userResp.GetUser().GetId(),
	}})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
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
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	postuserResp, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: postResp.GetPost().UserId})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
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
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var posts []*v1.Blog_Post
	var postUserIDs []uint64

	for _, post := range postResp.GetPosts() {
		postUserIDs = append(postUserIDs, post.GetUserId())
	}

	postUserResp, err := s.userClient.ListUsersByIDs(ctx, &userv1.ListUsersByIDsRequest{Ids: postUserIDs})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	for _, post := range postResp.GetPosts() {
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
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	postResp, err := s.postClient.GetPost(ctx, &postv1.GetPostRequest{Id: in.GetPost().GetId()})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	// 只有自己可以修改自己的文章
	if userResp.GetUser().GetId() != postResp.GetPost().GetUserId() {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
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
	updatePostResp, err := s.postClient.UpdatePost(ctx, &postv1.UpdatePostRequest{Post: updatedPost})
	if err != nil || !updatePostResp.GetSuccess() {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &v1.Blog_UpdatePostResponse{
		Success: true,
	}, nil
}
func (s *Server) DeletePost(ctx context.Context, in *v1.Blog_DeletePostRequest) (*v1.Blog_DeletePostResponse, error) {
	userID, ok := ctx.Value(interceptor.ContextKeyID).(uint64)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, exception.Msg.Blog.UserNotAuthenticated)
	}
	userResp, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: userID})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	postResp, err := s.postClient.GetPost(ctx, &postv1.GetPostRequest{Id: in.GetId()})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// 只能刪除自己的文章
	if userResp.GetUser().GetId() != postResp.GetPost().GetUserId() {
		return nil, status.Error(codes.PermissionDenied, exception.Msg.Blog.UserNotAuthenticated)
	}

	// 分佈式事務(Saga 模式)
	dtmGRPCServerAddr := s.conf.DTM.Server.Host + s.conf.DTM.Server.GRPC.Port
	gid := dtmgrpc.MustGenGid(dtmGRPCServerAddr)
	s.logger.Info("gid:", gid)
	saga := dtmgrpc.NewSagaGrpc(dtmGRPCServerAddr, gid).Add(
		s.conf.Post.Server.Host+s.conf.Post.Server.GRPC.Port+"/"+postv1.PostService_ServiceDesc.ServiceName+"/DeletePost",
		s.conf.Post.Server.Host+s.conf.Post.Server.GRPC.Port+"/"+postv1.PostService_ServiceDesc.ServiceName+"/DeletePostCompensate",
		&postv1.DeletePostRequest{
			Id: in.GetId(),
		},
	).Add(
		s.conf.Comment.Server.Host+s.conf.Comment.Server.GRPC.Port+"/"+commentv1.CommentService_ServiceDesc.ServiceName+"/DeleteCommentsByPostID",
		s.conf.Comment.Server.Host+s.conf.Comment.Server.GRPC.Port+"/"+commentv1.CommentService_ServiceDesc.ServiceName+"/DeleteCommentsByPostIDCompensate",
		&commentv1.DeleteCommentsByPostIDRequest{
			PostId: in.GetId(),
		},
	)
	saga.WaitResult = true
	err = saga.Submit()
	if err != nil {
		s.logger.Error("saga submit error:", err)
		return nil, status.Error(codes.Internal, "saga submit failed")
	}

	return &v1.Blog_DeletePostResponse{
		Success: true,
	}, nil
}
func (s *Server) CreateComment(ctx context.Context, in *v1.Blog_CreateCommentRequest) (*v1.Blog_CreateCommentResponse, error) {
	userID, ok := ctx.Value(interceptor.ContextKeyID).(uint64)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, exception.Msg.Blog.UserNotAuthenticated)
	}
	userResp, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: userID})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	postResp, err := s.postClient.GetPost(ctx, &postv1.GetPostRequest{Id: in.GetComment().GetPostId()})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	comment := &commentv1.Comment{
		Uuid:    uuid.New().String(),
		Content: in.GetComment().GetContent(),
		PostId:  postResp.GetPost().GetId(),
		UserId:  userResp.GetUser().GetId(),
	}

	// 分佈式事務(Saga 模式)
	dtmGRPCServerAddr := s.conf.DTM.Server.Host + s.conf.DTM.Server.GRPC.Port
	gid := dtmgrpc.MustGenGid(dtmGRPCServerAddr)
	s.logger.Info("gid:", gid)
	saga := dtmgrpc.NewSagaGrpc(dtmGRPCServerAddr, gid).Add(
		s.conf.Comment.Server.Host+s.conf.Comment.Server.GRPC.Port+"/"+commentv1.CommentService_ServiceDesc.ServiceName+"/CreateComment",
		s.conf.Post.Server.Host+s.conf.Post.Server.GRPC.Port+"/"+commentv1.CommentService_ServiceDesc.ServiceName+"/CreateCommentCompensate",
		&commentv1.CreateCommentRequest{
			Comment: comment,
		},
	).Add(
		s.conf.Post.Server.Host+s.conf.Post.Server.GRPC.Port+"/"+postv1.PostService_ServiceDesc.ServiceName+"/IncrementCommentsCount",
		s.conf.Post.Server.Host+s.conf.Post.Server.GRPC.Port+"/"+postv1.PostService_ServiceDesc.ServiceName+"/IncrementCommentsCountCompensate",
		&postv1.IncrementCommentsCountRequest{
			Id: postResp.GetPost().GetId(),
		},
	)
	saga.WaitResult = true
	err = saga.Submit()
	if err != nil {
		s.logger.Error("saga submit error:", err)
		return nil, status.Error(codes.Internal, "saga submit failed")
	}
	postUserResp, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{
		Id: postResp.GetPost().GetUserId(),
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	commentResp, err := s.commentClient.GetCommentByUUID(ctx, &commentv1.GetCommentByUUIDRequest{
		Uuid: comment.GetUuid(),
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.Blog_CreateCommentResponse{
		Comment: &v1.Blog_Comment{
			Id:        commentResp.GetComment().GetId(),
			Content:   commentResp.GetComment().GetContent(),
			PostId:    commentResp.GetComment().GetPostId(),
			UserId:    commentResp.GetComment().GetUserId(),
			CreatedAt: commentResp.GetComment().GetCreatedAt(),
			UpdatedAt: commentResp.GetComment().GetUpdatedAt(),
			Post: &v1.Blog_Post{
				Id:            postResp.GetPost().GetId(),
				Title:         postResp.GetPost().GetTitle(),
				Content:       postResp.GetPost().GetContent(),
				UserId:        postResp.GetPost().GetUserId(),
				CommentsCount: postResp.GetPost().GetCommentsCount(),
				CreatedAt:     postResp.GetPost().GetCreatedAt(),
				UpdatedAt:     postResp.GetPost().GetUpdatedAt(),
				User: &v1.Blog_User{
					Id:       postUserResp.GetUser().GetId(),
					Username: postUserResp.GetUser().GetUsername(),
					Avatar:   postUserResp.GetUser().GetAvatar(),
				},
			},
			User: &v1.Blog_User{
				Id:       userResp.GetUser().GetId(),
				Username: userResp.GetUser().GetUsername(),
				Avatar:   userResp.GetUser().GetAvatar(),
			},
		},
	}, nil
}
func (s *Server) DeleteComment(ctx context.Context, in *v1.Blog_DeleteCommentRequest) (*v1.Blog_DeleteCommentResponse, error) {
	userID, ok := ctx.Value(interceptor.ContextKeyID).(uint64)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, exception.Msg.Blog.UserNotAuthenticated)
	}
	userResp, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: userID})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	commentResp, err := s.commentClient.GetComment(ctx, &commentv1.GetCommentRequest{Id: in.GetId()})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	postResp, err := s.postClient.GetPost(ctx, &postv1.GetPostRequest{
		Id: commentResp.GetComment().GetPostId(),
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	// 只能更新自己的評論
	if commentResp.GetComment().GetUserId() != userResp.GetUser().GetId() || userResp.GetUser().GetId() != postResp.GetPost().GetUserId() {
		return nil, status.Error(codes.PermissionDenied, exception.Msg.Blog.UserNotAuthenticated)
	}

	// 分佈式(Saga 模式): 刪除評論 and 减少評論數目
	dtmGRPCServerAddr := s.conf.DTM.Server.Host + s.conf.DTM.Server.GRPC.Port
	gid := dtmgrpc.MustGenGid(dtmGRPCServerAddr)
	s.logger.Info("gid:", gid)
	saga := dtmgrpc.NewSagaGrpc(dtmGRPCServerAddr, gid).Add(
		s.conf.Comment.Server.Host+s.conf.Comment.Server.GRPC.Port+"/"+commentv1.CommentService_ServiceDesc.ServiceName+"/DeleteComment",
		s.conf.Comment.Server.Host+s.conf.Comment.Server.GRPC.Port+"/"+commentv1.CommentService_ServiceDesc.ServiceName+"/DeleteCommentCompensate",
		&commentv1.DeleteCommentRequest{
			Id: in.GetId(),
		},
	).Add(
		s.conf.Post.Server.Host+s.conf.Post.Server.GRPC.Port+"/"+postv1.PostService_ServiceDesc.ServiceName+"/DecrementCommentsCount",
		s.conf.Post.Server.Host+s.conf.Post.Server.GRPC.Port+"/"+postv1.PostService_ServiceDesc.ServiceName+"/DecrementCommentsCountCompensate",
		&postv1.DecrementCommentsCountRequest{
			Id: postResp.GetPost().GetId(),
		},
	)

	saga.WaitResult = true
	err = saga.Submit()
	if err != nil {
		return nil, status.Error(codes.Internal, "saga submit failed")
	}

	return &v1.Blog_DeleteCommentResponse{Success: true}, nil
}
func (s *Server) UpdateComment(ctx context.Context, in *v1.Blog_UpdateCommentRequest) (*v1.Blog_UpdateCommentResponse, error) {
	userID, ok := ctx.Value(interceptor.ContextKeyID).(uint64)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, exception.Msg.Blog.UserNotAuthenticated)
	}
	userResp, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: userID})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	commentResp, err := s.commentClient.GetComment(ctx, &commentv1.GetCommentRequest{Id: in.GetComment().Id})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	// 只能更新自己的評論
	if commentResp.GetComment().GetUserId() != userResp.GetUser().GetId() {
		return nil, status.Error(codes.PermissionDenied, "user not authorized")
	}

	comment := &commentv1.Comment{
		Id:      in.GetComment().GetId(),
		Content: in.GetComment().GetContent(),
	}
	_, err = s.commentClient.UpdateComment(ctx, &commentv1.UpdateCommentRequest{Comment: comment})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &v1.Blog_UpdateCommentResponse{
		Success: true,
	}, nil
}
func (s *Server) ListCommentsByPostID(ctx context.Context, in *v1.Blog_ListCommentsByPostIDRequest) (*v1.Blog_ListCommentsByPostIDResponse, error) {
	postID := in.GetPostId()
	offset := in.GetOffset()
	limit := in.GetLimit()
	commentResp, err := s.commentClient.ListCommentsByPostID(ctx, &commentv1.ListCommentsByPostIDRequest{
		PostId: postID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.Blog.ListPostFail)
	}
	var commentUserIDs []uint64
	for _, post := range commentResp.GetComments() {
		commentUserIDs = append(commentUserIDs, post.GetUserId())
	}

	commentUserResp, err := s.userClient.ListUsersByIDs(ctx, &userv1.ListUsersByIDsRequest{
		Ids: commentUserIDs,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var comments []*v1.Blog_Comment
	for _, comment := range commentResp.GetComments() {
		for _, user := range commentUserResp.GetUsers() {
			if user.GetId() == comment.GetUserId() {
				comments = append(comments, &v1.Blog_Comment{
					Id:        comment.GetId(),
					Content:   comment.GetContent(),
					PostId:    comment.GetPostId(),
					UserId:    comment.GetUserId(),
					CreatedAt: comment.GetCreatedAt(),
					UpdatedAt: comment.GetUpdatedAt(),
					User: &v1.Blog_User{
						Id:       user.GetId(),
						Username: user.GetUsername(),
						Avatar:   user.GetAvatar(),
					},
				})
			}
		}
	}

	return &v1.Blog_ListCommentsByPostIDResponse{
		Comments: comments,
		Total:    commentResp.GetTotal(),
	}, nil
}
