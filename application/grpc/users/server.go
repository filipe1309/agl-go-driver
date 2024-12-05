package users

import (
	"github.com/filipe1309/agl-go-driver/factories"
	domain "github.com/filipe1309/agl-go-driver/internal/users"
	pb "github.com/filipe1309/agl-go-driver/proto/v1/users"
	"github.com/filipe1309/agl-go-driver/repositories"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	repo    repositories.UserWriteRepository
	factory factories.UserFactory
}

func convertToUserResponse(user *domain.User) *pb.UserResponse {
	return &pb.UserResponse{
		User: &pb.User{
			Id:        user.ID,
			Name:      user.Name,
			Login:     user.Login,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}
}
