package meetings

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (c *Context) MeetingsHandler(w http.ResponseWriter, r *http.Request) {
	//CheckAuth(w, r, c)
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
	} else {
		http.Error(w, "Method must be POST", http.StatusMethodNotAllowed)
		return
	}
}
