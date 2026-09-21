package service

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbCommon "github.com/rdev-org/proto-contracts/pb/common"
	pbEmail "github.com/rdev-org/proto-contracts/pb/notification"
	pb "github.com/rdev-org/proto-contracts/pb/user"

	"github.com/rdev-org/users/internal/constant"
	"github.com/rdev-org/users/internal/models"
	"github.com/rdev-org/users/internal/shared"
)

func (s *UserService) MyProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	p, err := s.profile.MyProfile(ctx, req.GetUuid())
	if err != nil {
		log.Printf("[Profile]service - failed to get profile %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	profileProto, err := p.ToProto()
	if err != nil {
		log.Printf("[Profile]service - failed to convert to proto: %v", err)
		return nil, status.Error(codes.Internal, "failed to parse")
	}

	return &pb.GetProfileResponse{
		Profile: profileProto,
	}, nil
}

func (s *UserService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	p, err := s.profile.GetProfile(ctx, req.GetUuid())
	if err != nil {
		log.Printf("[Profile]service - failed to get profile %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	profileProto, err := p.ToProto()
	if err != nil {
		log.Printf("[Profile]service - failed to convert to proto: %v", err)
		return nil, status.Error(codes.Internal, "failed to parse")
	}

	return &pb.GetProfileResponse{
		Profile: profileProto,
	}, nil
}

func (s *UserService) ListProfile(ctx context.Context, req *pb.ListProfileRequest) (*pb.ListProfileResponse, error) {
	filter := &shared.Filter{
		Page:     int(req.Pagination.GetPage()),
		PageSize: int(req.Pagination.GetPageSize()),
		Sort:     req.Filter.GetSort(),
		Search:   req.Filter.GetSearch(),
	}

	profiles, count, err := s.profile.ListProfile(ctx, filter)
	if err != nil {
		log.Printf("[Profile]service - failed to list profile %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoProfile, err := models.ProfileToProto(profiles)
	if err != nil {
		log.Printf("[Profile]service - failed to convert profiles to proto: %v", err)
		return nil, status.Error(codes.Internal, "failed to parse profiles")
	}

	return &pb.ListProfileResponse{
		Profile: protoProfile,
		Pagination: &pbCommon.PaginationResponse{
			TotalCount: int64(count),
		},
	}, nil
}

func (s *UserService) CreateProfile(ctx context.Context, req *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {
	payload, err := models.ProtoToCreateProfile(req)
	if err != nil {
		log.Printf("[Address]service - failed to convert proto to user profile %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	uuid, err := s.profile.CreateProfile(ctx, payload)
	if err != nil {
		log.Printf("[Profile]service - failed to create profile %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	name := req.GetPayload().GetFirstName()
	if req.GetPayload().GetMiddleName() != "" {
		name += " " + req.GetPayload().GetMiddleName()
	}
	name += " " + req.GetPayload().GetLastName()

	_, err = s.email.SendEmail(ctx, &pbEmail.SendEmailRequest{
		Module:         "REGISTER_EMAIL",
		RequestContext: req.RequestContext,
		Email: &pbEmail.Email{
			Recipient: req.GetPayload().GetEmail(),
			Name:      name,
			Subject:   fmt.Sprintf("New sign in from %s", req.GetPayload().GetEmail()),
			Body:      constant.WelcomeMessage,
		},
	})
	if err != nil {
		log.Printf("[Profile]service - failed to send email login %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateProfileResponse{
		Profile: &pb.Profile{
			Meta: &pbCommon.Metadata{
				Uuid: uuid,
			},
		},
	}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	payload, err := models.ProtoToUpdateProfile(req)
	if err != nil {
		log.Printf("[Profile]service - failed to convert proto to user preference %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.profile.UpdateProfile(ctx, payload)
	if err != nil {
		log.Printf("[Profile]service - failed to update profile %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateProfileResponse{
		Profile: &pb.Profile{
			Meta: &pbCommon.Metadata{
				Uuid: req.GetUuid(),
			},
		},
	}, nil
}

func (s *UserService) DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest) (*pb.DeleteProfileResponse, error) {
	err := s.profile.DeleteProfile(ctx, req.GetUuid())
	if err != nil {
		log.Printf("[Profile]service - failed to delete profile %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteProfileResponse{
		Message: "Profile successfully deleted",
	}, nil
}
