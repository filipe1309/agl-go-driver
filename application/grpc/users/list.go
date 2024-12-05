package users

import (
	"context"

	pb "github.com/filipe1309/agl-go-driver/proto/v1/users"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserService) List(ctx context.Context, in *emptypb.Empty) (*pb.UserListResponse, error) {
	users, err := s.factory.RestoreAll()
	if err != nil {
		// TODO: Check if the error is sql.ErrNoRows and return 404
		return &pb.UserListResponse{Error: err.Error()}, err
	}

	data := make([]*pb.User, 0, len(users))
	for i, user := range users {
		data[i] = convertToUserPb(&user)
	}

	return &pb.UserListResponse{Users: data}, nil
}
