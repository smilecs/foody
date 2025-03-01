package requests

type UserReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	DOB      string `json:"date_of_birth"`
}
