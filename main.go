// Copyright (c) 2015, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	confIP := flag.String("ip", "", "this will overwrite Wide.IP if specified")
	confPort := flag.String("port", "", "this will overwrite Wide.Port if specified")
	confChannel := flag.String("channel", "", "this will overwrite Wide.Channel if specified")

	flag.Parse()

	loadConf(*confIP, *confPort, *confChannel)
	loadLocales()
	InitDocumentHolder()

	// TODO: fixedTimeXXX

	// static resources
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	serveSingle("/favicon.ico", "./static/favicon.ico")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/share", shareHandler)
	http.HandleFunc("/shareInfo", getShareInfoHandler)

	http.HandleFunc("/session/ws", coditorSessionWSHandler)
	http.HandleFunc("/notification/ws", notificationWSHandler)
	http.HandleFunc("/editor/ws", editorWSHandler)

	http.HandleFunc("/file/new", fileNew)
	http.HandleFunc("/file/del", fileDel)
	http.HandleFunc("/file/rename", fileRename)

	http.HandleFunc("/files", fileTreeHandler)
	http.HandleFunc("/shares", shareListHandler)

	http.HandleFunc("/doc/open", openDocHandler)
	http.HandleFunc("/doc/listCursors", listCursorsHandler)
	http.HandleFunc("/doc/setCursor", setCursorHandler)
	http.HandleFunc("/doc/commit", commitDocHandler)
	http.HandleFunc("/doc/fetch", fetchDocHandler)
	// TODO: 可能不需要
	http.HandleFunc("/doc/getHead", getHeadDocHandler)

	logger.Info(conf.Server)

	err := http.ListenAndServe(conf.Server, nil)
	if err != nil {
		logger.Error(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if conf.Context+"/" != r.RequestURI {
		http.NotFound(w, r)

		return
	}

	httpSession, _ := httpSessionStore.Get(r, "coditor-session")
	userSession := httpSession.Values[user_session]

	if userSession != nil {
		rand.Seed(time.Now().UnixNano())
		sid := strconv.Itoa(rand.Int())
		cSession := coditorSessions.new(httpSession, sid)
		user := userSession.(*User)
		model := map[string]interface{}{"session": cSession, "workspace": user.getWorkspace(),
			"pathSeparator": string(os.PathSeparator)}
		model["currentuser"] = user
		model["gravatar"] = toMd5(user.Email)

		toHtml(w, "coditor.html", model, user.Locale)

		return
	}

	model := map[string]interface{}{}
	toHtml(w, "login.html", model, conf.Locale)
}

func toHtml(w http.ResponseWriter, file string, model map[string]interface{}, locale string) {

	model["i18n"] = getMsgs(locale)
	model["conf"] = conf
	t, err := template.ParseFiles("views/" + file)
	if nil != err {
		logger.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	t.Execute(w, model)
}

func toJson(w http.ResponseWriter, model map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(model)
	if nil != err {
		logger.Error(err)
		return
	}

	w.Write(data)
}

// serveSingle registers the handler function for the given pattern and filename.
func serveSingle(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}
