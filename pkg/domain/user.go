package domain

type User struct {
	ID               string `json:"id"`
	Username         string `json:"username"`
	Region           string `json:"region"`
	OrganizationType string `json:"organization_type"`
	Role             string `json:"role"`
	Contacts         struct {
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		TelegramID  string `json:"telegram_id"`
	} `json:"contacts"`
	Permissions struct {
		Email    bool `json:"email"`
		Phone    bool `json:"phone"`
		Telegram bool `json:"telegram"`
		Push     bool `json:"push"`
	}
}
