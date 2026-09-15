package user_addresses

import (
	"context"

	pb "github.com/rdev-org/proto-contracts/pb/user_address"
	"github.com/rdev-org/users/internal/models"
	"github.com/rdev-org/users/internal/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserAddressService struct {
	pb.UnimplementedUserAddressServiceServer
	repo IUserAddressRepository
}

func NewUserAddressServiceServer(repo IUserAddressRepository) *UserAddressService {
	return &UserAddressService{
		repo: repo,
	}
}

func (s *UserAddressService) GetAddress(ctx context.Context, req *pb.GetAddressRequest) (*pb.GetAddressResponse, error) {
	addr, err := s.repo.GetAddress(ctx, req.GetUuid())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &pb.GetAddressResponse{
		Address: &pb.Address{
			AddressLine1:  addr.AddressLine1,
			AddressLine2:  addr.AddressLine2,
			City:          addr.City,
			StateProvince: addr.StateProvince,
			PostalCode:    addr.PostalCode,
			Country:       addr.Country,
			IsDefault:     addr.IsDefault,
			AddressType:   addr.AddressType,
			Status:        addr.Status,
			Uuid:          addr.Uuid,
			CreatedAt:     timestamppb.New(addr.CreatedAt),
			UpdatedAt:     timestamppb.New(addr.UpdatedAt),
			DeletedAt:     timestamppb.New(addr.DeletedAt),
		},
	}

	return resp, nil
}

func (s *UserAddressService) GetAllAddress(ctx context.Context, req *pb.GetAllAddressRequest) (*pb.GetAllAddressResponse, error) {

	filter := &shared.Filter{
		UserId:   req.GetUserId(),
		Page:     int(req.Pagination.GetPage()),
		PageSize: int(req.Pagination.GetPageSize()),
		Sort:     req.Filter.GetSort(),
		Status:   req.Filter.GetStatus(),
		Search:   req.Filter.GetSearch(),
	}

	userAddrs, count, err := s.repo.GetAllAddress(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &pb.GetAllAddressResponse{
		Address: []*pb.Address{},
		Pagination: &pb.PaginationResponse{
			TotalCount: int64(count),
		},
	}

	if len(userAddrs) > 0 {
		for _, addr := range userAddrs {
			resp.Address = append(resp.Address, &pb.Address{
				AddressLine1:  addr.AddressLine1,
				AddressLine2:  addr.AddressLine2,
				City:          addr.City,
				StateProvince: addr.StateProvince,
				PostalCode:    addr.PostalCode,
				Country:       addr.Country,
				IsDefault:     addr.IsDefault,
				AddressType:   addr.AddressType,
				Status:        addr.Status,
				Uuid:          addr.Uuid,
				CreatedAt:     timestamppb.New(addr.CreatedAt),
				UpdatedAt:     timestamppb.New(addr.UpdatedAt),
				DeletedAt:     timestamppb.New(addr.DeletedAt),
			})
		}
	}

	return resp, nil
}

func (s *UserAddressService) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.CreateAddressResponse, error) {

	addr := req.GetAddress()
	payload := &models.UserAddress{
		UserID:        addr.GetUserId(),
		AddressLine1:  addr.GetAddressLine1(),
		AddressLine2:  addr.GetAddressLine2(),
		City:          addr.GetCity(),
		StateProvince: addr.GetStateProvince(),
		PostalCode:    addr.GetPostalCode(),
		Country:       addr.GetCountry(),
		IsDefault:     addr.GetIsDefault(),
		AddressType:   addr.GetAddressType(),
	}

	_, err := s.repo.CreateAddress(ctx, payload)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &pb.CreateAddressResponse{
		Address: req.GetAddress(),
	}

	resp.Address.Uuid = payload.Uuid

	return resp, nil
}

func (s *UserAddressService) UpdateAddress(ctx context.Context, req *pb.UpdateAddressRequest) (*pb.UpdateAddressResponse, error) {
	return nil, nil
}

func (s *UserAddressService) DeleteAddress(ctx context.Context, req *pb.DeleteAddressRequest) (*pb.DeleteAddressResponse, error) {
	return nil, nil
}

func (s *UserAddressService) SetDefaultAddress(ctx context.Context, req *pb.SetDefaultAddressRequest) (*pb.SetDefaultAddressResponse, error) {
	return nil, nil
}
