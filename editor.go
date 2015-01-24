package main

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/b3log/wide/util"
	"github.com/gorilla/websocket"
)

func editorWSHandler(w http.ResponseWriter, r *http.Request) {
	httpSession, _ := httpSessionStore.Get(r, "coditor-session")
	userSession := httpSession.Values[user_session]

	if nil == userSession {
		http.Error(w, "Forbidden", http.StatusForbidden)

		return
	}

	sid := r.URL.Query()["sid"][0]

	cSession := coditorSessions.get(sid)
	if nil == cSession {
		return
	}

	conn, _ := websocket.Upgrade(w, r, nil, 1024, 1024)
	wsChan := util.WSChannel{Sid: sid, Conn: conn, Request: r, Time: time.Now()}

	ret := map[string]interface{}{"editor": "Editor initialized", "cmd": "init-editor"}
	err := wsChan.WriteJSON(&ret)
	if nil != err {
		return
	}

	editorWS[sid] = &wsChan

	logger.Tracef("Open a new [Editor] with session [%s], %d", sid, len(editorWS))

	input := map[string]interface{}{}

	for {
		if err := wsChan.ReadJSON(&input); err != nil {
			return
		}

		logger.Trace(input)

		docName := input["docName"].(string)
		docName = filepath.Clean(docName)

		doc := documentHolder.getDoc(docName)

		for _, cursor := range doc.cursors {
			if cursor.sid == sid { // skip the current session itself
				continue
			}

			delta := input["delta"].(map[string]interface{})
			ops := delta["ops"].([]interface{})

			offset := float64(0)
			for _, op := range ops {
				ope := op.(map[string]interface{})

				if retain, ok := ope["retain"]; ok {
					offset = retain.(float64)
				}

				if _, ok := ope["insert"]; ok {
					offset = offset + 1
				}
			}

			ret = map[string]interface{}{"output": input["delta"], "cmd": "text-change",
				"docName": docName, "user": input["user"].(string), "offset": offset, "color": input["color"].(string)}

			if err := editorWS[cursor.sid].WriteJSON(&ret); err != nil {
				logger.Error("[Editor Channel] ERROR: " + err.Error())

				return
			}
		}

		wsChan.Time = time.Now()
	}
}
