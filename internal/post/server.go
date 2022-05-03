package post

import (
	"context"

	v1 "blog-grpc-microservices/api/protobuf/post/v1"
	"blog-grpc-microservices/internal/pkg/exception"
	"blog-grpc-microservices/internal/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func NewServer(logger *log.Logger, repo Repository) v1.PostServiceServer {
	return &Server{
		logger: logger,
		repo:   repo,
	}
}

type Server struct {
	v1.UnimplementedPostServiceServer
	logger *log.Logger
	repo   Repository
}

func (s *Server) GetPost(ctx context.Context, in *v1.GetPostRequest) (*v1.GetPostResponse, error) {
	post, err := s.repo.Get(ctx, in.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Post.GetPostByIDFail, err)
	}
	return &v1.GetPostResponse{
		Post: entityToProtobuf(post),
	}, nil
}

func (s *Server) CreatePost(ctx context.Context, in *v1.CreatePostRequest) (*v1.CreatePostResponse, error) {
	post := &Post{
		UUID:    in.GetPost().GetUuid(),
		Title:   in.GetPost().GetTitle(),
		Content: in.GetPost().GetContent(),
		UserID:  in.GetPost().GetUserId(),
	}
	err := s.repo.Create(ctx, post)
	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.Post.CreatePostFail, err)
	}

	return &v1.CreatePostResponse{
		Post: entityToProtobuf(post),
	}, nil
}
func (s *Server) UpdatePost(ctx context.Context, in *v1.UpdatePostRequest) (*v1.UpdatePostResponse, error) {
	postID := in.GetPost().GetId()
	post, err := s.repo.Get(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Post.GetPostByIDFail, err)
	}

	if in.GetPost().GetTitle() != "" {
		post.Title = in.GetPost().GetTitle()
	}

	if in.GetPost().GetContent() != "" {
		post.Content = in.GetPost().GetContent()
	}

	s.logger.Info("update post", post)

	err = s.repo.Update(ctx, post)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Post.UpdatePostFail, err)
	}

	return &v1.UpdatePostResponse{
		Success: true,
	}, nil
}

func (s *Server) DeletePost(ctx context.Context, in *v1.DeletePostRequest) (*v1.DeletePostResponse, error) {
	post, err := s.repo.Get(ctx, in.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Post.GetPostByIDFail, err)
	}

	err = s.repo.Delete(ctx, post.ID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Post.DeletePostFail, err)
	}

	return &v1.DeletePostResponse{
		Success: true,
	}, nil
}

func (s *Server) DeletePostCompensate(ctx context.Context, in *v1.DeletePostRequest) (*v1.DeletePostResponse, error) {
	post, err := s.repo.GetWithUnscoped(ctx, in.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Post.GetPostByIDFail, err)
	}

	post.DeletedAt = gorm.DeletedAt{}

	err = s.repo.UpdateWithUnscoped(ctx, post)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Post.UpdatePostFail, err)
	}

	return &v1.DeletePostResponse{
		Success: true,
	}, nil
}
func (s *Server) ListPosts(ctx context.Context, in *v1.ListPostsRequest) (*v1.ListPostsResponse, error) {
	list, err := s.repo.List(ctx, int(in.GetOffset()), int(in.GetLimit()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Post.ListPostFail, err)
	}

	var posts []*v1.Post
	for _, post := range list {
		posts = append(posts, entityToProtobuf(post))
	}

	count, err := s.repo.Count(ctx)

	return &v1.ListPostsResponse{
		Posts: posts,
		Count: count,
	}, nil
}

func (s *Server) IncrementCommentsCount(ctx context.Context, in *v1.IncrementCommentsCountRequest) (*v1.IncrementCommentsCountResponse, error) {
	id := in.GetId()
	post, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Post.GetPostByIDFail, err)
	}
	post.CommentsCount++
	err = s.repo.Update(ctx, post)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Post.UpdatePostFail, err)
	}

	return &v1.IncrementCommentsCountResponse{
		Success: true,
	}, nil
}
func (s *Server) IncrementCommentsCountCompensate(ctx context.Context, in *v1.IncrementCommentsCountRequest) (*v1.IncrementCommentsCountResponse, error) {
	id := in.GetId()
	post, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Post.GetPostByIDFail, err)
	}
	post.CommentsCount--
	err = s.repo.Update(ctx, post)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Post.UpdatePostFail, err)
	}

	return &v1.IncrementCommentsCountResponse{
		Success: true,
	}, nil
}

func (s *Server) DecrementCommentsCount(ctx context.Context, in *v1.DecrementCommentsCountRequest) (*v1.DecrementCommentsCountResponse, error) {
	id := in.GetId()
	post, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Post.GetPostByIDFail, err)
	}
	post.CommentsCount--
	err = s.repo.Update(ctx, post)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Post.UpdatePostFail, err)
	}

	return &v1.DecrementCommentsCountResponse{
		Success: true,
	}, nil
}

func (s *Server) DecrementCommentsCountCompensate(ctx context.Context, in *v1.DecrementCommentsCountRequest) (*v1.DecrementCommentsCountResponse, error) {
	id := in.GetId()
	post, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Post.GetPostByIDFail, err)
	}
	post.CommentsCount++
	err = s.repo.Update(ctx, post)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Post.UpdatePostFail, err)
	}

	return &v1.DecrementCommentsCountResponse{
		Success: true,
	}, nil
}

func entityToProtobuf(post *Post) *v1.Post {
	return &v1.Post{
		Id:            post.ID,
		Uuid:          post.UUID,
		Title:         post.Title,
		Content:       post.Content,
		CommentsCount: post.CommentsCount,
		UserId:        post.UserID,
		CreatedAt:     timestamppb.New(post.CreatedAt),
		UpdatedAt:     timestamppb.New(post.UpdatedAt),
	}
}
