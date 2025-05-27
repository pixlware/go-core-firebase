package firebaseauth

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type FirebaseConfig struct {
	ProjectID          string
	TenantID           string
	EnforceTenant      bool
	BlacklistTenantIDs []string
}

var Config FirebaseConfig = FirebaseConfig{
	ProjectID:          "",
	EnforceTenant:      false,
	TenantID:           "",
	BlacklistTenantIDs: []string{},
}

func init() {
	env := getEnv("ENV", "default")
	var envFilePath string
	if env == "default" {
		return
	} else {
		envFilePath = ".env." + env
	}

	err := godotenv.Load(envFilePath)
	if err != nil {
		return
	}

	Config.ProjectID = getEnv("GOOGLE_CLOUD_PROJECT", Config.ProjectID)
	Config.EnforceTenant = getEnv("FIREBASE_ENFORCE_TENANT", "false") == "true"
	Config.TenantID = getEnv("FIREBASE_TENANT_ID", Config.TenantID)
	Config.BlacklistTenantIDs = strings.Split(getEnv("FIREBASE_BLACKLIST_TENANT_IDS", ""), ",")
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
