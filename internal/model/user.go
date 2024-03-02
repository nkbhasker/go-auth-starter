package model

import (
	"github.com/nkbhasker/go-auth-starter/internal/enum"
	"github.com/nkbhasker/go-auth-starter/internal/uid"
)

type User struct {
	ID               uid.Identifier             `json:"id" gorm:"primaryKey;type:bigint;serializer:id;" kind:"user"`
	FirstName        string                     `json:"firstName" gorm:"not null"`
	LastName         *string                    `json:"lastName"`
	Email            *string                    `json:"email" gorm:"index:idx_email,unique,where:email IS NOT NULL"`
	Phone            *string                    `json:"phone" gorm:"index:idx_phone,unique,where:phone IS NOT NULL"`
	IsEmailVerified  bool                       `json:"isEmailVerified"`
	IsPhoneVerified  bool                       `json:"isPhoneVerified"`
	IsBot            bool                       `json:"-"`
	Gender           *enum.GenderEnum           `json:"gender" gorm:"type:gender"`
	IdentityProvider *enum.IdentityProviderEnum `json:"identityProvider" gorm:"type:identity_provider"`
}
