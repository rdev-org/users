package models

import (
	"context"
	"fmt"
	"time"
	"uuid"

	"github.com/rdev-org/helpers/structutil"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbCommon "github.com/rdev-org/proto-contracts/pb/common"
	pb "github.com/rdev-org/proto-contracts/pb/user"

	"github.com/rdev-org/users/internal/shared"
	"github.com/uptrace/bun"
)

type UserPreference struct {
	bun.BaseModel `bun:"table:user_preference,alias:upr"`
	shared.Meta

	UserID     string            `bun:"user_id" json:"user_id"`
	Language   string            `json:"language"`
	Currency   string            `json:"currency"`
	EmailOptIn bool              `json:"email_opt_in"`
	SmsOptIn   bool              `json:"sms_opt_in"`
	Settings   map[string]string `json:"settings"`

	Profile *UserProfile `bun:"rel:belongs-to,join:user_id=user_id" json:"profile"`
}

func (u *UserPreference) ToProto() (*pb.Preference, error) {
	if u == nil {
		return nil, nil
	}

	settingsStruct, err := structutil.ToProtoStruct(u.Settings)
	if err != nil {
		return nil, err
	}

	meta := &pbCommon.Metadata{
		Uuid:      u.Meta.Uuid,
		Status:    u.Meta.Status,
		CreatedAt: timestamppb.New(u.Meta.CreatedAt),
		UpdatedAt: timestamppb.New(u.Meta.UpdatedAt),
	}
	if !u.Meta.DeletedAt.IsZero() {
		meta.DeletedAt = timestamppb.New(u.Meta.DeletedAt)
	}

	return &pb.Preference{
		Language:   u.Language,
		Currency:   u.Currency,
		EmailOptIn: u.EmailOptIn,
		SmsOptIn:   u.SmsOptIn,
		Settings:   settingsStruct,
		Meta:       meta,
	}, nil
}

func ProtoToUserPreference(preference *pb.Preference) *UserPreference {
	if preference == nil {
		return nil
	}

	return &UserPreference{
		UserID:     preference.GetUserId(),
		Language:   preference.GetLanguage(),
		Currency:   preference.GetCurrency(),
		EmailOptIn: preference.GetEmailOptIn(),
		SmsOptIn:   preference.GetSmsOptIn(),
		Settings:   structutil.ToStringMap(preference.GetSettings()),
	}
}

func ProtoToCreateUserPreference(proto *pb.CreatePreferenceRequest) (*UserPreference, error) {
	if proto == nil || proto.GetPreference() == nil {
		return nil, fmt.Errorf("preference data is required")
	}

	return ProtoToUserPreference(proto.GetPreference()), nil
}

func ProtoToUpdateUserPreference(proto *pb.UpdatePreferenceRequest) (*UserPreference, error) {
	if proto == nil || proto.GetPreference() == nil {
		return nil, fmt.Errorf("preference data is required")
	}

	preference := ProtoToUserPreference(proto.GetPreference())
	if uuidVal := proto.GetUuid(); uuidVal != "" {
		preference.Uuid = uuidVal
	}

	return preference, nil
}

func UserPreferencesToProto(preferences []*UserPreference) ([]*pb.Preference, error) {
	if len(preferences) == 0 {
		return []*pb.Preference{}, nil
	}

	protos := make([]*pb.Preference, len(preferences))
	for i, p := range preferences {
		proto, err := p.ToProto()
		if err != nil {
			return nil, err
		}
		protos[i] = proto
	}

	return protos, nil
}

func (u *UserPreference) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if u.Status == "" {
			u.Status = "A"
		}
		if u.Uuid == "" {
			u.Uuid = uuid.New().String()
		}
	case *bun.UpdateQuery:
		if u.UpdatedAt.IsZero() {
			u.UpdatedAt = time.Now()
		}
	}
	return nil
}
