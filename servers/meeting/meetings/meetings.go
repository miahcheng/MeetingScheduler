package meetings

import (
	"encoding/json"
	"info441-finalproj/servers/gateway/models/users"
	"io/ioutil"
	"net/http"
)

func (c *Context) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	CheckAuth(w, r, c)
	ws := Week{}
	references := []*[]string{&ws.Sunday, &ws.Monday, &ws.Tuesday, &ws.Wednesday, &ws.Thursday, &ws.Friday, &ws.Saturday}
	if r.Method == "GET" {
		rows, getErr := c.CalendarStore.Query("SELECT DayID, TimeStart FROM UserTimes WHERE UserID = ?", c.UserID)
		if getErr != nil {
			http.Error(w, "Could not find that user", http.StatusBadRequest)
		}
		for rows.Next() {
			temp := &Holder{}
			if err := rows.Scan(&temp.dayID, &temp.timeString); err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			*references[temp.dayID] = append(*references[temp.dayID], temp.timeString)
		}
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

		rows, getErr := c.CalendarStore.Query("SELECT DayID, TimeStart FROM UserTimes WHERE UserID = ?", c.UserID)
		if getErr != nil {
			http.Error(w, "Could not find that user", http.StatusBadRequest)
		}
		for rows.Next() {
			temp := &Holder{}
			if err := rows.Scan(&temp.dayID, &temp.timeString); err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			*references[temp.dayID] = append(*references[temp.dayID], temp.timeString)
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
	}
}

func (c *Context) SpecificMeetingHandler(w http.ResponseWriter, r *http.Request) {

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
