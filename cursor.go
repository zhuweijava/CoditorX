package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/b3log/wide/util"
)

type Cursor struct {
	Sid      string `json:"sid"`
	Offset   int    `json:"offset"`
	Username string `json:"username"`
	Color    string `json:"color"`
	Email    string `json:"email"`
	Md5Email string `json:"md5Email"`
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

	md5Email := toMd5(user.Email)

	cursor := &Cursor{Sid: sid, Offset: int(offset), Username: user.Username, Color: color, Email: user.Email, Md5Email: md5Email}

	docName = filepath.Clean(docName)

	doc := documentHolder.getDoc(docName)
	doc.addCursor(cursor)
}

func listCursorsHandler(w http.ResponseWriter, r *http.Request) {
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

	docName := args["docName"].(string)
	doc := documentHolder.getDoc(docName)
	if doc == nil {
		data["succ"] = false
		data["msg"] = "Can not find the document."
		return
	}

	cursors := doc.getCursors()
	data["cursors"] = cursors
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
		if c.Sid != cursorId { // in case of dupilicated id, remove them all
			newCursors = append(newCursors, c)
		}
	}

	doc.cursors = newCursors
}

func (doc *Document) getCursors() []*Cursor {
	doc.cursorLock.Lock()
	defer doc.cursorLock.Unlock()

	return doc.cursors
}
