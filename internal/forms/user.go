package forms

// UserRegisterForm - User register form struct
type UserRegisterForm struct {
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
	Email     string `bson:"email" json:"email"`
	Password  string `bson:"password" json:"password"`
}

// UserLoginForm - User login form struct
type UserLoginForm struct {
	Email    string
	Password string
}

// ForgotPasswordForm - Forgot password form struct
type ForgotPasswordForm struct {
	Email string
}
