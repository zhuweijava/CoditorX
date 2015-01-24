package main

import (
	"encoding/json"
	"github.com/b3log/wide/util"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	EDITABLE = 1 // editable
	VIEWABLE = 0 // viewable
)

type Share struct {
	Owner     string `json:"owner"`
	DocName   string `json:"docName"`
	ShareType int    `json:"shareType"` // 0 - view, 1 - edit
}

func shareHandler(w http.ResponseWriter, r *http.Request) {
	if "GET" == r.Method {
		model := map[string]interface{}{}
		toHtml(w, "share.html", model, conf.Locale)
		return
	} else if "POST" == r.Method {
		data := map[string]interface{}{"succ": true}
		defer util.RetJSON(w, r, data)

		httpSession, _ := httpSessionStore.Get(r, "coditor-session")
		userSession := httpSession.Values[user_session]
		if nil == userSession {
			data["succ"] = false
			data["msg"] = "permission denied"
			return
		}

		user := userSession.(*User)

		var args map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
			logger.Error(err)
			data["succ"] = false
			data["msg"] = "args decode error!"
			return
		}

		fileName := args["fileName"].(string)
		editorsStr := ""
		if args["editors"] != nil {
			editorsStr = args["editors"].(string)
		}
		viewersStr := ""
		if args["viewers"] != nil {
			viewersStr = args["viewers"].(string)
		}
		isPublic := 0
		if args["isPublic"] != nil {
			isPublic = int(args["isPublic"].(float64))
		}

		doc := documentHolder.getDoc(fileName)
		if doc == nil {
			data["succ"] = false
			data["msg"] = "File Not exist!"
			return
		}
		// get old editors and old viewers.To del the invalids.
		// check permission
		oldEditors, err := doc.getEditors(user.Username)
		if err != nil {
			data["succ"] = false
			data["msg"] = err.Error()
			return
		}
		oldViewers, err := doc.getViewers(user.Username)
		if err != nil {
			data["succ"] = false
			data["msg"] = err.Error()
			return
		}

		// only get the file name
		index := strings.LastIndex(fileName, "/")
		if index > -1 {
			fileName = fileName[index+1:]
		}

		tempEditors := strings.Split(editorsStr, ",")
		tempViewers := strings.Split(viewersStr, ",")

		// check user exist and save to share.json
		editors := []string{}
		viewers := []string{}
		for _, editor := range tempEditors {
			if editor == "" {
				continue
			}
			share := &Share{}
			share.Owner = user.Username
			share.DocName = fileName
			share.ShareType = EDITABLE
			editors, _ = checkAndSave(editor, editors, share)
		}
		for _, viewer := range tempViewers {
			if viewer == "" {
				continue
			}
			share := &Share{}
			share.Owner = user.Username
			share.DocName = fileName
			share.ShareType = VIEWABLE
			viewers, _ = checkAndSave(viewer, viewers, share)
		}

		delEditors := []string{}
		delViewers := []string{}
		for _, oldEditor := range oldEditors {
			delAble := true
			for _, tempEditor := range tempEditors {
				if oldEditor == tempEditor {
					delAble = false
					break
				}
			}
			if delAble {
				delEditors = append(delEditors, oldEditor)
			}
		}
		for _, oldViewer := range oldViewers {
			delAble := true
			for _, tempViewer := range tempViewers {
				if oldViewer == tempViewer {
					delAble = false
					break
				}
			}
			if delAble {
				delViewers = append(delViewers, oldViewer)
			}
		}
		share := &Share{}
		share.Owner = user.Username
		share.DocName = fileName
		share.ShareType = VIEWABLE
		logger.Debugf("%s cancle share file %s to %v.\n", user.Username, fileName, delEditors)
		for _, delEditor := range delEditors {
			checkAndDel(delEditor, share)
		}
		logger.Debugf("%s cancle share file %s to %v.\n", user.Username, fileName, delViewers)
		for _, delViewer := range delViewers {
			checkAndDel(delViewer, share)
		}

		err = doc.setIsPublic(isPublic, user.Username)
		if err != nil {
			data["succ"] = false
			data["msg"] = err.Error()
			return
		}
		logger.Debugf("%s share file %s to %v.\n", user.Username, fileName, editors)
		err = doc.setEditors(editors, user.Username)
		if err != nil {
			data["succ"] = false
			data["msg"] = err.Error()
			return
		}
		logger.Debugf("%s share file %s to %v.\n", user.Username, fileName, viewers)
		err = doc.setViewers(viewers, user.Username)
		if err != nil {
			data["succ"] = false
			data["msg"] = err.Error()
			return
		}

	}
}

func checkAndSave(user string, users []string, share *Share) ([]string, error) {
	shareFilePath := filepath.Join(conf.Workspace, user, "share.json")
	data := func(shareFilePath string) (data []byte) {
		file, err := os.Open(shareFilePath)
		if err != nil {
			data = []byte{}
			logger.Errorf("share file error, %v.\n", err)
			return data
		}
		data, err = ioutil.ReadAll(file)
		defer file.Close()
		if err != nil {
			data = []byte{}
			logger.Errorf("share file error, %v.\n", err)
		}
		return data
	}(shareFilePath)
	shareList := []*Share{}
	if len(data) > 2 {
		// not empty!
		err := json.Unmarshal(data, &shareList)
		if err != nil {
			logger.Errorf("share file error, %v.\n", err)
			return nil, err
		}
	}
	// check if this share is exist!
	index := -1
	for i, oShare := range shareList {
		if share.Owner == oShare.Owner && share.DocName == oShare.DocName {
			index = i
		}
	}
	if index == -1 {
		shareList = append(shareList, share)
	} else {
		shareList[index] = share
	}
	data, err := json.Marshal(shareList)
	if err != nil {
		logger.Errorf("share file error, %v.\n", err)
		return nil, err
	}
	err = ioutil.WriteFile(shareFilePath, data, 0644)
	if err != nil {
		logger.Errorf("share file error, %v.\n", err)
		return nil, err
	}
	users = append(users, user)
	return users, nil
}

