package service

import (
	pb "github.com/rdev-org/proto-contracts/pb/user"
	"github.com/rdev-org/users/internal/repository"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	profile    repository.UserProfileRepository
	address    repository.UserAddressRepository
	preference repository.UserPreferenceRepository
}

func NewUserService(
	profile repository.UserProfileRepository,
	address repository.UserAddressRepository,
	preference repository.UserPreferenceRepository,
) *UserService {
	return &UserService{
		profile:    profile,
		address:    address,
		preference: preference,
	}
}
