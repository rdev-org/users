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
)

func (s *UserService) GetAddress(ctx context.Context, req *pb.GetAddressRequest) (*pb.GetAddressResponse, error) {
	address, err := s.address.GetAddress(ctx, req.GetUuid())
	if err != nil {
		log.Printf("[Address]service - failed to get address %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetAddressResponse{
		Address: address.ToProto(),
	}, nil
}

func (s *UserService) ListAddress(ctx context.Context, req *pb.ListAddressRequest) (*pb.ListAddressResponse, error) {
	filter := &shared.Filter{
		Page:     int(req.Pagination.GetPage()),
		PageSize: int(req.Pagination.GetPageSize()),
		Sort:     req.Filter.GetSort(),
		Search:   req.Filter.GetSearch(),
	}

	addresses, count, err := s.address.ListAddress(ctx, filter)
	if err != nil {
		log.Printf("[Address]service - failed to list address %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoAddress := models.AddressesToProto(addresses)

	return &pb.ListAddressResponse{
		Address: protoAddress,
		Pagination: &pbCommon.PaginationResponse{
			TotalCount: int64(count),
		},
	}, nil
}

func (s *UserService) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.CreateAddressResponse, error) {
	payload, err := models.ProtoToCreateAddress(req)
	if err != nil {
		log.Printf("[Address]service - failed to convert proto to user address %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	_, err = s.address.CreateAddress(ctx, payload)
	if err != nil {
		log.Printf("[Address]service - failed to create address %v", err)
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
	payload, err := models.ProtoToUpdateAddress(req)
	if err != nil {
		log.Printf("[Address]service - failed to convert proto to user address %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.address.UpdateAddress(ctx, payload)
	if err != nil {
		log.Printf("[Address]service - failed to update address %v", err)
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
		log.Printf("[Address]service - failed to delete address %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteAddressResponse{
		Message: "Address successfully deleted",
	}, nil
}

func (s *UserService) DefaultAddress(ctx context.Context, req *pb.DefaultAddressRequest) (*pb.DefaultAddressResponse, error) {
	err := s.address.DefaultAddress(ctx, req.GetUserId(), req.GetUuid())
	if err != nil {
		log.Printf("[Address]service - failed to default address %v", err)
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