func checkAndDel(user string, share *Share) error {
	shareFilePath := filepath.Join(conf.Workspace, user, "share.json")
	data := func(shareFilePath string) (data []byte) {
		file, err := os.Open(shareFilePath)
		if err != nil {
			data = []byte{}
			logger.Errorf("share file error, %v.\n", err)
			return data
		}
		data, err = ioutil.ReadAll(file)
		defer file.Close()
		if err != nil {
			data = []byte{}
			logger.Errorf("share file error, %v.\n", err)
		}
		return data
	}(shareFilePath)
	shareList := []*Share{}
	if len(data) < 2 {
		// empty!
		return nil
	} else {
		// not empty!
		err := json.Unmarshal(data, &shareList)
		if err != nil {
			logger.Errorf("share file error, %v.\n", err)
			return err
		}
	}
	// check if this share is exist!
	index := -1
	for i, oShare := range shareList {
		if share.Owner == oShare.Owner && share.DocName == oShare.DocName {
			index = i
		}
	}
	if index == -1 {
		return nil
	} else {
		shareList = append(shareList[:index], shareList[index+1:]...)
	}
	data, err := json.Marshal(shareList)
	if err != nil {
		logger.Errorf("share file error, %v.\n", err)
		return err
	}
	err = ioutil.WriteFile(shareFilePath, data, 0644)
	if err != nil {
		logger.Errorf("share file error, %v.\n", err)
	}
	return err
}

func checkAndUpdate(user string, oldShare, newShare *Share) error {
	shareFilePath := filepath.Join(conf.Workspace, user, "share.json")
	data := func(shareFilePath string) (data []byte) {
		file, err := os.Open(shareFilePath)
		if err != nil {
			data = []byte{}
			logger.Errorf("share file error, %v.\n", err)
			return data
		}
		data, err = ioutil.ReadAll(file)
		defer file.Close()
		if err != nil {
			data = []byte{}
			logger.Errorf("share file error, %v.\n", err)
		}
		return data
	}(shareFilePath)
	shareList := []*Share{}
	if len(data) < 2 {
		// empty!
		return nil
	} else {
		// not empty!
		err := json.Unmarshal(data, &shareList)
		if err != nil {
			logger.Errorf("share file error, %v.\n", err)
			return err
		}
	}
	// check if this share is exist!
	index := -1
	for i, oShare := range shareList {
		if oldShare.Owner == oShare.Owner && oldShare.DocName == oShare.DocName {
			index = i
		}
	}
	if index == -1 {
		return nil
	} else {
		newShare.ShareType = oldShare.ShareType
		shareList[index] = newShare
	}
	data, err := json.Marshal(shareList)
	if err != nil {
		logger.Errorf("share file error, %v.\n", err)
		return err
	}
	err = ioutil.WriteFile(shareFilePath, data, 0644)
	if err != nil {
		logger.Errorf("share file error, %v.\n", err)
	}
	return err
}

func shareListHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"succ": true}
	defer util.RetJSON(w, r, data)

	httpSession, _ := httpSessionStore.Get(r, "coditor-session")
	userSession := httpSession.Values[user_session]
	if nil == userSession {
		data["succ"] = false
		data["msg"] = "permission denied"
		return
	}

	user := userSession.(*User)
	shareList, err := getOrInitShareFiles(user)
	if err != nil {
		data["succ"] = false
		data["msg"] = err.Error()
		return
	}

	data["shares"] = shareList
	logger.Infof("shares is %v.\n", shareList)
}

func getShareInfoHandler(w http.ResponseWriter, r *http.Request) {
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
		data["msg"] = "args decode error!"
		return
	}

	docName := args["docName"]
	if docName == nil || len(docName.(string)) == 0 {
		data["succ"] = false
		data["msg"] = "docName can not be null!"
		return
	}
	filePath := filepath.Join(conf.Workspace, docName.(string))
	dmd, err := newDocumentMetaData(filePath)
	if err != nil {
		logger.Error(err)
		data["succ"] = false
		data["msg"] = err.Error()
		return
	}

	data["shareInfo"] = dmd
}

func getOrInitShareFiles(u *User) ([]*Share, error) {
	shareFilePath := u.getFileBasePath("share.json")
	file, err := os.Open(shareFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// new and init share.json
			shareList := []*Share{}
			data, err := json.Marshal(shareList)
			if err != nil {
				return nil, err
			}
			err = ioutil.WriteFile(shareFilePath, data, 0644)
			if err != nil {
				return nil, err
			}
			file, err = os.Open(shareFilePath)
			if err != nil {
				return nil, err
			}
		}
	}
	data, err := ioutil.ReadAll(file)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	shareList := []*Share{}
	if len(data) > 2 {
		// empty!
		err = json.Unmarshal(data, &shareList)
		if err != nil {
			return nil, err
		}
	}
	return shareList, nil
}
