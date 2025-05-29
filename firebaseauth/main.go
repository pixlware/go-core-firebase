package firebaseauth

import (
	"context"
	"log"
	"net/http"
	"slices"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/pixlware/go-core-fiber/fiberutils"
)

var (
	App        *firebase.App
	Auth       *auth.Client
	TenantAuth *auth.TenantClient
)

func init() {
	log.Printf("[Firebase] Initializing Firebase...")

	firebaseConfig := &firebase.Config{
		ProjectID: Config.ProjectID,
	}
	firebaseApp, err := firebase.NewApp(context.Background(), firebaseConfig)
	if err != nil {
		log.Fatalf("[Firebase] Error initializing Firebase App: %v\n", err)
	}
	App = firebaseApp
	log.Printf("[Firebase] Firebase App Initialized")

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("[Firebase] Error initializing Auth Client: %v\n", err)
	}
	Auth = authClient
	log.Printf("[Firebase] Auth Client Initialized")

	if Config.EnforceTenant && Config.TenantID != "" {
		tenantAuthClient, err := Auth.TenantManager.AuthForTenant(Config.TenantID)
		if err != nil {
			log.Fatalf("[Firebase] Error initializing Tenant Auth Client: %v\n", err)
		}
		TenantAuth = tenantAuthClient
		log.Printf("[Firebase] Tenant Auth Client Initialized")
	}

	log.Printf("[Firebase] Firebase Initialized")
}

func ValidateAuthHandler(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization", "")
	idToken := strings.TrimPrefix(authHeader, "Bearer ")
	if idToken == "" {
		responseBody := fiberutils.NewErrorResponseBody(http.StatusUnauthorized, "Missing authentication token", nil, ERROR_CODE_MISSING_AUTH_TOKEN)
		return c.Status(http.StatusUnauthorized).JSON(responseBody)
	}

	var verifyIdTokenFunc func(context.Context, string) (*auth.Token, error)

	if TenantAuth != nil {
		verifyIdTokenFunc = TenantAuth.VerifyIDToken
	} else if !Config.EnforceTenant {
		verifyIdTokenFunc = Auth.VerifyIDToken
	} else {
		tenantId := c.Get("X-Tenant-Id", "")
		if tenantId == "" {
			responseBody := fiberutils.NewErrorResponseBody(http.StatusUnauthorized, "Missing Tenant ID", nil, ERROR_CODE_MISSING_TENANT_ID)
			return c.Status(http.StatusUnauthorized).JSON(responseBody)
		}

		if slices.Contains(Config.BlacklistTenantIDs, tenantId) {
			responseBody := fiberutils.NewErrorResponseBody(http.StatusUnauthorized, "Invalid Tenant ID", nil, ERROR_CODE_INVALID_TENANT_ID)
			return c.Status(http.StatusUnauthorized).JSON(responseBody)
		}

		tenantAuthClient, err := Auth.TenantManager.AuthForTenant(tenantId)
		if err != nil {
			responseBody := fiberutils.NewErrorResponseBody(http.StatusUnauthorized, "Invalid Tenant ID", nil, ERROR_CODE_INVALID_TENANT_ID)
			return c.Status(http.StatusUnauthorized).JSON(responseBody)
		}
		verifyIdTokenFunc = tenantAuthClient.VerifyIDToken
	}

	token, err := verifyIdTokenFunc(context.Background(), idToken)
	if err != nil {
		if strings.Contains(err.Error(), "ID token has expired") {
			responseBody := fiberutils.NewErrorResponseBody(http.StatusUnauthorized, "Expired authentication token", nil, ERROR_CODE_EXPIRED_AUTH_TOKEN)
			return c.Status(http.StatusUnauthorized).JSON(responseBody)
		}
		responseBody := fiberutils.NewErrorResponseBody(http.StatusUnauthorized, "Invalid authentication token", nil, ERROR_CODE_INVALID_AUTH_TOKEN)
		return c.Status(http.StatusUnauthorized).JSON(responseBody)
	}

	authUser := AuthUser{
		UserID:        token.UID,
		Email:         token.Claims["email"].(string),
		EmailVerified: token.Claims["email_verified"].(bool),
		AuthTime:      token.AuthTime,
		TenantID:      token.Firebase.Tenant,
	}

	c.Locals(FIBER_AUTH_USER_KEY, authUser)

	return c.Next()
}

func CreateUser(user *auth.UserToCreate, tenantId string) (*auth.UserRecord, error) {
	var createUserFunc func(context.Context, *auth.UserToCreate) (*auth.UserRecord, error)

	if tenantId != "" {
		if slices.Contains(Config.BlacklistTenantIDs, tenantId) {
			return nil, ErrorInvalidTenantID
		}

		tenantAuthClient, err := Auth.TenantManager.AuthForTenant(tenantId)
		if err != nil {
			return nil, ErrorInvalidTenantID
		}
		createUserFunc = tenantAuthClient.CreateUser
	} else if TenantAuth != nil {
		createUserFunc = TenantAuth.CreateUser
	} else if !Config.EnforceTenant {
		createUserFunc = Auth.CreateUser
	} else {
		return nil, ErrorMissingTenantID
	}

	return createUserFunc(context.Background(), user)
}

func DeleteUser(uid string, tenantId string) error {
	var deleteUserFunc func(context.Context, string) error

	if tenantId != "" {
		if slices.Contains(Config.BlacklistTenantIDs, tenantId) {
			return ErrorInvalidTenantID
		}

		tenantAuthClient, err := Auth.TenantManager.AuthForTenant(tenantId)
		if err != nil {
			return ErrorInvalidTenantID
		}
		deleteUserFunc = tenantAuthClient.DeleteUser
	} else if TenantAuth != nil {
		deleteUserFunc = TenantAuth.DeleteUser
	} else if !Config.EnforceTenant {
		deleteUserFunc = Auth.DeleteUser
	} else {
		return ErrorMissingTenantID
	}

	return deleteUserFunc(context.Background(), uid)
}

func UpdateTenant(tenantId string, updateProps *auth.TenantToUpdate) error {
	if tenantId == "" {
		return ErrorMissingTenantID
	}

	if tenantId != "" {
		if slices.Contains(Config.BlacklistTenantIDs, tenantId) {
			return ErrorInvalidTenantID
		}
	}

	_, err := Auth.TenantManager.UpdateTenant(context.Background(), tenantId, updateProps)
	if err != nil {
		return err
	}

	return nil
}
