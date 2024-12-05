package users

import (
	"context"

	pb "github.com/filipe1309/agl-go-driver/proto/v1/users"
)

func (s *UserService) Get(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	user, err := s.factory.RestoreOne(in.Id)
	if err != nil {
		// TODO: Check if the error is sql.ErrNoRows and return 404
		return &pb.UserResponse{Error: err.Error()}, err
	}

	return convertToUserResponse(user), nil
}
