package restmodel

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Result    bool   `json:"result"`
	Role      string `json:"role"`
	Username  string `json:"username"`
	Tingkat   string `json:"tingkat"`
}

type RegisterRequest struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Msisdn   string `json:"msisdn"`
	Username string `json:"username"`
	Password string `json:"password"`
	Status   int    `json:"status"`
	Role     int    `json:"role"`
}
