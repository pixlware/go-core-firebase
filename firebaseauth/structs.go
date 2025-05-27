package firebaseauth

type AuthUser struct {
	UserID        string `json:"userId"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	AuthTime      int64  `json:"authTime"`
	TenantID      string `json:"tenantId"`
}
