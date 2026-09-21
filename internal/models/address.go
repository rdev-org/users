package models

import (
	"context"
	"fmt"
	"time"
	"uuid"

	pbCommon "github.com/rdev-org/proto-contracts/pb/common"
	pb "github.com/rdev-org/proto-contracts/pb/user"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/rdev-org/users/internal/shared"
	"github.com/uptrace/bun"
)

type Address struct {
	bun.BaseModel `bun:"table:user_address,alias:ua"`
	shared.Meta

	UserID        string `bun:"user_id" json:"user_id"`
	AddressLine1  string `json:"address_line1"`
	AddressLine2  string `json:"address_line2,omitempty"`
	City          string `json:"city"`
	StateProvince string `json:"state_province"`
	PostalCode    string `json:"postal_code"`
	Country       string `json:"country"`
	IsDefault     bool   `json:"is_default"`
	AddressType   string `json:"address_type"`

	Profile *Profile `bun:"rel:belongs-to,join:user_id=user_id" json:"profile"`
}

func (u *Address) ToProto() *pb.Address {
	if u == nil {
		return nil
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

	return &pb.Address{
		AddressLine1:  u.AddressLine1,
		AddressLine2:  u.AddressLine2,
		City:          u.City,
		StateProvince: u.StateProvince,
		PostalCode:    u.PostalCode,
		Country:       u.Country,
		IsDefault:     u.IsDefault,
		AddressType:   u.AddressType,
		Meta:          meta,
	}
}

func AddressesToProto(addresses []*Address) []*pb.Address {
	if len(addresses) == 0 {
		return nil
	}

	protos := make([]*pb.Address, len(addresses))
	for i, addr := range addresses {
		protos[i] = addr.ToProto()
	}

	return protos
}

func ProtoToAddress(addr *pb.Address) *Address {
	if addr == nil {
		return nil
	}

	return &Address{
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
}

func ProtoToCreateAddress(proto *pb.CreateAddressRequest) (*Address, error) {
	if proto == nil || proto.GetAddress() == nil {
		return nil, fmt.Errorf("address data is required")
	}

	return ProtoToAddress(proto.GetAddress()), nil
}

func ProtoToUpdateAddress(proto *pb.UpdateAddressRequest) (*Address, error) {
	if proto == nil || proto.GetAddress() == nil {
		return nil, fmt.Errorf("address data is required")
	}

	address := ProtoToAddress(proto.GetAddress())
	if uuidVal := proto.GetUuid(); uuidVal != "" {
		address.Uuid = uuidVal
	}

	return address, nil
}

func (u *Address) BeforeAppendModel(ctx context.Context, query bun.Query) error {
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
