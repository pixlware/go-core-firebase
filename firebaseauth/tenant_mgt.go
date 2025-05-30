package firebaseauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"

	"firebase.google.com/go/v4/auth"
)

func updateTenant(tenantId string, updateProps *auth.TenantToUpdate) error {
	if tenantId == "" {
		return ErrorMissingTenantID
	}

	if slices.Contains(Config.BlacklistTenantIDs, tenantId) {
		return ErrorInvalidTenantID
	}

	_, err := Auth.TenantManager.UpdateTenant(context.Background(), tenantId, updateProps)
	if err != nil {
		return err
	}

	return nil
}

func updateTenantInheritance(tenantId string, inheritance *TenantInheritance) error {
	if tenantId == "" {
		return ErrorMissingTenantID
	}

	if slices.Contains(Config.BlacklistTenantIDs, tenantId) {
		return ErrorInvalidTenantID
	}

	token, err := IdentityToolkitCredentials.TokenSource.Token()
	if err != nil {
		log.Printf("[FirebaseAuth] Error generating token for tenant '%s': %v", tenantId, err)
		return ErrorGeneratingToken
	}

	jsonData, err := json.Marshal(IdentityPlatformTenantUpdater{
		Inheritance: inheritance,
	})
	if err != nil {
		log.Printf("[FirebaseAuth] Error marshaling request for tenant '%s': %v", tenantId, err)
		return ErrorMarshalingRequest
	}

	url := fmt.Sprintf(
		"https://identitytoolkit.googleapis.com/v2/projects/%s/tenants/%s?updateMask=inheritance.emailSendingConfig",
		Config.ProjectID,
		tenantId,
	)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[FirebaseAuth] Error creating request for tenant '%s': %v", tenantId, err)
		return ErrorCreatingRequest
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("[%s] Error reading response body: %v", resp.Status, err)
		}
		return fmt.Errorf("[%s] Response body: %s", resp.Status, string(bodyBytes))
	}

	return nil
}
