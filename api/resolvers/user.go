package resolvers

import (
	"context"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"imabad.dev/do/api/handlers"
	messages "imabad.dev/do/api/messages"
	"imabad.dev/do/lib/models"
	"imabad.dev/do/lib/utils"
)

type createUserArgs struct {
	Username string
	Email    string
	Password string
}

type loginUserArgs struct {
	Email    string
	Password string
}

type updateUserArgs struct {
	Username         *string
	Email            *string
	Password         *string
	PreviousPassword *string
}

type verifyEmailArgs struct {
	Code *string
}

type SignupResponse struct {
	token string
}

func (r *SignupResponse) Token() string {
	return r.token
}

type VerifyEmailResponse struct {
	verified bool
}

func (r *VerifyEmailResponse) Verified() bool {
	return r.verified
}

type LoginResponse struct {
	token string
}

func (r *LoginResponse) Token() string {
	return r.token
}

func (r *RootResolver) CreateUser(args createUserArgs) (*SignupResponse, error) {
	if len(args.Password) <= 0 {
		return nil, fmt.Errorf("missing password")
	}
	if len(args.Email) <= 0 {
		return nil, fmt.Errorf("missing email address")
	} else {
		var matched, err = regexp.Match("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])", []byte(args.Email))
		if err != nil {
			return nil, err
		}
		if !matched {
			return nil, fmt.Errorf("invalid email address")
		}
	}
	if len(args.Username) <= 0 {
		return nil, fmt.Errorf("missing username")
	}
	var count int
	r.Db.Model(&models.User{}).Where("email = ?", args.Email).Count(&count)
	if count > 0 {
		return nil, fmt.Errorf("user already exists with email address")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(args.Password), 14)
	if err != nil {
		return nil, err
	}
	user := models.User{
		Username: args.Username,
		Email:    args.Email,
		Role:     "user",
		Password: string(bytes),
	}
	if err := r.Db.Create(&user).Error; err != nil {
		return nil, err
	}
	code := utils.GenerateRandomString(8)
	emailVerify := models.EmailVerification{
		Email: args.Email,
		Code:  code,
	}
	if err := r.Db.Create(&emailVerify).Error; err != nil {
		return nil, err
	}
	err = messages.SendEmail(messages.EmailRequest{
		To:      args.Email,
		Subject: "DynamicOverlay - Please Verify Email Address",
		Content: fmt.Sprintf("Please verify your email address by clicking the link below<br/><br/><a href=\"%s/verify?code=%s\">Verify Email</a>", utils.GetConfig().APIURL, code),
	})
	if err != nil {
		fmt.Println("Error sending email")
	}
	userID := fmt.Sprintf("%v", user.ID)
	token, err := utils.SignToken(&userID)
	if err == nil {
		return &SignupResponse{token: utils.StrValue(token)}, nil
	} else {
		return nil, err
	}
}

func (r *RootResolver) Login(args loginUserArgs) (*LoginResponse, error) {
	if len(args.Email) <= 0 || len(args.Password) <= 0 {
		return nil, fmt.Errorf("incorrect details")
	}
	var user models.User
	r.Db.First(&user, "email = ?", args.Email)
	if &user == nil {
		return nil, fmt.Errorf("user does not exist")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(args.Password))
	if err != nil {
		return nil, fmt.Errorf("incorrect details")
	}
	userID := fmt.Sprintf("%v", user.ID)
	token, err := utils.SignToken(&userID)
	if err == nil {
		return &LoginResponse{token: utils.StrValue(token)}, nil
	} else {
		return nil, err
	}
}

func (r *RootResolver) UpdateUser(ctx context.Context, args updateUserArgs) (*models.UserResolver, error) {
	if ctx.Value(handlers.ContextKey("UserID")) == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	userID := ctx.Value(handlers.ContextKey("UserID")).(*int)
	var user models.User
	err := r.Db.First(&user, utils.IntValue(userID)).Error
	if err != nil {
		return nil, err
	}
	if &user != nil {
		if args.Username != nil {
			if utils.StrValue(args.Username) != user.Username {
				user.Username = utils.StrValue(args.Username)
			}
		}
		if args.Email != nil {
			if utils.StrValue(args.Email) != user.Email {
				var matched, err = regexp.Match("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])", []byte(utils.StrValue(args.Email)))
				if err != nil {
					return nil, err
				}
				if !matched {
					return nil, fmt.Errorf("invalid email address")
				}
				var count int
				r.Db.Model(&models.User{}).Where("email = ?", args.Email).Count(&count)
				if count > 0 {
					return nil, fmt.Errorf("user already exists with email address")
				}
				user.Email = utils.StrValue(args.Email)
			}
		}
		if args.Password != nil {
			if args.PreviousPassword == nil {
				return nil, fmt.Errorf("updating your password requires your current password")
			}
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(utils.StrValue(args.PreviousPassword)))
			if err != nil {
				return nil, fmt.Errorf("incorrect password")
			}
			newPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(utils.StrValue(args.Password)), 14)
			if err != nil {
				return nil, fmt.Errorf("failed to hash new password %e", err)
			}
			user.Password = string(newPasswordBytes)
		}
		r.Db.Save(&user)
		return &models.UserResolver{U: user}, nil
	}
	return nil, nil
}

func (r *RootResolver) User(ctx context.Context) (*models.UserResolver, error) {
	if ctx.Value(handlers.ContextKey("UserID")) == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	userID := ctx.Value(handlers.ContextKey("UserID")).(*int)
	var user models.User
	err := r.Db.First(&user, utils.IntValue(userID)).Error
	if err != nil {
		return nil, err
	}
	if &user != nil {
		return &models.UserResolver{U: user}, nil
	}
	return nil, nil
}

func (r *RootResolver) VerifyEmail(ctx context.Context, args verifyEmailArgs) (*VerifyEmailResponse, error) {
	if ctx.Value(handlers.ContextKey("UserID")) == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	userID := ctx.Value(handlers.ContextKey("UserID")).(*int)
	var user models.User
	err := r.Db.First(&user, utils.IntValue(userID)).Error
	if err != nil {
		return nil, err
	}
	if &user != nil && !user.EmailVerified {
		var emailVerification models.EmailVerification
		err := r.Db.Find(&emailVerification, "email = ? AND code = ?", user.Email, args.Code).Error
		if err != nil {
			return nil, err
		}
		if &emailVerification != nil {
			r.Db.Delete(&emailVerification)
			user.EmailVerified = true
			r.Db.Save(&user)
			return &VerifyEmailResponse{verified: true}, nil
		} else {
			return &VerifyEmailResponse{verified: false}, nil
		}
	}
	return nil, fmt.Errorf("unable to verify")
}
