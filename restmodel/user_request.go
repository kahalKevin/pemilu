package restmodel

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AddUserRequest struct {
	Name      string `json:"name"`
	Tingkat   string `json:"tingkat"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type Response struct {
	Result    bool   `json:"result"`
	Role      string `json:"role"`
	Username  string `json:"username"`
	Tingkat   string `json:"tingkat"`
}

type ResponseGetUser struct {
	IDCalon  string `json:"idCalon"`
	Name     string `json:"name"`
	Tingkat  string `json:"tingkat"`
}

type ResponseGeneral struct {
	Result    bool   `json:"result"`	
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
