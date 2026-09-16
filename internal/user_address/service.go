package user_address

import (
	"context"

	pb "github.com/rdev-org/proto-contracts/pb/user_address"
)

type UserAddressService struct {
	pb.UnimplementedUserAddressServiceServer
	r UserAddressRepository
}

func NewUserAddressService(r UserAddressRepository) *UserAddressService {
	return &UserAddressService{
		r: r,
	}
}

func (s *UserAddressService) GetAddress(ctx context.Context, req *pb.GetAddressRequest) (*pb.GetAddressResponse, error) {
	return nil, nil
}

func (s *UserAddressService) ListAddress(ctx context.Context, req *pb.ListAddressRequest) (*pb.ListAddressResponse, error) {
	return nil, nil
}

func (s *UserAddressService) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.CreateAddressResponse, error) {
	return nil, nil
}

func (s *UserAddressService) UpdateAddress(ctx context.Context, req *pb.UpdateAddressRequest) (*pb.UpdateAddressResponse, error) {
	return nil, nil
}

func (s *UserAddressService) DeleteAddress(ctx context.Context, req *pb.DeleteAddressRequest) (*pb.DeleteAddressResponse, error) {
	return nil, nil
}

func (s *UserAddressService) DefaultAddress(ctx context.Context, req *pb.DefaultAddressRequest) (*pb.DefaultAddressResponse, error) {
	return nil, nil
}
