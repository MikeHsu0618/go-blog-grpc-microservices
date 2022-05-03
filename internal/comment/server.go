package comment

import (
	"context"

	v1 "blog-grpc-microservices/api/protobuf/comment/v1"
	"blog-grpc-microservices/internal/pkg/exception"
	"blog-grpc-microservices/internal/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func NewServer(logger *log.Logger, repo Repository) v1.CommentServiceServer {
	return &Server{
		logger: logger,
		repo:   repo,
	}
}

type Server struct {
	v1.UnimplementedCommentServiceServer
	logger *log.Logger
	repo   Repository
}

func (s *Server) CreateComment(ctx context.Context, in *v1.CreateCommentRequest) (*v1.CreateCommentResponse, error) {
	// 如果該評論已存在直接返回
	orgComment, err := s.repo.GetByUUID(ctx, in.GetComment().GetUuid())
	if err == nil {
		return &v1.CreateCommentResponse{
			Comment: entityToProtobuf(orgComment),
		}, nil
	}

	comment := &Comment{
		UUID:    in.GetComment().GetUuid(),
		Content: in.GetComment().GetContent(),
		PostID:  in.GetComment().GetPostId(),
		UserID:  in.GetComment().GetUserId(),
	}

	err = s.repo.Create(ctx, comment)
	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.Comment.CreateCommentFail, err)
	}

	return &v1.CreateCommentResponse{
		Comment: entityToProtobuf(comment),
	}, nil
}
func (s *Server) CreateCommentCompensate(ctx context.Context, in *v1.CreateCommentRequest) (*v1.CreateCommentResponse, error) {
	err := s.repo.DeleteByUUID(ctx, in.GetComment().GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.Comment.DeleteCommentFail, err)
	}
	return &v1.CreateCommentResponse{}, nil
}
func (s *Server) GetComment(ctx context.Context, in *v1.GetCommentRequest) (*v1.GetCommentResponse, error) {
	id := in.GetId()
	comment, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Comment.GetCommentByIDFail, err)
	}

	return &v1.GetCommentResponse{
		Comment: entityToProtobuf(comment),
	}, nil
}
func (s *Server) GetCommentByUUID(ctx context.Context, in *v1.GetCommentByUUIDRequest) (*v1.GetCommentByUUIDResponse, error) {
	comment, err := s.repo.GetByUUID(ctx, in.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Comment.GetCommentByIDFail, err)
	}

	return &v1.GetCommentByUUIDResponse{
		Comment: entityToProtobuf(comment),
	}, nil
}
func (s *Server) UpdateComment(ctx context.Context, in *v1.UpdateCommentRequest) (*v1.UpdateCommentResponse, error) {
	id := in.GetComment().GetId()
	comment, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Comment.GetCommentByIDFail, err)
	}

	if in.GetComment().GetContent() != "" {
		comment.Content = in.GetComment().GetContent()
	}

	err = s.repo.Update(ctx, comment)
	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.Comment.UpdateCommentFail, err)
	}

	return &v1.UpdateCommentResponse{
		Success: true,
	}, nil
}
func (s *Server) DeleteComment(ctx context.Context, in *v1.DeleteCommentRequest) (*v1.DeleteCommentResponse, error) {
	comment, err := s.repo.Get(ctx, in.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Comment.GetCommentByIDFail, err)
	}

	err = s.repo.Delete(ctx, comment.ID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Comment.DeleteCommentFail, err)
	}

	return &v1.DeleteCommentResponse{
		Success: true,
	}, nil
}
func (s *Server) DeleteCommentCompensate(ctx context.Context, in *v1.DeleteCommentRequest) (*v1.DeleteCommentResponse, error) {
	comment, err := s.repo.GetWithUnscoped(ctx, in.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Comment.GetCommentByIDFail, err)
	}

	comment.DeletedAt = gorm.DeletedAt{}

	err = s.repo.UpdateWithUnscoped(ctx, comment)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Comment.UpdateCommentFail, err)
	}

	return &v1.DeleteCommentResponse{
		Success: true,
	}, nil
}
func (s *Server) DeleteCommentsByPostID(ctx context.Context, in *v1.DeleteCommentsByPostIDRequest) (*v1.DeleteCommentsByPostIDResponse, error) {
	err := s.repo.DeleteByPostID(ctx, in.GetPostId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Comment.DeleteCommentFail, err)
	}

	return &v1.DeleteCommentsByPostIDResponse{
		Success: true,
	}, nil
}
func (s Server) DeleteCommentsByPostIDCompensate(ctx context.Context, req *v1.DeleteCommentsByPostIDRequest) (*v1.DeleteCommentsByPostIDResponse, error) {
	postID := req.GetPostId()
	err := s.repo.UpdateByPostIDWithUnscoped(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.Comment.UpdateCommentFail, err)
	}
	return &v1.DeleteCommentsByPostIDResponse{
		Success: true,
	}, nil
}
func (s *Server) ListCommentsByPostID(ctx context.Context, in *v1.ListCommentsByPostIDRequest) (*v1.ListCommentsByPostIDResponse, error) {
	postID := in.GetPostId()
	offset := in.GetOffset()
	limit := in.GetLimit()
	list, err := s.repo.ListByPostID(ctx, postID, int(offset), int(limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.Comment.ListCommentFail, err)
	}

	var comments []*v1.Comment
	for _, comment := range list {
		comments = append(comments, entityToProtobuf(comment))
	}

	total, err := s.repo.CountByPostID(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.Comment.ListCommentFail, err)
	}

	return &v1.ListCommentsByPostIDResponse{
		Comments: comments,
		Total:    total,
	}, nil
}

func entityToProtobuf(comment *Comment) *v1.Comment {
	return &v1.Comment{
		Id:        comment.ID,
		Uuid:      comment.UUID,
		Content:   comment.Content,
		PostId:    comment.PostID,
		UserId:    comment.UserID,
		CreatedAt: timestamppb.New(comment.CreatedAt),
		UpdatedAt: timestamppb.New(comment.UpdatedAt),
	}
}
