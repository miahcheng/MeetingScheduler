package meetings

import "database/sql"

type Context struct {
	UserID        int64
	CalendarStore *sql.DB
}
