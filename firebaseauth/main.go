package firebaseauth

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"golang.org/x/oauth2/google"
)

var (
	App                        *firebase.App
	Auth                       *auth.Client
	TenantAuth                 *auth.TenantClient
	IdentityToolkitCredentials *google.Credentials
)

var UserManager = &UserService{
	Create: createUser,
	Delete: deleteUser,
}

var TenantManager = &TenantService{
	Create:            createTenant,
	Update:            updateTenant,
	UpdateInheritance: updateTenantInheritance,
}

var HandlersManager = &HandlerService{
	ValidateAuth: validateAuthHandler,
}

func init() {
	log.Printf("[FirebaseAuth] Initializing FirebaseAuth...")
	ctx := context.Background()

	firebaseConfig := &firebase.Config{
		ProjectID: Config.ProjectID,
	}
	firebaseApp, err := firebase.NewApp(context.Background(), firebaseConfig)
	if err != nil {
		log.Fatalf("[FirebaseAuth] Error initializing Firebase App: %v\n", err)
	}
	App = firebaseApp
	log.Printf("[FirebaseAuth] Firebase App Initialized")

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("[FirebaseAuth] Error initializing Auth Client: %v\n", err)
	}
	Auth = authClient
	log.Printf("[FirebaseAuth] Auth Client Initialized")

	if Config.EnforceTenant && Config.TenantID != "" {
		tenantAuthClient, err := Auth.TenantManager.AuthForTenant(Config.TenantID)
		if err != nil {
			log.Fatalf("[FirebaseAuth] Error initializing Tenant Auth Client: %v\n", err)
		}
		TenantAuth = tenantAuthClient
		log.Printf("[FirebaseAuth] Tenant Auth Client Initialized")
	}

	credentials, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/identitytoolkit")
	if err != nil {
		log.Fatalf("[FirebaseAuth] Error initializing Identity Platform Credentials: %v\n", err)
	}
	IdentityToolkitCredentials = credentials

	log.Printf("[FirebaseAuth] FirebaseAuth Initialized")
}
