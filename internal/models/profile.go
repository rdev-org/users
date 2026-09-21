package models

import (
	"context"
	"fmt"
	"time"
	"uuid"

	pbCommon "github.com/rdev-org/proto-contracts/pb/common"
	pb "github.com/rdev-org/proto-contracts/pb/user"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/rdev-org/helpers/structutil"
	"github.com/rdev-org/users/internal/dto"
	"github.com/rdev-org/users/internal/jwt"
	"github.com/rdev-org/users/internal/shared"
	"github.com/uptrace/bun"
)

type Profile struct {
	bun.BaseModel `bun:"table:user_profile,alias:up"`
	shared.Meta

	UserID       string `json:"user_id"`
	FirstName    string `json:"first_name"`
	MiddleName   string `json:"middle_name"`
	LastName     string `json:"last_name"`
	MobileNumber string `json:"mobile_number"`
	DateOfBirth  string `json:"date_of_birth"`
	Gender       string `json:"gender"`
	Avatar       string `json:"avatar"`

	Addresses  []*Address  `bun:"rel:has-many,join:user_id=user_id" json:"addresses,omitempty"`
	Preference *Preference `bun:"rel:has-one,join:user_id=user_id" json:"preference,omitempty"`
	User       *Credential `bun:"rel:belongs-to,join:user_id=uuid" json:"-"`
}

func (u *Profile) ToProto() (*pb.Profile, error) {
	if u == nil {
		return nil, nil
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

	resp := &pb.Profile{
		FirstName:    u.FirstName,
		MiddleName:   u.MiddleName,
		LastName:     u.LastName,
		MobileNumber: u.MobileNumber,
		DateOfBirth:  u.DateOfBirth,
		Gender:       u.Gender,
		Avatar:       u.Avatar,
		Addresses:    AddressesToProto(u.Addresses),
		Meta:         meta,
	}

	if u.Preference != nil {
		settingsStruct, err := structutil.ToProtoStruct(u.Preference.Settings)
		if err != nil {
			return nil, err
		}

		resp.Preference = &pb.Preference{
			Language:   u.Preference.Language,
			Currency:   u.Preference.Currency,
			EmailOptIn: u.Preference.EmailOptIn,
			SmsOptIn:   u.Preference.SmsOptIn,
			Settings:   settingsStruct,
		}
	}

	return resp, nil
}

func ProtoToProfile(profile *pb.Profile) *Profile {
	if profile == nil {
		return nil
	}

	return &Profile{
		UserID:       profile.GetUserId(),
		FirstName:    profile.GetFirstName(),
		MiddleName:   profile.GetMiddleName(),
		LastName:     profile.GetLastName(),
		MobileNumber: profile.GetMobileNumber(),
		DateOfBirth:  profile.GetDateOfBirth(),
		Gender:       profile.GetGender(),
		Avatar:       profile.GetAvatar(),
	}
}

func ProtoToCreateProfile(proto *pb.CreateProfileRequest) (*dto.ProfileRequestDTO, error) {
	if proto == nil || proto.GetPayload() == nil {
		return nil, fmt.Errorf("preference data is required")
	}

	reqProfile := proto.GetPayload()
	payload := &dto.ProfileRequestDTO{
		FirstName:    reqProfile.GetFirstName(),
		MiddleName:   reqProfile.GetMiddleName(),
		LastName:     reqProfile.GetLastName(),
		MobileNumber: reqProfile.GetMobileNumber(),
		DateOfBirth:  reqProfile.GetDateOfBirth(),
		Gender:       reqProfile.GetGender(),
		Avatar:       reqProfile.GetAvatar(),
		Address:      &dto.AddressRequestDTO{},

		Token:    reqProfile.GetToken(),
		Username: reqProfile.GetUsername(),
		Password: jwt.HashPassword(reqProfile.GetPassword()),
	}

	if reqProfile.GetAddress() != nil {
		reqAddress := reqProfile.GetAddress()

		payload.Address = &dto.AddressRequestDTO{
			AddressLine1:  reqAddress.AddressLine1,
			AddressLine2:  reqAddress.AddressLine2,
			City:          reqAddress.City,
			StateProvince: reqAddress.StateProvince,
			PostalCode:    reqAddress.PostalCode,
			Country:       reqAddress.Country,
			IsDefault:     reqAddress.IsDefault,
			AddressType:   reqAddress.AddressType,
		}
	}

	return payload, nil
}

func ProtoToUpdateProfile(proto *pb.UpdateProfileRequest) (*Profile, error) {
	if proto == nil || proto.GetProfile() == nil {
		return nil, fmt.Errorf("preference data is required")
	}

	preference := ProtoToProfile(proto.GetProfile())
	if uuidVal := proto.GetUuid(); uuidVal != "" {
		preference.Uuid = uuidVal
	}

	return preference, nil
}

func ProfileToProto(profiles []*Profile) ([]*pb.Profile, error) {
	if len(profiles) == 0 {
		return []*pb.Profile{}, nil
	}

	protos := make([]*pb.Profile, len(profiles))
	for i, p := range profiles {
		proto, err := p.ToProto()
		if err != nil {
			return nil, err
		}
		protos[i] = proto
	}

	return protos, nil
}

func (u *Profile) BeforeAppendModel(ctx context.Context, query bun.Query) error {
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
