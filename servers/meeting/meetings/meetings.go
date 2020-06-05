package meetings

import (
	"encoding/json"
	"errors"
	"fmt"
	"info441-finalproj/servers/gateway/models/users"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/juliangruber/go-intersect"
)

func (c *Context) MeetingsHandler(w http.ResponseWriter, r *http.Request) {
	err := CheckAuth(w, r, c)
	if err != nil {
		return
	}
	if r.Method == "POST" {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Wrong content type, must be application/json", http.StatusUnsupportedMediaType)
			return
		}
		data, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			http.Error(w, "Request body could not be read", http.StatusBadRequest)
			return
		}
		r.Body.Close()
		meeting := Meeting{}
		json.Unmarshal(data, &meeting)
		insq := "INSERT INTO Meeting(CreatorID, MeetingName, MeetingDesc) VALUES(?, ?, ?)"
		res, err := c.CalendarStore.Exec(insq, c.UserID, meeting.MeetingName, meeting.MeetingDesc)
		if err != nil {
			http.Error(w, "Error with connecting to database", http.StatusInternalServerError)
			return
		}
		cid, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Could not insert into database", http.StatusInternalServerError)
			return
		}
		query := "INSERT INTO MeetingMembers(UserID, MeetingID) VALUES(?, ?)"
		if len(meeting.Members) > 0 {
			for _, member := range meeting.Members {
				c.CalendarStore.Exec(query, member, cid)
			}
		}
		c.CalendarStore.Exec(query, c.UserID, cid)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Meeting created successfully"))
	} else {
		http.Error(w, "Method must be POST", http.StatusMethodNotAllowed)
		return
	}
}

