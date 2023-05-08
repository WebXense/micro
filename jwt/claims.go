package jwt

type Claims struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	IP        string `json:"ip"`
	TokenType string `json:"token_type"`
}

// NewClaims creates a new Claims.
func NewClaims(id uint, username, ip string, tokenType string) *Claims {
	return &Claims{
		ID:        id,
		Username:  username,
		IP:        ip,
		TokenType: tokenType,
	}
}
