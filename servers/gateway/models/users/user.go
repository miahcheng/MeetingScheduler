package users

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//gravatarBasePhotoURL is the base URL for Gravatar image requests.
//See https://id.gravatar.com/site/implement/images/ for details
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"-"` //never JSON encoded/decoded
	PassHash  []byte `json:"-"` //never JSON encoded/decoded
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	email, emErr := mail.ParseAddress(nu.Email)
	if emErr != nil {
		return fmt.Errorf("Invalid Email: got %s", email)
	}

	if len(nu.Password) < 6 {
		return fmt.Errorf("Invalid Password: Password is not long enough")
	}

	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("Invalid Conf Password: Confirmation password does not match password")
	}

	if len(nu.UserName) == 0 {
		return fmt.Errorf("Invalid Username: username needs to be at least 1 character")
	}

	if strings.ContainsAny(nu.UserName, " ") {
		return fmt.Errorf("Invalid Username: Username cannot have spaces")
	}

	return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	valErr := nu.Validate()
	if valErr != nil {
		return nil, valErr
	}
	email := nu.Email
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	hash := md5.New()
	hash.Write([]byte(email))
	hashEmail := hex.EncodeToString(hash.Sum(nil))

	pURL := gravatarBasePhotoURL + hashEmail
	// SetPassword(usr.PassHash)
	usr := &User{Email: nu.Email, UserName: nu.UserName, FirstName: nu.FirstName,
		LastName: nu.LastName, PhotoURL: pURL}
	usr.SetPassword(nu.Password)
	return usr, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put between the names. If both are missing,
//this returns an empty string
func (u *User) FullName() string {
	fName := u.FirstName
	lName := u.LastName
	if fName == "" && lName == "" {
		return ""
	} else if fName == "" {
		return lName
	} else if lName == "" {
		return fName
	}
	return fName + " " + lName
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	u.PassHash = hash
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	if password == "" {
		return fmt.Errorf("Error: Please Enter a password")
	}
	p := []byte(password)
	// log.Println(u.PassHash)
	err := bcrypt.CompareHashAndPassword(u.PassHash, p)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("passwords don't match")
	}
	return nil
}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	upF := updates.FirstName
	upL := updates.LastName
	if len(upF) == 0 && len(upL) == 0 {
		return fmt.Errorf("Updates cannot be of length 0")
	}
	if len(upF) > 0 && u.FirstName != upF {
		u.FirstName = upF
	}
	if len(upL) > 0 && u.LastName != upL {
		u.LastName = upL
	}

	return nil
}
