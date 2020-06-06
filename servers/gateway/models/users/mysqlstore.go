package users

import (
	"database/sql"
	"time"
)

type Database struct {
	DB *sql.DB
}

func GetNewStore(db *sql.DB) *Database {
	return &Database{db}
}

func (db *Database) GetByID(id int64) (*User, error) {
	row := db.DB.QueryRow("SELECT * FROM Users WHERE UserID = ?", id)
	user := User{}
	if err := row.Scan(&user.ID, &user.Email, &user.PassHash, &user.FirstName, &user.LastName); err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (db *Database) GetByEmail(email string) (*User, error) {
	row := db.DB.QueryRow("SELECT * FROM Users WHERE Email = ?", email)
	user := User{}
	if err := row.Scan(&user.ID, &user.Email, &user.PassHash, &user.FirstName, &user.LastName); err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (db *Database) Insert(user *User) (*User, error) {
	insq := "INSERT INTO Users(Email, PassHash, FirstName, LastName) VALUES (?,?,?,?)"
	res, err := db.DB.Exec(insq, user.Email, user.PassHash, user.FirstName, user.LastName)
	if err != nil {
		return nil, err
	}
	id, err2 := res.LastInsertId()
	if err2 != nil {
		return nil, err2
	}
	user.ID = id
	return user, nil
}

func (db *Database) Update(id int64, updates *Updates) (*User, error) {
	insq := "UPDATE Users SET FirstName = ?, LastName = ? WHERE UserID = ?"
	_, err := db.DB.Exec(insq, updates.FirstName, updates.LastName, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	user, err2 := db.GetByID(id)
	if err2 != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (db *Database) Delete(id int64) error {
	insq := "DELETE FROM Users WHERE UserID = ?"
	_, err := db.DB.Exec(insq, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) TrackLogin(id int64, ip string, time time.Time) error {
	query := "INSERT INTO SignIns(UserID, IPAddress, SignInDate) VALUES (?,?,?)"
	_, err := db.DB.Exec(query, id, ip, time)
	if err != nil {
		return err
	}
	return nil
}
