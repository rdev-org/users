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

func (s *UserService) GetPreference(ctx context.Context, req *pb.GetPreferenceRequest) (*pb.GetPreferenceResponse, error) {
	p, err := s.preference.GetPreference(ctx, req.GetUuid())
	if err != nil {
		log.Printf("[UserPreference]service - failed to get preference %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	preferenceProto, err := p.ToProto()
	if err != nil {
		log.Printf("[UserPreference]service - failed to convert preferences to proto: %v", err)
		return nil, status.Error(codes.Internal, "failed to parse preferences")
	}

	return &pb.GetPreferenceResponse{
		Preference: preferenceProto,
	}, nil
}

func (s *UserService) ListPreference(ctx context.Context, req *pb.ListPreferenceRequest) (*pb.ListPreferenceResponse, error) {
	filter := &shared.Filter{
		Page:     int(req.Pagination.GetPage()),
		PageSize: int(req.Pagination.GetPageSize()),
		Sort:     req.Filter.GetSort(),
		Search:   req.Filter.GetSearch(),
	}

	preferences, count, err := s.preference.ListPreference(ctx, filter)
	if err != nil {
		log.Printf("[UserPreference]service - failed to list preference %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoPreferences, err := models.UserPreferencesToProto(preferences)
	if err != nil {
		log.Printf("[UserPreference]service - failed to convert preferences to proto: %v", err)
		return nil, status.Error(codes.Internal, "failed to parse preferences")
	}

	return &pb.ListPreferenceResponse{
		Preference: protoPreferences,
		Pagination: &pbCommon.PaginationResponse{
			TotalCount: int64(count),
		},
	}, nil
}

func (s *UserService) CreatePreference(ctx context.Context, req *pb.CreatePreferenceRequest) (*pb.CreatePreferenceResponse, error) {
	payload, err := models.ProtoToCreateUserPreference(req)
	if err != nil {
		log.Printf("[UserAddress]service - failed to convert proto to user preference %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	_, err = s.preference.CreatePreference(ctx, payload)
	if err != nil {
		log.Printf("[UserPreference]service - failed to create preference %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreatePreferenceResponse{
		Preference: &pb.Preference{
			Meta: &pbCommon.Metadata{
				Uuid: payload.Meta.Uuid,
			},
		},
	}, nil
}

func (s *UserService) UpdatePreference(ctx context.Context, req *pb.UpdatePreferenceRequest) (*pb.UpdatePreferenceResponse, error) {
	payload, err := models.ProtoToUpdateUserPreference(req)
	if err != nil {
		log.Printf("[UserAddress]service - failed to convert proto to user preference %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.preference.UpdatePreference(ctx, payload)
	if err != nil {
		log.Printf("[UserPreference]service - failed to update preference %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdatePreferenceResponse{
		Preference: &pb.Preference{
			Meta: &pbCommon.Metadata{
				Uuid: req.GetUuid(),
			},
		},
	}, nil
}

func (s *UserService) DeletePreference(ctx context.Context, req *pb.DeletePreferenceRequest) (*pb.DeletePreferenceResponse, error) {
	err := s.preference.DeletePreference(ctx, req.GetUuid())
	if err != nil {
		log.Printf("[UserPreference]service - failed to delete preference %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeletePreferenceResponse{
		Message: "Preference successfully deleted",
	}, nil
}
