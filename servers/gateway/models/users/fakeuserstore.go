package users

import "time"

// FakeStore is used for testing handler functions
// and the User Store
type FakeStore interface {
	//GetByID returns the User with the given ID
	GetByID(id int64) (*User, error)

	//GetByEmail returns the User with the given email
	GetByEmail(email string) (*User, error)

	//GetByUserName returns the User with the given Username
	GetByUserName(username string) (*User, error)

	//Insert inserts the user into the database, and returns
	//the newly-inserted User, complete with the DBMS-assigned ID
	Insert(user *User) (*User, error)

	//Update applies UserUpdates to the given user ID
	//and returns the newly-updated user
	Update(id int64, updates *Updates) (*User, error)

	//Delete deletes the user with the given ID
	Delete(id int64) error

	// InsertSignIN inserts into the sign in table to track
	// user sign in attempts
	InsertSignIn(user *User, signinTime time.Time, ipAddy string) (*User, error)
}

// FakeUserStore blahs
type FakeUserStore struct {
	User *User
	Err  error
}

// NewFakeUserStore blah
func NewFakeUserStore(usr *User) *FakeUserStore {
	return &FakeUserStore{
		User: usr,
		Err:  nil,
	}
}

// GetByID blah
func (fakeUser *FakeUserStore) GetByID(ID int64) (*User, error) {
	if ID != fakeUser.User.ID {
		return InvalidUser, ErrUserNotFound
	}
	return fakeUser.User, nil
}

// GetByEmail blah
func (fakeUser *FakeUserStore) GetByEmail(email string) (*User, error) {
	if email != fakeUser.User.Email {
		return InvalidUser, ErrUserNotFound
	}
	return fakeUser.User, nil
}

// GetByUserName blah
func (fakeUser *FakeUserStore) GetByUserName(username string) (*User, error) {
	if username != fakeUser.User.UserName {
		return InvalidUser, ErrUserNotFound
	}
	return fakeUser.User, nil
}

// Insert blah
func (fakeUser *FakeUserStore) Insert(user *User) (*User, error) {
	user.ID = 1
	fakeUser.User = user
	return fakeUser.User, nil
}

// Update blah
func (fakeUser *FakeUserStore) Update(id int64, updates *Updates) (*User, error) {
	if id == fakeUser.User.ID {
		up := fakeUser.User.ApplyUpdates(updates)
		if up != nil {
			return nil, up
		}
		return fakeUser.User, nil
	}
	return InvalidUser, ErrUserNotFound
}

// Delete blah
func (fakeUser *FakeUserStore) Delete(id int64) error {
	if id == fakeUser.User.ID {
		fakeUser.User = InvalidUser
		return nil
	}
	return ErrUserNotFound
}

// InsertSignIn blah
func (fakeUser *FakeUserStore) InsertSignIn(user *User, signinTime time.Time, ipAddy string) (*User, error) {
	return user, nil
}
