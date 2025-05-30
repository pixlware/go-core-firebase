package firebaseauth

import "errors"

const FIBER_AUTH_USER_KEY = "authUser"
const FIBER_TENANT_ID_HEADER_KEY = "X-Tenant-Id"

const (
	ERROR_CODE_MISSING_AUTH_TOKEN = "MISSING_AUTH_TOKEN"
	ERROR_CODE_EXPIRED_AUTH_TOKEN = "EXPIRED_AUTH_TOKEN"
	ERROR_CODE_INVALID_AUTH_TOKEN = "INVALID_AUTH_TOKEN"
	ERROR_CODE_MISSING_TENANT_ID  = "MISSING_TENANT_ID"
	ERROR_CODE_INVALID_TENANT_ID  = "INVALID_TENANT_ID"
)

var (
	ErrorInvalidTenantID            = errors.New("invalid tenant id")
	ErrorMissingTenantID            = errors.New("missing tenant id")
	ErrorGeneratingToken            = errors.New("error generating token")
	ErrorMarshalingRequest          = errors.New("error marshaling request")
	ErrorCreatingRequest            = errors.New("error creating request")
	ErrorUpdatingEmailSendingConfig = errors.New("error updating email sending configuration")
)
