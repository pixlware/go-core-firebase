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
	Update func(uid string, user *auth.UserToUpdate, tenantId string) (*auth.UserRecord, error)
	Delete func(uid string, tenantId string) error
}

type TenantService struct {
	Create            func(tenant *auth.TenantToCreate) (*auth.Tenant, error)
	Update            func(tenantId string, updateProps *auth.TenantToUpdate) error
	UpdateInheritance func(tenantId string, inheritance *TenantInheritance) error
	Delete            func(tenantId string) error
}

type HandlerService struct {
	ValidateAuth func(c *fiber.Ctx) error
}

type IdentityPlatformTenantUpdater struct {
	Inheritance *TenantInheritance `json:"inheritance,omitempty"`
}

type TenantInheritance struct {
	EmailSendingConfig bool `json:"emailSendingConfig,omitempty"`
}
