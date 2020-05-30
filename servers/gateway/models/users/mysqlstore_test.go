package users

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestGetByID is a test function for the SQLStore's GetByID
func TestGetByID(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		idToGet      int64
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1,
			false,
		},
		{
			"User Not Found",
			&User{},
			2,
			true,
		},
		{
			"User With Large ID Found",
			&User{
				1234567890,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1234567890,
			false,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()
		mainSQLStore := &SQLStore{db}
		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := "SELECT ID, Email, PassHash, UserName, FirstName, LastName, PhotoURL FROM Users WHERE ID = ?"

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnError(ErrUserNotFound)

			// Test GetByID()
			user, err := mainSQLStore.GetByID(c.idToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnRows(row)

			// Test GetByID()
			user, err := mainSQLStore.GetByID(c.idToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetByEmail(t *testing.T) {
	cases := []struct {
		name         string
		expectedUser *User
		emailToGet   string
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"test@test.com",
			false,
		},
		{
			"User not Found",
			&User{},
			"test@test.com",
			true,
		},
	}

	for _, c := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()
		mainSQLStore := &SQLStore{db}
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := "SELECT ID, Email, PassHash, UserName, FirstName, LastName, PhotoURL FROM Users WHERE Email = ?"
		if c.expectError {
			mock.ExpectQuery(query).WithArgs(c.emailToGet).WillReturnError(ErrUserNotFound)

			user, err := mainSQLStore.GetByEmail(c.emailToGet)
			if user != nil || err == nil {
				t.Errorf("Expected Error %v, but got %v", ErrUserNotFound, err)
			}
		} else {
			mock.ExpectQuery(query).WithArgs(c.emailToGet).WillReturnRows(row)

			user, err := mainSQLStore.GetByEmail(c.emailToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test[%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test[%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expecations: %s", err)
		}
	}

}

func TestGetByUserName(t *testing.T) {
	cases := []struct {
		name         string
		expectedUser *User
		usernToGet   string
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"username",
			false,
		},
		{
			"User not found",
			&User{},
			"username",
			true,
		},
		{
			"Username not found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"usrname",
			true,
		},
	}
	for _, c := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("there was a problem opening a db connection: [%v", err)
		}
		defer db.Close()
		mainSQLStore := &SQLStore{db}
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := "SELECT ID, Email, PassHash, UserName, FirstName, LastName, PhotoURL FROM Users WHERE UserName = ?"
		if c.expectError {
			mock.ExpectQuery(query).WithArgs(c.usernToGet).WillReturnError(ErrUserNotFound)

			user, err := mainSQLStore.GetByUserName(c.usernToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error %v, but got %v", ErrUserNotFound, err)
			}
		} else {
			mock.ExpectQuery(query).WithArgs(c.usernToGet).WillReturnRows(row)
			user, err := mainSQLStore.GetByUserName(c.usernToGet)
			if err != nil {
				t.Errorf("Unexpected error on test[%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test[%s]", c.name)
			}
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expecations: %s", err)
		}
	}
}

func TestInsert(t *testing.T) {
	cases := []struct {
		name         string
		expectedUser *User
		expectError  bool
	}{
		{
			"Correct Case",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"hello",
				"world",
				"photourl",
			},
			false,
		},
		{
			"Empty fields Case",
			&User{},
			true,
		},
	}
	for _, c := range cases {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("there was a problem opening a database connection: %v", err)
		}
		defer db.Close()
		mainSQLStore := &SQLStore{db}
		mock.ExpectBegin()
		query := "INSERT INTO Users(Email, PassHash, UserName, FirstName, LastName, PhotoURL) VALUES (?,?,?,?,?,?)"
		if c.expectError {
			mock.ExpectPrepare(query).ExpectExec().WithArgs(c.expectedUser.Email,
				c.expectedUser.PassHash,
				c.expectedUser.UserName,
				c.expectedUser.FirstName,
				c.expectedUser.LastName,
				c.expectedUser.PhotoURL).WillReturnError(errors.New("Error inserting row"))
			mock.ExpectRollback()
			user, err := mainSQLStore.Insert(c.expectedUser)
			if user != InvalidUser || err == nil {
				t.Errorf("Expected error %v, but got %v", errors.New("Error inserting row"), err)
			}
		} else {
			mock.ExpectPrepare(query)
			mock.ExpectExec(query).WithArgs(c.expectedUser.Email,
				c.expectedUser.PassHash,
				c.expectedUser.UserName,
				c.expectedUser.FirstName,
				c.expectedUser.LastName,
				c.expectedUser.PhotoURL).WillReturnResult(sqlmock.NewResult(c.expectedUser.ID, 1))
			mock.ExpectCommit()
			user, err := mainSQLStore.Insert(c.expectedUser)
			if err != nil {
				t.Errorf("error was not expected: %v", err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("error, invalid match in test %s", c.name)
			}
		}
		if err2 := mock.ExpectationsWereMet(); err2 != nil {
			t.Errorf("there were unfulfilled expecations: %s", err2)
		}
	}

}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name         string
		startUser    *User
		expectedUser *User
		updates      *Updates
		idToGet      int64
		expectError  bool
	}{
		{
			"Correct case",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"first",
				"last",
				"photourl",
			},
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"hello",
				"world",
				"photourl",
			},
			&Updates{
				"hello",
				"world",
			},
			1,
			false,
		},
		{
			"Incorrect update",
			&User{},
			&User{
				2,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"hello",
				"world",
				"photourl",
			},
			&Updates{
				"first",
				"there",
			},
			3,
			true,
		},
	}

	for _, c := range cases {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()
		mainSQLStore := &SQLStore{db}
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.startUser.ID,
			c.startUser.Email,
			c.startUser.PassHash,
			c.startUser.UserName,
			c.startUser.FirstName,
			c.startUser.LastName,
			c.startUser.PhotoURL,
		)
		query := "UPDATE Users SET FirstName = ? AND LastName = ? WHERE ID = ?"
		getUser := "SELECT ID, Email, PassHash, UserName, FirstName, LastName, PhotoURL FROM Users WHERE ID = ?"
		mock.ExpectBegin()
		if c.expectError {
			mock.ExpectQuery(getUser).WithArgs(c.idToGet).WillReturnRows(row)
			mock.ExpectPrepare(query)
			mock.ExpectExec(query).WithArgs(c.updates.FirstName, c.updates.LastName, c.idToGet).WillReturnError(errors.New("error applying update to user"))
			mock.ExpectRollback()
			user, err := mainSQLStore.Update(c.idToGet, c.updates)
			if user != nil || err == nil {
				t.Errorf("expected %v, got %v", ErrUserNotFound, err)
			}
		} else {
			mock.ExpectQuery(getUser).WithArgs(c.idToGet).WillReturnRows(row)
			mock.ExpectPrepare(query)
			mock.ExpectExec(query).WithArgs(c.updates.FirstName, c.updates.LastName, c.idToGet).WillReturnResult(sqlmock.NewResult(c.idToGet, 1))
			mock.ExpectCommit()

			user, err := mainSQLStore.Update(c.idToGet, c.updates)

			if err != nil {
				t.Errorf("Unexpected error on test %s: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("error, invalid match in test %s", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name           string
		givenUser      *User
		idToGet        int64
		expectedResult error
		expectError    bool
	}{
		{
			"Correct Case",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"hello",
				"world",
				"photourl",
			},
			1,
			nil,
			false,
		},
		{
			"Invalid ID num",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"hello",
				"world",
				"photourl",
			},
			2,
			errors.New("Error Deleting User"),
			true,
		},
	}

	for _, c := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: %v", err)
		}
		defer db.Close()
		mainSQLStore := &SQLStore{db}
		query := "DELETE FROM Users WHERE ID = ?"
		if c.expectError {
			mock.ExpectExec(query).WithArgs(c.idToGet).WillReturnError(c.expectedResult)
			err := mainSQLStore.Delete(c.idToGet)
			if err == nil {
				t.Errorf("expected error, but got nil")
			}
		} else {
			mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
			delErr := mainSQLStore.Delete(c.idToGet)

			if delErr != nil {
				t.Errorf("unexpected error on successful test: %v", delErr)
			}

			_, getErr := mainSQLStore.GetByID(c.idToGet)
			if getErr == nil {
				t.Errorf("unexpected error, user was not deleted")
			}
		}

		if err2 := mock.ExpectationsWereMet(); err2 != nil {
			t.Errorf("There were unfulfilled expectations: %s", err2)
		}

	}
}
