package firebaseauth

import (
	"context"
	"slices"

	"firebase.google.com/go/v4/auth"
)

func createUser(user *auth.UserToCreate, tenantId string) (*auth.UserRecord, error) {
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

func deleteUser(uid string, tenantId string) error {
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
