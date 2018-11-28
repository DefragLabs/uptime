package forms

// UserRegisterForm - User register form struct
type UserRegisterForm struct {
	FirstName   string `bson:"firstName" json:"firstName"`
	LastName    string `bson:"lastName" json:"lastName"`
	PhoneNumber string `bson:"phoneNumber" json:"phoneNumber"`
	CompanyName string `bson:"CompanyName" json:"CompanyName"`
	Email       string `bson:"email" json:"email"`
	Password    string `bson:"password" json:"password"`
}

// Validate user registration form.
func (userRegisterForm UserRegisterForm) Validate() string {
	if userRegisterForm.FirstName == "" {
		return "first name is required"
	} else if userRegisterForm.LastName == "" {
		return "last name is required"
	} else if userRegisterForm.Email == "" {
		return "email is required"
	} else if userRegisterForm.Password == "" {
		return "password is required"
	} else if userRegisterForm.CompanyName == "" {
		return "company name is required"
	}

	return ""
}

// UserLoginForm - User login form struct
type UserLoginForm struct {
	Email    string
	Password string
}

// Validate user login form.
func (userLoginForm UserLoginForm) Validate() string {
	if userLoginForm.Email == "" {
		return "email is required"
	} else if userLoginForm.Password == "" {
		return "password is required"
	}

	return ""
}

// ForgotPasswordForm - Forgot password form struct
type ForgotPasswordForm struct {
	Email string
}

// Validate forgot password form
func (forgotPasswordForm ForgotPasswordForm) Validate() string {
	if forgotPasswordForm.Email == "" {
		return "email is required"
	}

	return ""
}

// ResetPasswordForm is used for reset password
type ResetPasswordForm struct {
	UID         string `json:"uid"`
	Code        string `json:"code"`
	NewPassword string `json:"newPassword"`
}

// Validate reset password form.
func (resetPasswordForm ResetPasswordForm) Validate() string {
	if resetPasswordForm.UID == "" {
		return "uid is required"
	} else if resetPasswordForm.Code == "" {
		return "code is required"
	} else if resetPasswordForm.NewPassword == "" {
		return "newPassword is required"
	}

	return ""
}
