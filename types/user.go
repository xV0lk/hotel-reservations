package types

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost  = 12
	minFNameLen = 2
	minLNameLen = 2
	minPassLen  = 7
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string             `bson:"firstName" json:"firstName"`
	LastName  string             `bson:"lastName" json:"lastName"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
}

type NewUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (params NewUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < minFNameLen {
		errors["firstName"] = fmt.Sprintf("first name must be at least %d characters long", minFNameLen)
	}
	if len(params.LastName) < minLNameLen {
		errors["lastName"] = fmt.Sprintf("last name must be at least %d characters long", minLNameLen)
	}
	if err := ValidatePassword(params.Password); len(err) > 0 {
		errors["password"] = strings.Join(err, ", ")
	}
	if err := ValidateEmail(params.Email); err != nil {
		errors["email"] = err.Error()
	}

	return errors
}

func ValidateEmail(email string) error {
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
		return errors.New("invalid email address")
	}
	return nil
}

func ValidatePassword(password string) []string {
	errors := []string{}
	if len(password) < minPassLen {
		errors = append(errors, fmt.Sprintf("password must be at least %d characters long", minPassLen))
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		errors = append(errors, "password must contain at least one lowercase letter")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		errors = append(errors, "password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		errors = append(errors, "password must contain at least one number")
	}
	if !regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password) {
		errors = append(errors, "password must contain at least one special character")
	}
	return errors
}

func NewUserFromParams(params *NewUserParams) (*User, error) {
	cryptPass, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Password:  string(cryptPass),
	}, nil
}

func (params UpdateUserParams) ToBson() bson.M {
	bson := bson.M{}
	if params.FirstName != "" {
		bson["firstName"] = params.FirstName
	}
	if params.LastName != "" {
		bson["lastName"] = params.LastName
	}
	return bson
}

func (params UpdateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if params.FirstName != "" && len(params.FirstName) < minFNameLen {
		errors["firstName"] = fmt.Sprintf("first name must be at least %d characters long", minFNameLen)
	}
	if params.LastName != "" && len(params.LastName) < minLNameLen {
		errors["lastName"] = fmt.Sprintf("last name must be at least %d characters long", minLNameLen)
	}
	return errors
}

// CheckUserBody checks if the json body contains the correct keys
func (params UpdateUserParams) CheckBody(jsonBody map[string]any) error {
	paramsMap := params.ToBson()
	for key := range jsonBody {
		_, ok := paramsMap[key]
		if !ok {
			return fmt.Errorf("the key '%s' is not a valid update property", key)
		}
	}
	if len(jsonBody) == 0 {
		return errors.New("no valid user properties were provided")
	}
	return nil
}