func (c *Context) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	err := CheckAuth(w, r, c)
	if err != nil {
		return
	}
	if r.Method == "GET" {
		user := &User{}
		ws, err := c.GetUserTimes(c.UserID, w)
		if err != nil {
			return
		}
		user.Week = ws
		row := c.CalendarStore.QueryRow("SELECT Email, FirstName, LastName FROM Users WHERE UserID = ?", c.UserID)
		userErr := row.Scan(&user.Email, &user.FirstName, &user.LastName)
		if userErr != nil {
			http.Error(w, "Could not get user", http.StatusInternalServerError)
			return
		}
		meetings, meetErr := c.CalendarStore.Query("SELECT MeetingID FROM MeetingMembers WHERE UserID = ?", c.UserID)
		if meetErr != nil {
			http.Error(w, "Could not get meetings", http.StatusInternalServerError)
			return
		}
		meetingToAdd := make([]int64, 0)
		for meetings.Next() {
			holder := &Holder{}
			err := meetings.Scan(&holder.meetingID)
			if err != nil {
				http.Error(w, "Problem occurred when getting meetings", http.StatusInternalServerError)
				return
			}
			meetingToAdd = append(meetingToAdd, holder.meetingID)
		}
		user.Meetings = meetingToAdd
		weekJSON, jsonErr := json.Marshal(user)
		if jsonErr != nil {
			http.Error(w, "Data could not be returned", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(weekJSON)
	} else if r.Method == "PATCH" {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Wrong content type, must be application/json", http.StatusUnsupportedMediaType)
			return
		}
		data, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			http.Error(w, "Request body could not be read", http.StatusBadRequest)
			return
		}
		r.Body.Close()
		tempWeek := Week{}
		json.Unmarshal(data, &tempWeek)
		delete := "DELETE FROM UserTimes WHERE UserID = ?"
		c.CalendarStore.Exec(delete, c.UserID)
		c.InsertHelper(tempWeek.Sunday, "Sunday", w)
		c.InsertHelper(tempWeek.Monday, "Monday", w)
		c.InsertHelper(tempWeek.Tuesday, "Tuesday", w)
		c.InsertHelper(tempWeek.Wednesday, "Wednesday", w)
		c.InsertHelper(tempWeek.Thursday, "Thursday", w)
		c.InsertHelper(tempWeek.Friday, "Friday", w)
		c.InsertHelper(tempWeek.Saturday, "Saturday", w)

		ws, err := c.GetUserTimes(c.UserID, w)
		if err != nil {
			return
		}
		weekJSON, jsonErr := json.Marshal(ws)
		if jsonErr != nil {
			http.Error(w, "Data could not be returned", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(weekJSON)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func (c *Context) SpecificMeetingHandler(w http.ResponseWriter, r *http.Request) {
	err := CheckAuth(w, r, c)
	if err != nil {
		return
	}
	vars := mux.Vars(r)
	meetingID := vars["id"]
	if r.Method == "GET" {
		members := make([]int64, 0)
		result := &Meeting{}
		row := c.CalendarStore.QueryRow("SELECT MeetingName, MeetingDesc, CreatorID FROM Meeting WHERE MeetingID = ?", meetingID)
		if err := row.Scan(&result.MeetingName, &result.MeetingDesc, &result.CreatorID); err != nil {
			http.Error(w, "Meeting does not exist", http.StatusBadRequest)
			return
		}
		rows, queryErr := c.CalendarStore.Query("SELECT UserID FROM MeetingMembers WHERE MeetingID = ?", meetingID)
		if queryErr != nil {
			http.Error(w, "Could not get members", http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			temp := &Holder{}
			rows.Scan(&temp.userID)
			members = append(members, temp.userID)
		}
		result.Members = members
		weeks := make([]Week, len(members))
		for i, id := range members {
			ws, err := c.GetUserTimes(id, w)
			if err != nil {
				return
			}
			weeks[i] = ws
		}
		if len(weeks) > 0 {
			firstUser := weeks[0]
			for i := 1; i < len(weeks); i++ {
				firstUser.Sunday = InterfaceToString(intersect.Hash(firstUser.Sunday, weeks[i].Sunday))
				firstUser.Monday = InterfaceToString(intersect.Hash(firstUser.Monday, weeks[i].Monday))
				firstUser.Tuesday = InterfaceToString(intersect.Hash(firstUser.Tuesday, weeks[i].Tuesday))
				firstUser.Wednesday = InterfaceToString(intersect.Hash(firstUser.Wednesday, weeks[i].Wednesday))
				firstUser.Thursday = InterfaceToString(intersect.Hash(firstUser.Thursday, weeks[i].Thursday))
				firstUser.Friday = InterfaceToString(intersect.Hash(firstUser.Friday, weeks[i].Friday))
				firstUser.Saturday = InterfaceToString(intersect.Hash(firstUser.Saturday, weeks[i].Saturday))
			}
			result.Sunday = firstUser.Sunday
			result.Monday = firstUser.Monday
			result.Tuesday = firstUser.Tuesday
			result.Wednesday = firstUser.Wednesday
			result.Thursday = firstUser.Thursday
			result.Friday = firstUser.Friday
			result.Saturday = firstUser.Saturday
		}
		meetingJSON, jsonErr := json.Marshal(result)
		if jsonErr != nil {
			http.Error(w, "Data could not be returned", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(meetingJSON)
	} else if r.Method == "POST" {
		data, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			http.Error(w, "Request body could not be read", http.StatusBadRequest)
			return
		}
		temp := &Holder{}
		json.Unmarshal(data, &temp)
		row := c.CalendarStore.QueryRow("SELECT UserID FROM Users WHERE Email = ?", temp.Email)
		scanErr := row.Scan(&temp.userID)
		if scanErr != nil {
			http.Error(w, "Error reading from database", http.StatusInternalServerError)
			return
		}
		insq := "INSERT INTO MeetingMembers(UserID, MeetingID) VALUES(?,?)"
		_, insErr := c.CalendarStore.Exec(insq, temp.userID, meetingID)
		if insErr != nil {
			http.Error(w, "Problem inserting into database", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "string")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User added successfully"))
	} else {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
}

func (c *Context) GetUserTimes(userID int64, w http.ResponseWriter) (Week, error) {
	ws := Week{}
	references := []*[]string{&ws.Sunday, &ws.Monday, &ws.Tuesday, &ws.Wednesday, &ws.Thursday, &ws.Friday, &ws.Saturday}
	rows, getErr := c.CalendarStore.Query("SELECT DayID, TimeStart FROM UserTimes UT INNER JOIN `Time` T ON UT.TimeID = T.TimeID WHERE UserID = ?", userID)
	if getErr != nil {
		http.Error(w, "No free times have been added", http.StatusBadRequest)
		return Week{}, errors.New("Could not contact database")
	}
	for rows.Next() {
		temp := &Holder{}
		if err := rows.Scan(&temp.dayID, &temp.timeString); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return Week{}, errors.New("Database error")
		}
		*references[temp.dayID-1] = append(*references[temp.dayID-1], temp.timeString)
	}
	for i := 0; i < len(references); i++ {
		if *references[i] == nil {
			*references[i] = make([]string, 0)
		}
	}
	return ws, nil
}

func (c *Context) InsertHelper(times []string, dayName string, w http.ResponseWriter) {
	insq := "INSERT INTO UserTimes(UserID, TimeID, DayID) VALUES(?,?,?)"
	dayRow := c.CalendarStore.QueryRow("SELECT DayID FROM `Day` WHERE DayName = ?", dayName)
	var dayID int64
	err := dayRow.Scan(&dayID)
	if err != nil {
		http.Error(w, "Error retrieving day", http.StatusBadRequest)
		return
	}
	for _, timeStart := range times {
		var timeID int64
		row := c.CalendarStore.QueryRow("SELECT TimeID FROM `Time` WHERE TimeStart = ?", timeStart)
		if err := row.Scan(&timeID); err != nil {
			http.Error(w, "Error retrieving times", http.StatusBadRequest)
			return
		}
		c.CalendarStore.Exec(insq, c.UserID, timeID, dayID)
	}
}

// Helper function to check if the current user is authenticated. If so, the user's ID is saved
// in a context struct used to keep track of the session.
func CheckAuth(w http.ResponseWriter, r *http.Request, c *Context) error {
	if r.Header.Get("X-User") == "" {
		http.Error(w, "User is not authenticated", http.StatusUnauthorized)
		return errors.New("Unauthenticated user")
	}
	userInfo := &users.User{}
	json.Unmarshal([]byte(r.Header.Get("X-User")), userInfo)
	c.UserID = userInfo.ID

	return nil
}

func InterfaceToString(interfaces []interface{}) []string {
	strings := make([]string, 0)
	for _, i := range interfaces {
		strings = append(strings, fmt.Sprintf("%v", i))
	}
	return strings
}
