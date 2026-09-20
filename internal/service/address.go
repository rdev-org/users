package service

import (
	"context"
	"log"

	pbCommon "github.com/rdev-org/proto-contracts/pb/common"
	pb "github.com/rdev-org/proto-contracts/pb/user"
	"github.com/rdev-org/users/internal/models"
	"github.com/rdev-org/users/internal/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *UserService) GetAddress(ctx context.Context, req *pb.GetAddressRequest) (*pb.GetAddressResponse, error) {
	addr, err := s.address.GetAddress(ctx, req.GetUuid())
	if err != nil {
		log.Printf("[UserAddress]service - failed to get address %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetAddressResponse{
		Address: &pb.Address{
			AddressLine1:  addr.AddressLine1,
			AddressLine2:  addr.AddressLine2,
			City:          addr.City,
			StateProvince: addr.StateProvince,
			PostalCode:    addr.PostalCode,
			Country:       addr.Country,
			IsDefault:     addr.IsDefault,
			AddressType:   addr.AddressType,
		},
	}, nil
}

func (s *UserService) ListAddress(ctx context.Context, req *pb.ListAddressRequest) (*pb.ListAddressResponse, error) {
	filter := &shared.Filter{
		Page:     int(req.Pagination.GetPage()),
		PageSize: int(req.Pagination.GetPageSize()),
		Sort:     req.Filter.GetSort(),
		Search:   req.Filter.GetSearch(),
	}

	addrs, count, err := s.address.ListAddress(ctx, filter)
	if err != nil {
		log.Printf("[UserAddress]service - failed to list address %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &pb.ListAddressResponse{
		Address: []*pb.Address{},
		Pagination: &pbCommon.PaginationResponse{
			TotalCount: int64(count),
		},
	}

	if len(addrs) > 0 {
		for _, addr := range addrs {
			resp.Address = append(resp.Address, &pb.Address{
				AddressLine1:  addr.AddressLine1,
				AddressLine2:  addr.AddressLine2,
				City:          addr.City,
				StateProvince: addr.StateProvince,
				PostalCode:    addr.PostalCode,
				Country:       addr.Country,
				IsDefault:     addr.IsDefault,
				AddressType:   addr.AddressType,
				Meta: &pbCommon.Metadata{
					Uuid:      addr.Meta.Uuid,
					Status:    addr.Meta.Status,
					CreatedAt: timestamppb.New(addr.Meta.CreatedAt),
					UpdatedAt: timestamppb.New(addr.Meta.UpdatedAt),
					DeletedAt: timestamppb.New(addr.Meta.DeletedAt),
				},
			})
		}
	}

	return resp, nil
}

func (s *UserService) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.CreateAddressResponse, error) {
	reqAddr := req.GetAddress()
	payload := &models.UserAddress{
		UserID:        reqAddr.GetUserId(),
		AddressLine1:  reqAddr.GetAddressLine1(),
		AddressLine2:  reqAddr.GetAddressLine2(),
		City:          reqAddr.GetCity(),
		StateProvince: reqAddr.GetStateProvince(),
		PostalCode:    reqAddr.GetPostalCode(),
		Country:       reqAddr.GetCountry(),
		IsDefault:     reqAddr.GetIsDefault(),
		AddressType:   reqAddr.GetAddressType(),
	}

	err := s.address.CreateAddress(ctx, payload)
	if err != nil {
		log.Printf("[UserAddress]service - failed to create address %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateAddressResponse{
		Address: &pb.Address{
			Meta: &pbCommon.Metadata{
				Uuid: payload.Meta.Uuid,
			},
		},
	}, nil
}

func (s *UserService) UpdateAddress(ctx context.Context, req *pb.UpdateAddressRequest) (*pb.UpdateAddressResponse, error) {
	reqAddr := req.GetAddress()
	payload := &models.UserAddress{
		Uuid:          req.GetUuid(),
		UserID:        reqAddr.GetUserId(),
		AddressLine1:  reqAddr.GetAddressLine1(),
		AddressLine2:  reqAddr.GetAddressLine2(),
		City:          reqAddr.GetCity(),
		StateProvince: reqAddr.GetStateProvince(),
		PostalCode:    reqAddr.GetPostalCode(),
		Country:       reqAddr.GetCountry(),
		IsDefault:     reqAddr.GetIsDefault(),
		AddressType:   reqAddr.GetAddressType(),
	}

	err := s.address.UpdateAddress(ctx, payload)
	if err != nil {
		log.Printf("[UserAddress]service - failed to update address %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateAddressResponse{
		Address: &pb.Address{
			Meta: &pbCommon.Metadata{
				Uuid: req.GetUuid(),
			},
		},
	}, nil
}

func (s *UserService) DeleteAddress(ctx context.Context, req *pb.DeleteAddressRequest) (*pb.DeleteAddressResponse, error) {
	err := s.address.DeleteAddress(ctx, req.GetUuid())
	if err != nil {
		log.Printf("[UserAddress]service - failed to delete address %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteAddressResponse{
		Message: "Address successfully deleted",
	}, nil
}

func (s *UserService) DefaultAddress(ctx context.Context, req *pb.DefaultAddressRequest) (*pb.DefaultAddressResponse, error) {
	err := s.address.DefaultAddress(ctx, req.GetUserId(), req.GetUuid())
	if err != nil {
		log.Printf("[UserAddress]service - failed to default address %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DefaultAddressResponse{
		Address: &pb.Address{
			Meta: &pbCommon.Metadata{
				Uuid: req.GetUuid(),
			},
		},
	}, nil
}
