package users

import (
	"context"

	domain "github.com/filipe1309/agl-go-driver/internal/users"
	pb "github.com/filipe1309/agl-go-driver/proto/v1/users"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *UserService) Create(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	user, err := domain.New(in.Id, in.Name, in.Login, in.Password)
	if err != nil {
		return &pb.UserResponse{Error: err.Error()}, err
	}

	id, err := s.repo.InsertDB(user)
	if err != nil {
		return &pb.UserResponse{Error: err.Error()}, err
	}

	ur := &pb.UserResponse{
		User: &pb.User{
			Id:        id,
			Name:      user.Name,
			Login:     user.Login,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}

	return ur, nil
}
