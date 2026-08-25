package goldgym

type TokenRedis struct {
	UserID    int    `json:"user_id"`
	SessionId string `json:"session_id"`
	ExpiresAt int64  `json:"expires_at"`
	IP        string `json:"ip"`
	Device    string `json:"device"`
	Role      string `json:"role"`
}
