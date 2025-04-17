package users

type UserType string

const (
	UserTypeFree       UserType = "free"
	UserTypeSubscribed UserType = "subscribed"
	UserTypeAdmin      UserType = "admin"
)
