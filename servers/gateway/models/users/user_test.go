package users

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// Test the (nu *NewUser) Validate() function to ensure it catches all possible validation errors,
// and returns no error when the new user is valid.
func TestValidate(t *testing.T) {
	cases := []struct {
		input       NewUser
		expectedErr bool
		// output error
	}{
		{
			// good case
			NewUser{
				Email:        "hello@uw.edu",
				Password:     "abcdefghi",
				PasswordConf: "abcdefghi",
				UserName:     "dubsthehusky",
				FirstName:    "Dubs",
				LastName:     "Husky"},
			false,
			// nil,
		},
		{
			// bad email
			NewUser{
				Email:        "1930323924248",
				Password:     "abcdefghi",
				PasswordConf: "abcdefghi",
				UserName:     "dubsthehusky",
				FirstName:    "Dubs",
				LastName:     "Husky"},
			true,
		},
		{
			// bad password
			NewUser{
				Email:        "hello@uw.edu",
				Password:     "ab",
				PasswordConf: "abcdefghi",
				UserName:     "dubsthehusky",
				FirstName:    "Dubs",
				LastName:     "Husky"},
			true,
		},
		{
			// bad conf password
			NewUser{
				Email:        "hello@uw.edu",
				Password:     "abcdefghi",
				PasswordConf: "abcdefgh",
				UserName:     "dubsthehusky",
				FirstName:    "Dubs",
				LastName:     "Husky"},
			true,
		},
		{
			// 0 length username
			NewUser{
				Email:        "hello@uw.edu",
				Password:     "abcdefghi",
				PasswordConf: "abcdefghi",
				UserName:     "",
				FirstName:    "Dubs",
				LastName:     "Husky"},
			true,
		},
		{
			// spaces in username
			NewUser{
				Email:        "hello@uw.edu",
				Password:     "abcdefghi",
				PasswordConf: "abcdefghi",
				UserName:     "du bsth ehusky ",
				FirstName:    "Dubs",
				LastName:     "Husky"},
			true,
		},
		{
			// just a space username
			NewUser{
				Email:        "hello@uw.edu",
				Password:     "abcdefghi",
				PasswordConf: "abcdefghi",
				UserName:     "  ",
				FirstName:    "Dubs",
				LastName:     "Husky"},
			true,
		},
	}

	for _, c := range cases {
		err := c.input.Validate()
		if !c.expectedErr && err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}
}

// Test to user
func TestToUser(t *testing.T) {
	cases := []struct {
		newUsr NewUser
		usr    User
	}{
		{
			NewUser{
				Email:        "MyEmailAddress@example.com ",
				UserName:     "hellohello",
				Password:     "helloworld",
				PasswordConf: "helloworld",
			},
			User{
				Email:    "MyEmailAddress@example.com ",
				UserName: "hellohello",
				PhotoURL: gravatarBasePhotoURL + "0bc83cb571cd1c50ba6f3e8a78ef1346",
			},
		},
	}
	for _, c := range cases {
		newPerson, err := c.newUsr.ToUser()
		if err != nil {
			t.Errorf("Invalid new user, error: %v", err)
		}
		if newPerson.PhotoURL != c.usr.PhotoURL {
			t.Errorf("Error creating hash for photourl, got: %s", newPerson.PhotoURL)
		}
		passErr := bcrypt.CompareHashAndPassword(newPerson.PassHash, []byte(c.newUsr.Password))
		if passErr != nil {
			t.Errorf("Error in creating password hash")
		}
	}
}

// Test Full Name
func TestFullName(t *testing.T) {
	cases := []struct {
		user     User
		expected string
	}{
		{
			// normal
			User{
				FirstName: "Hellen",
				LastName:  "Wang",
			},
			"Hellen Wang",
		},
		{
			// No first name
			User{
				LastName: "Wu",
			},
			"Wu",
		},
		{
			// No last name
			User{
				FirstName: "Wes",
			},
			"Wes",
		},
		{
			// no name
			User{},
			"",
		},
	}

	for _, c := range cases {
		name := c.user.FullName()
		if name != c.expected {
			t.Errorf("Full name is not correct")
		}
	}
}

// Test Authenticate
func TestAuth(t *testing.T) {
	cases := []struct {
		usr         NewUser
		inputPass   string
		expectedErr bool
	}{
		{
			// normal
			NewUser{
				Email:        "abcd@uw.edu",
				Password:     "helloworld",
				PasswordConf: "helloworld",
				UserName:     "bobthebuilder",
				FirstName:    "Bob",
				LastName:     "Builder",
			},
			"helloworld",
			false,
		},
		{
			// incorrect password
			NewUser{
				Email:        "abcd@uw.edu",
				Password:     "helloworld",
				PasswordConf: "helloworld",
				UserName:     "bobthebuilder",
				FirstName:    "Bob",
				LastName:     "Builder",
			},
			"heloworld",
			true,
		},
		{
			// empty string
			NewUser{
				Email:        "abcd@uw.edu",
				Password:     "helloworld",
				PasswordConf: "helloworld",
				UserName:     "bobthebuilder",
				FirstName:    "Bob",
				LastName:     "Builder",
			},
			"",
			true,
		},
	}
	for _, c := range cases {
		newUse, _ := c.usr.ToUser()
		authErr := newUse.Authenticate(c.inputPass)
		if !c.expectedErr && authErr != nil {
			t.Errorf("Unexpected error: %v", authErr)
		}
	}
}

func TestApplyUpdates(t *testing.T) {
	cases := []struct {
		user     User
		update   Updates
		expected string
	}{
		// good
		{
			User{
				FirstName: "John",
				LastName:  "Oliver",
			},
			Updates{
				FirstName: "Jonathan",
				LastName:  "Twist",
			},
			"Jonathan Twist",
		},
		{
			// invalid first name update
			User{
				FirstName: "Liam",
				LastName:  "Smith",
			},
			Updates{
				LastName: "Louie",
			},
			"Liam Louie",
		},
		{
			// invalid last name update
			User{
				FirstName: "Emily",
				LastName:  "Lee",
			},
			Updates{
				FirstName: "Emilia",
			},
			"Emilia Lee",
		},
		{
			// no user name
			User{},
			Updates{
				FirstName: "Angela",
				LastName:  "Smith",
			},
			"Angela Smith",
		},
	}
	for _, c := range cases {
		err := c.user.ApplyUpdates(&c.update)
		if err != nil {
			t.Errorf("Update not applied, error: %v", err)
		}
		fn := c.user.FullName()
		if fn != c.expected {
			t.Errorf("Update not applied correctly, expected %s, got %s ", c.expected, fn)
		}

	}
}
