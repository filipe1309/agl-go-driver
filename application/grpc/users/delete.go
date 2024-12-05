package users

import (
	"context"

	pb "github.com/filipe1309/agl-go-driver/proto/v1/users"
)

func (s *UserService) Delete(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	err := s.repo.SoftDeleteDB(in.Id)
	if err != nil {
		return &pb.UserResponse{Error: err.Error()}, err
	}

	return &pb.UserResponse{}, nil
}
