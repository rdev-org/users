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

type Preference struct {
	bun.BaseModel `bun:"table:user_preference,alias:upr"`
	shared.Meta

	UserID     string            `bun:"user_id" json:"user_id"`
	Language   string            `json:"language"`
	Currency   string            `json:"currency"`
	EmailOptIn bool              `json:"email_opt_in"`
	SmsOptIn   bool              `json:"sms_opt_in"`
	Settings   map[string]string `json:"settings"`

	Profile *Profile `bun:"rel:belongs-to,join:user_id=user_id" json:"profile"`
}

func (u *Preference) ToProto() (*pb.Preference, error) {
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

func ProtoToPreference(preference *pb.Preference) *Preference {
	if preference == nil {
		return nil
	}

	return &Preference{
		UserID:     preference.GetUserId(),
		Language:   preference.GetLanguage(),
		Currency:   preference.GetCurrency(),
		EmailOptIn: preference.GetEmailOptIn(),
		SmsOptIn:   preference.GetSmsOptIn(),
		Settings:   structutil.ToStringMap(preference.GetSettings()),
	}
}

func ProtoToCreatePreference(proto *pb.CreatePreferenceRequest) (*Preference, error) {
	if proto == nil || proto.GetPreference() == nil {
		return nil, fmt.Errorf("preference data is required")
	}

	return ProtoToPreference(proto.GetPreference()), nil
}

func ProtoToUpdatePreference(proto *pb.UpdatePreferenceRequest) (*Preference, error) {
	if proto == nil || proto.GetPreference() == nil {
		return nil, fmt.Errorf("preference data is required")
	}

	preference := ProtoToPreference(proto.GetPreference())
	if uuidVal := proto.GetUuid(); uuidVal != "" {
		preference.Uuid = uuidVal
	}

	return preference, nil
}

func PreferencesToProto(preferences []*Preference) ([]*pb.Preference, error) {
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

func (u *Preference) BeforeAppendModel(ctx context.Context, query bun.Query) error {
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
