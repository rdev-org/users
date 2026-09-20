package service

import (
	pbEmail "github.com/rdev-org/proto-contracts/pb/notification"
	pb "github.com/rdev-org/proto-contracts/pb/user"
	"github.com/rdev-org/users/internal/repository"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	email      pbEmail.NotificationServiceClient
	profile    repository.UserProfileRepository
	address    repository.UserAddressRepository
	preference repository.UserPreferenceRepository
}

func NewUserService(
	email pbEmail.NotificationServiceClient,
	profile repository.UserProfileRepository,
	address repository.UserAddressRepository,
	preference repository.UserPreferenceRepository,
) *UserService {
	return &UserService{
		email:      email,
		profile:    profile,
		address:    address,
		preference: preference,
	}
}
