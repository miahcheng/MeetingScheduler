package meetings

import "database/sql"

type Week struct {
	Sunday    []string
	Monday    []string
	Tuesday   []string
	Wednesday []string
	Thursday  []string
	Friday    []string
	Saturday  []string
}

type Meeting struct {
	MeetingName string
	MeetingDesc string
	CreatorID   int64
	Members     []int64
	Sunday      []string
	Monday      []string
	Tuesday     []string
	Wednesday   []string
	Thursday    []string
	Friday      []string
	Saturday    []string
}

type Holder struct {
	userID     int64
	dayID      int64
	meetingID  int64
	timeString string
	Email      string
}

type Context struct {
	UserID        int64
	CalendarStore *sql.DB
}

type User struct {
	Email     string
	FirstName string
	LastName  string
	Meetings  []int64
	Week      Week
}
