package users

import (
	"context"

	pb "github.com/filipe1309/agl-go-driver/proto/v1/users"
)

func (s *UserService) Update(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	user, err := s.factory.RestoreOne(in.Id)
	if err != nil {
		return &pb.UserResponse{Error: err.Error()}, err
	}

	err = user.ChangeName(in.Name)
	if err != nil {
		return &pb.UserResponse{Error: err.Error()}, err
	}

	_, err = s.repo.UpdateDB(in.Id, user)
	if err != nil {
		return &pb.UserResponse{Error: err.Error()}, err
	}

	return convertToUserResponse(user), nil
}
