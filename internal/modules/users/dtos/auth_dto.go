package users

type OAuthRequestDTO struct {
	Provider string `json:"provider"`
}

type SendOTPViaSMSRequestDTO struct {
	Phonenumber string `json:"phone_number" exampe:"966512345678"`
}

type SendOTPViaEmailRequestDTO struct {
	Email string `json:"email" exampe:"mshari@alhaimah.sa"`
}

type IdentifierType string

const (
	IdentifierEmail IdentifierType = "email"
	IdentifierPhone IdentifierType = "phone_number"
)

type VerifyOTPRequestDTO struct {
	Identifier IdentifierType `json:"identifier" example:"email"`
	Type       string         `json:"type" example:"mshari@alhaimah.sa"`
	OTP        int            `json:"otp" example:"8402"`
}

type SendSMSRequest struct {
	To      []string `json:"to"`
	Message string   `json:"message"`
	Token   string   `json:"token"`
}
