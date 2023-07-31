package kratos

import "time"

type Claims struct {
	ID                          string    `json:"id"`
	Active                      bool      `json:"active"`
	ExpiresAt                   time.Time `json:"expires_at"`
	AuthenticatedAt             time.Time `json:"authenticated_at"`
	AuthenticatorAssuranceLevel string    `json:"authenticator_assurance_level"`
	AuthenticationMethods       []struct {
		Method      string    `json:"method"`
		Aal         string    `json:"aal"`
		CompletedAt time.Time `json:"completed_at"`
	} `json:"authentication_methods"`
	IssuedAt time.Time `json:"issued_at"`
	Identity struct {
		ID             string    `json:"id"`
		SchemaID       string    `json:"schema_id"`
		SchemaURL      string    `json:"schema_url"`
		State          string    `json:"state"`
		StateChangedAt time.Time `json:"state_changed_at"`
		Traits         struct {
			Name      string `json:"name"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
		} `json:"traits"`
		VerifiableAddresses []struct {
			ID        string    `json:"id"`
			Value     string    `json:"value"`
			Verified  bool      `json:"verified"`
			Via       string    `json:"via"`
			Status    string    `json:"status"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"verifiable_addresses"`
		RecoveryAddresses []struct {
			ID        string    `json:"id"`
			Value     string    `json:"value"`
			Via       string    `json:"via"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"recovery_addresses"`
		MetadataPublic interface{} `json:"metadata_public"`
		CreatedAt      time.Time   `json:"created_at"`
		UpdatedAt      time.Time   `json:"updated_at"`
	} `json:"identity"`
}

type Error struct {
	Error struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Reason  string `json:"reason"`
		Message string `json:"message"`
	} `json:"error"`
}
