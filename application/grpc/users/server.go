package users

import (
	pb "github.com/filipe1309/agl-go-driver/proto/v1/users"
	"github.com/filipe1309/agl-go-driver/repositories"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	repo repositories.UserWriteRepository
}
