package user_profile

import (
	"context"
	"log"

	pbCommon "github.com/rdev-org/proto-contracts/pb/common"
	pb "github.com/rdev-org/proto-contracts/pb/user_profile"
	"github.com/rdev-org/users/internal/models"
	"github.com/rdev-org/users/internal/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserProfileService struct {
	pb.UnimplementedUserProfileServiceServer
	r UserProfileRepository
}

func NewUserProfileService(r UserProfileRepository) *UserProfileService {
	return &UserProfileService{
		r: r,
	}
}

func (s *UserProfileService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	p, err := s.r.GetProfile(ctx, req.GetUuid())
	if err != nil {
		log.Printf("[UserProfile]service - failed to get profile %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetProfileResponse{
		Profile: &pb.Profile{
			FirstName:    p.FirstName,
			MiddleName:   p.MiddleName,
			LastName:     p.LastName,
			MobileNumber: p.MobileNumber,
			DateOfBirth:  p.DateOfBirth,
			Gender:       p.Gender,
			Avatar:       p.Avatar,
		},
	}, nil
}

func (s *UserProfileService) ListProfile(ctx context.Context, req *pb.ListProfileRequest) (*pb.ListProfileResponse, error) {
	filter := &shared.Filter{
		Page:     int(req.Pagination.GetPage()),
		PageSize: int(req.Pagination.GetPageSize()),
		Sort:     req.Filter.GetSort(),
		Search:   req.Filter.GetSearch(),
	}

	profiles, count, err := s.r.ListProfile(ctx, filter)
	if err != nil {
		log.Printf("[UserProfile]service - failed to list profile %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &pb.ListProfileResponse{
		Profile: []*pb.Profile{},
		Pagination: &pbCommon.PaginationResponse{
			TotalCount: int64(count),
		},
	}

	if len(profiles) > 0 {
		for _, p := range profiles {
			resp.Profile = append(resp.Profile, &pb.Profile{
				FirstName:    p.FirstName,
				MiddleName:   p.MiddleName,
				LastName:     p.LastName,
				MobileNumber: p.MobileNumber,
				DateOfBirth:  p.DateOfBirth,
				Gender:       p.Gender,
				Avatar:       p.Avatar,
				Meta: &pbCommon.Metadata{
					Uuid:      p.Meta.Uuid,
					Status:    p.Meta.Status,
					CreatedAt: timestamppb.New(p.Meta.CreatedAt),
					UpdatedAt: timestamppb.New(p.Meta.UpdatedAt),
					DeletedAt: timestamppb.New(p.Meta.DeletedAt),
				},
			})
		}
	}

	return resp, nil
}

func (s *UserProfileService) CreateProfile(ctx context.Context, req *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {
	reqProfile := req.GetProfile()
	payload := &models.UserProfile{
		UserID:       reqProfile.GetUserId(),
		FirstName:    reqProfile.GetFirstName(),
		MiddleName:   reqProfile.GetMiddleName(),
		LastName:     reqProfile.GetLastName(),
		MobileNumber: reqProfile.GetMobileNumber(),
		DateOfBirth:  reqProfile.GetDateOfBirth(),
		Gender:       reqProfile.GetGender(),
		Avatar:       reqProfile.GetAvatar(),
	}

	err := s.r.CreateProfile(ctx, payload)
	if err != nil {
		log.Printf("[UserProfile]service - failed to create profile %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateProfileResponse{
		Profile: &pb.Profile{
			Meta: &pbCommon.Metadata{
				Uuid: payload.Meta.Uuid,
			},
		},
	}, nil
}

func (s *UserProfileService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	reqProfile := req.GetProfile()
	payload := &models.UserProfile{
		Uuid:         req.GetUuid(),
		UserID:       reqProfile.GetUserId(),
		FirstName:    reqProfile.GetFirstName(),
		MiddleName:   reqProfile.GetMiddleName(),
		LastName:     reqProfile.GetLastName(),
		MobileNumber: reqProfile.GetMobileNumber(),
		DateOfBirth:  reqProfile.GetDateOfBirth(),
		Gender:       reqProfile.GetGender(),
		Avatar:       reqProfile.GetAvatar(),
	}

	err := s.r.UpdateProfile(ctx, payload)
	if err != nil {
		log.Printf("[UserProfile]service - failed to update profile %v", err)
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

func (s *UserProfileService) DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest) (*pb.DeleteProfileResponse, error) {
	err := s.r.DeleteProfile(ctx, req.GetUuid())
	if err != nil {
		log.Printf("[UserProfile]service - failed to delete profile %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteProfileResponse{
		Message: "Profile successfully deleted",
	}, nil
}
