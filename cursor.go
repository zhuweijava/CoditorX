package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/b3log/wide/util"
)

type Cursor struct {
	sid      string
	offset   int
	username string
	color    string
}

func setCursorHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"succ": true}
	defer util.RetJSON(w, r, data)

	httpSession, _ := httpSessionStore.Get(r, "coditor-session")
	userSession := httpSession.Values[user_session]
	if nil == userSession {
		data["succ"] = false
		data["msg"] = "permission denied"
		return
	}

	var args map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		logger.Error(err)
		data["succ"] = false
		return
	}

	sid := args["sid"].(string)
	docName := args["docName"].(string)
	offset := args["offset"].(float64)
	color := args["color"].(string)
	user := userSession.(*User)

	cursor := &Cursor{sid: sid, offset: int(offset), username: user.Username, color: color}

	docName = filepath.Clean(docName)

	doc := documentHolder.getDoc(docName)
	doc.addCursor(cursor)
}

func (doc *Document) addCursor(cursor *Cursor) {
	doc.cursorLock.Lock()
	defer doc.cursorLock.Unlock()

	doc.cursors = append(doc.cursors, cursor)
}

func (doc *Document) removeCursor(cursorId string) {
	doc.cursorLock.Lock()
	defer doc.cursorLock.Unlock()

	var newCursors []*Cursor
	for _, c := range doc.cursors {
		if c.sid != cursorId { // in case of dupilicated id, remove them all
			newCursors = append(newCursors, c)
		}
	}

	doc.cursors = newCursors
}
