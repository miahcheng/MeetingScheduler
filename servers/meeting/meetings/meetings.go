package meetings

import (
	"encoding/json"
	"info441-finalproj/servers/gateway/models/users"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"

	"github.com/juliangruber/go-intersect"
)

func (c *Context) MeetingsHandler(w http.ResponseWriter, r *http.Request) {
	CheckAuth(w, r, c)
	if r.Method == "POST" {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Wrong content type, must be application/json", http.StatusUnsupportedMediaType)
		}
		data, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			http.Error(w, "Request body could not be read", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Method must be POST", http.StatusMethodNotAllowed)
	}
}

func (c *Context) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	CheckAuth(w, r, c)
	if r.Method == "GET" {
		ws := c.GetUserTimes(c.UserID, w)
		weekJSON, jsonErr := json.Marshal(ws)
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
		c.InsertHelper(tempWeek.Sunday, 1, w)
		c.InsertHelper(tempWeek.Monday, 2, w)
		c.InsertHelper(tempWeek.Tuesday, 3, w)
		c.InsertHelper(tempWeek.Wednesday, 4, w)
		c.InsertHelper(tempWeek.Thursday, 5, w)
		c.InsertHelper(tempWeek.Friday, 6, w)
		c.InsertHelper(tempWeek.Saturday, 7, w)

		ws := c.GetUserTimes(c.UserID, w)
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
	}
}

func (c *Context) SpecificMeetingHandler(w http.ResponseWriter, r *http.Request) {
	CheckAuth(w, r, c)
	meetingID, idErr := strconv.Atoi(path.Base(r.URL.Path))
	if idErr != nil {
		http.Error(w, "Invalid ID passed, cannot parse", http.StatusBadRequest)
	}
	if r.Method == "GET" {
		members := make([]int64, 0)
		result := &Meeting{}
		row := c.CalendarStore.QueryRow("SELECT MeetingName, MeetingDesc, CreatorID FROM Meeting WHERE MeetingID = ?", meetingID)
		if err := row.Scan(&result.MeetingName, &result.MeetingDesc, &result.CreatorID); err != nil {
			http.Error(w, "Database could not be queried", http.StatusInternalServerError)
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
		weeks := []Week{}
		for _, id := range members {
			ws := c.GetUserTimes(id, w)
			weeks = append(weeks, ws)
		}
		if len(weeks) > 0 {
			firstUser := weeks[0]
			for i := 1; i < len(weeks); i++ {
				firstUser.Sunday = intersect.Hash(firstUser.Sunday, weeks[i].Sunday).([]string)
				firstUser.Monday = intersect.Hash(firstUser.Monday, weeks[i].Monday).([]string)
				firstUser.Tuesday = intersect.Hash(firstUser.Tuesday, weeks[i].Tuesday).([]string)
				firstUser.Wednesday = intersect.Hash(firstUser.Wednesday, weeks[i].Wednesday).([]string)
				firstUser.Thursday = intersect.Hash(firstUser.Thursday, weeks[i].Thursday).([]string)
				firstUser.Friday = intersect.Hash(firstUser.Friday, weeks[i].Friday).([]string)
				firstUser.Saturday = intersect.Hash(firstUser.Saturday, weeks[i].Saturday).([]string)
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
		json.Unmarshal(data, &temp.userID)
		insq := "INSERT INTO MeetingMembers(UserID, MeetingID) VALUES(?,?)"
		c.CalendarStore.Exec(insq, temp.userID, meetingID)
		w.Header().Set("Content-Type", "string")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User added succesfully"))
	} else {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
}

func (c *Context) GetUserTimes(userID int64, w http.ResponseWriter) Week {
	ws := Week{}
	references := []*[]string{&ws.Sunday, &ws.Monday, &ws.Tuesday, &ws.Wednesday, &ws.Thursday, &ws.Friday, &ws.Saturday}
	rows, getErr := c.CalendarStore.Query("SELECT DayID, TimeStart FROM UserTimes WHERE UserID = ?", c.UserID)
	if getErr != nil {
		http.Error(w, "Could not find that user", http.StatusBadRequest)
		return Week{}
	}
	for rows.Next() {
		temp := &Holder{}
		if err := rows.Scan(&temp.dayID, &temp.timeString); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return Week{}
		}
		*references[temp.dayID] = append(*references[temp.dayID], temp.timeString)
	}
	return ws
}

func (c *Context) InsertHelper(times []string, dayID int, w http.ResponseWriter) {
	insq := "INSERT INTO UserTimes(UserID, TimeID, DayID) VALUES(?,?,?)"
	for _, timeStart := range times {
		var timeID int64
		row := c.CalendarStore.QueryRow("SELECT TimeID FROM [Time] WHERE TimeStart = ?", timeStart)
		if err := row.Scan(&timeID); err != nil {
			http.Error(w, "Error retreiving times", http.StatusBadRequest)
			return
		}
		c.CalendarStore.Exec(insq, c.UserID, timeID, dayID)
	}
}

// Helper function to check if the current user is authenticated. If so, the user's ID is saved
// in a context struct used to keep track of the session.
func CheckAuth(w http.ResponseWriter, r *http.Request, c *Context) {
	if r.Header.Get("X-User") == "" {
		http.Error(w, "User is not authenticated", http.StatusUnauthorized)
		return
	} else {
		userInfo := &users.User{}
		json.Unmarshal([]byte(r.Header.Get("X-User")), userInfo)
		c.UserID = userInfo.ID
	}
}
