package innards

type LoginRequest struct {
	Id       string `json:"id" form:"id"`
	Password string `json:"password" form:"password"`
}

type SignupRequest struct {
	Id       string `json:"id" form:"id"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type LinkRequest struct {
	Alias string `json:"alias"`
	Url   string `json:"url"`
}
