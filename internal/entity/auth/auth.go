package auth

// UserID ...
type UserID struct {
	UserID int `json:"userid"`
}

// Metadata ...
type Metadata struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

// Err ...
type Err struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	Code   int    `json:"code"`
}

// Auth ...
type Auth struct {
	Data     interface{} `json:"data"`
	Metadata Metadata    `json:"metadata"`
	Error    Err         `json:"error"`
}

// -----------------------------------------------------------------------
type ClaimKey string

type Token struct {
	AccessToken         string `json:"access_token"`
	ExpiresIn           int64  `json:"expires_in"`
	ExpiresAt           int64  `json:"expires_at"`
	TokenType           string `json:"token_type"`
	ForceChangePassword int    `json:"force_change_password"`
}

type LoginRequest struct {
	NIP      string `json:"nip"`
	Password string `json:"password"`
}

type AccessListRequest struct {
	PTID  int    `json:"pt_id"`
	AppID int    `json:"app_id"`
	NIP   string `json:"nip"`
}

// -----------------------------------------------------------------------
