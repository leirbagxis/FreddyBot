package types

type ConnectStartRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type ConnectVerifyRequest struct {
	Code string `json:"code" binding:"required"`
}

type Connect2FARequest struct {
	Password string `json:"password" binding:"required"`
}

type ConnectStatusResponse struct {
	Connected bool   `json:"connected"`
	PhoneHash string `json:"phoneHash,omitempty"`
	UserID    int64  `json:"userId"`
}
