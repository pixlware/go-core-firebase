package firebaseauth

import (
	"context"
	"log"
	"net/http"
	"slices"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/pixlware/go-core-fiber/fiberutils"
)

func validateAuthHandler(c *fiber.Ctx) error {
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
		tenantId := c.Get(FIBER_TENANT_ID_HEADER_KEY, "")
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
			log.Printf("[FirebaseAuth] Error getting tenant auth client for tenant '%s': %v", tenantId, err)
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
