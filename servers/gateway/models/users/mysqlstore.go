package users

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// InvalidUser is returned when there's an error finding user in database
var InvalidUser = &User{}

// MySQLStore represents a user in MySQL
type MySQLStore interface {
	// GetByID returns the User associated with the given ID
	GetByID(ID int64) (*User, error)

	// GetByEmail returns the User with given emaiil
	GetByEmail(email string) (*User, error)

	// GetByUserName returns the User with given username
	GetByUserName(username string) (*User, error)

	// Insert inserts new user into the database and returns a copy of the
	// user with the ID value filled
	Insert(user *User) (*User, error)

	// Update updates the user by finding with ID
	Update(ID int64, updates *Updates) (*User, error)

	// Delete removes user with the ID and returns error
	Delete(ID int64) error
}

// SQLStore is the struct that has the abilities of mysqlstore?
type SQLStore struct {
	Db *sql.DB
}

func GetNewStore(db *sql.DB) *SQLStore {
	return &SQLStore{db}
}

// GetByID gets the User by given ID value
func (database SQLStore) GetByID(ID int64) (*User, error) {
	info, err := database.Db.Query("SELECT ID, Email, PassHash, UserName, FirstName, LastName, FROM Users WHERE ID = ?", ID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	usr := User{}
	for info.Next() {
		scanErr := info.Scan(&usr.ID, &usr.Email, &usr.PassHash, &usr.FirstName,
			&usr.LastName)
		if scanErr != nil {
			return nil, scanErr
		}
	}
	//
	return &usr, nil
}

// GetByEmail gets the User by given email value
func (database SQLStore) GetByEmail(email string) (*User, error) {
	info, err := database.Db.Query("SELECT ID, Email, PassHash, UserName, FirstName, LastName, FROM Users WHERE Email = ?", email)
	if err != nil {
		return nil, fmt.Errorf("unexpected error querying db: %v", err)
	}
	usr := User{}
	for info.Next() {
		scanErr := info.Scan(&usr.ID, &usr.Email, &usr.PassHash, &usr.FirstName, &usr.LastName)
		if scanErr != nil {
			return nil, scanErr
		}
	}
	usr.Email = email
	return &usr, nil
}

// Insert adds a new user into the database
func (database SQLStore) Insert(user *User) (*User, error) {
	trx, errTR := database.Db.Begin()
	if errTR != nil {
		return nil, fmt.Errorf("Error beginning transaction: %v", errTR)
	}
	insertQ := "INSERT INTO Users(Email, PassHash, UserName, FirstName, LastName) VALUES (?,?,?,?,?,?)"
	q, errQ := trx.Prepare(insertQ)
	if errQ != nil {
		return InvalidUser, fmt.Errorf("error preparing")
	}
	defer q.Close()
	result, err := q.Exec(user.Email, user.PassHash, user.FirstName, user.LastName)
	if err != nil {
		trx.Rollback()
		return InvalidUser, fmt.Errorf("Error inserting row")
	}
	id, idErr := result.LastInsertId()
	if idErr != nil {
		return InvalidUser, fmt.Errorf("Error getting new ID: %v", idErr)
	}
	// refactor this
	InsertByDay(database, trx, 1, user.Sunday, id)
	InsertByDay(database, trx, 2, user.Monday, id)
	InsertByDay(database, trx, 3, user.Tuesday, id)
	InsertByDay(database, trx, 4, user.Wednesday, id)
	InsertByDay(database, trx, 5, user.Thursday, id)
	InsertByDay(database, trx, 6, user.Friday, id)
	InsertByDay(database, trx, 7, user.Saturday, id)
	user.ID = id
	trx.Commit()
	return user, nil
}

// InsertSignIn inserts into the sign in table
func (database SQLStore) InsertSignIn(user *User, signinTime time.Time, ipAddy string) (*User, error) {
	trx, errTR := database.Db.Begin()
	if errTR != nil {
		return nil, fmt.Errorf("Error beginning transaction: %v", errTR)
	}
	insertQ := "INSERT INTO SignIns(UserID, SignInDate, IPAddress) VALUES (?,?,?)"
	q, errQ := trx.Prepare(insertQ)
	if errQ != nil {
		return InvalidUser, fmt.Errorf("error preparing")
	}
	defer q.Close()
	result, err := q.Exec(user.ID, signinTime, ipAddy)
	if err != nil {
		trx.Rollback()
		return InvalidUser, fmt.Errorf("Error inserting row")
	}
	_, idErr := result.LastInsertId()
	if idErr != nil {
		return InvalidUser, fmt.Errorf("Error gettign new ID: %v", idErr)
	}
	trx.Commit()
	return user, nil
}

// Update updates's a users information
func (database SQLStore) Update(ID int64, updates *Updates) (*User, error) {
	trx, errTR := database.Db.Begin()
	if errTR != nil {
		return nil, fmt.Errorf("error beginning transaction %v", errTR)
	}
	curr, err := database.GetByID(ID)
	if err != nil {
		return nil, err
	}
	upErr := curr.ApplyUpdates(updates)
	if upErr != nil {
		return nil, upErr
	}

	updateQ := "UPDATE Users SET FirstName = ?, LastName = ? WHERE ID = ?"
	q, errQ := trx.Prepare(updateQ)
	if errQ != nil {
		return InvalidUser, fmt.Errorf("error preparing")
	}
	defer q.Close()
	_, sqlErr := q.Exec(updates.FirstName, updates.LastName, ID)
	if sqlErr != nil {
		trx.Rollback()
		return nil, fmt.Errorf("error applying update to user: %v", sqlErr)
	}
	trx.Commit()
	return curr, nil
}

// Delete deletes user from database
func (database SQLStore) Delete(ID int64) error {
	deleteQ := "DELETE FROM Users WHERE ID = ?"
	_, err := database.Db.Exec(deleteQ, ID)
	if err != nil {
		return errors.New("Error Deleting User")
	}
	return nil
}

func InsertByDay(db SQLStore, trans *sql.Tx, dayID int, times []string, userID int64) error {
	for _, time := range times {
		timeID, err := db.Db.Query("SELECT TimeID FROM [Time] WHERE TimeRange = ?", time)
		if err != nil {
			return err
		}
		insQ := "INSERT INTO UserTimes(DayID, TimeID, UserID) VALUES(?, ?, ?)"
		q, qerr := trans.Prepare(insQ)
		if qerr != nil {
			return qerr
		}
		defer q.Close()
		_, serr := q.Exec(dayID, timeID, userID)
		if serr != nil {
			return serr
		}
	}
	return nil
}
