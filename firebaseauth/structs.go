package firebaseauth

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
)

type AuthUser struct {
	UserID        string `json:"userId"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	AuthTime      int64  `json:"authTime"`
	TenantID      string `json:"tenantId"`
}

type UserService struct {
	Create func(user *auth.UserToCreate, tenantId string) (*auth.UserRecord, error)
	Delete func(uid string, tenantId string) error
}

type TenantService struct {
	Update                   func(tenantId string, updateProps *auth.TenantToUpdate) error
	UpdateEmailSendingConfig func(tenantId string, emailSendingConfig *TenantEmailSendingConfig) error
}

type HandlerService struct {
	ValidateAuth func(c *fiber.Ctx) error
}

type IdentityPlatformTenantUpdater struct {
	Inheritance *TenantInheritance `json:"inheritance,omitempty"`
}

type TenantInheritance struct {
	EmailSendingConfig *TenantEmailSendingConfig `json:"emailSendingConfig,omitempty"`
}

type TenantEmailSendingConfig struct {
	Enabled bool `json:"enabled,omitempty"`
}
