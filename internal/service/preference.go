package service

import (
	"context"
	"log"

	"github.com/rdev-org/helpers/structutil"
	pbCommon "github.com/rdev-org/proto-contracts/pb/common"
	pb "github.com/rdev-org/proto-contracts/pb/user"
	"github.com/rdev-org/users/internal/models"
	"github.com/rdev-org/users/internal/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *UserService) GetPreference(ctx context.Context, req *pb.GetPreferenceRequest) (*pb.GetPreferenceResponse, error) {
	p, err := s.preference.GetPreference(ctx, req.GetUuid())
	if err != nil {
		log.Printf("[UserPreference]service - failed to get preference %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	settingsStruct, err := structutil.ToProtoStruct(p.Settings)
	if err != nil {
		log.Printf("[UserPreference]service - failed to convert settings to structpb: %v", err)
		return nil, status.Error(codes.Internal, "failed to parse settings")
	}

	return &pb.GetPreferenceResponse{
		Preference: &pb.Preference{
			Language:   p.Language,
			Currency:   p.Currency,
			EmailOptIn: p.EmailOptIn,
			SmsOptIn:   p.SmsOptIn,
			Settings:   settingsStruct,
		},
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

	resp := &pb.ListPreferenceResponse{
		Preference: []*pb.Preference{},
		Pagination: &pbCommon.PaginationResponse{
			TotalCount: int64(count),
		},
	}

	if len(preferences) > 0 {
		for _, p := range preferences {
			settingsStruct, err := structutil.ToProtoStruct(p.Settings)
			if err != nil {
				log.Printf("[UserPreference]service - failed to convert settings to structpb: %v", err)
				return nil, status.Error(codes.Internal, "failed to parse settings")
			}

			resp.Preference = append(resp.Preference, &pb.Preference{
				Language:   p.Language,
				Currency:   p.Currency,
				EmailOptIn: p.EmailOptIn,
				SmsOptIn:   p.SmsOptIn,
				Settings:   settingsStruct,
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

func (s *UserService) CreatePreference(ctx context.Context, req *pb.CreatePreferenceRequest) (*pb.CreatePreferenceResponse, error) {
	reqPreference := req.GetPreference()

	payload := &models.UserPreference{
		UserID:     reqPreference.GetUserId(),
		Language:   reqPreference.GetLanguage(),
		Currency:   reqPreference.GetCurrency(),
		EmailOptIn: reqPreference.GetEmailOptIn(),
		SmsOptIn:   reqPreference.GetSmsOptIn(),
		Settings:   structutil.ToStringMap(reqPreference.GetSettings()),
	}

	err := s.preference.CreatePreference(ctx, payload)
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
	reqPreference := req.GetPreference()

	payload := &models.UserPreference{
		Uuid:       req.GetUuid(),
		UserID:     reqPreference.GetUserId(),
		Language:   reqPreference.GetLanguage(),
		Currency:   reqPreference.GetCurrency(),
		EmailOptIn: reqPreference.GetEmailOptIn(),
		SmsOptIn:   reqPreference.GetSmsOptIn(),
		Settings:   structutil.ToStringMap(reqPreference.GetSettings()),
	}

	err := s.preference.UpdatePreference(ctx, payload)
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
