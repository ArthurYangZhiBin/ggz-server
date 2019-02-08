package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/golang/glog"
	"ggz-server/object"
	"ggz-server/util"
	"ggz-server/store"
	"github.com/dgraph-io/badger"
	"gitlabClient"
	"io/ioutil"
)


func CreateGitlabClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	group := vars["group"]
	token := vars["token"]

	if group == "" || token == "" {
		glog.Error("group or token 为空")
		util.WriteJsonString(w, object.NewParamErrReturnObj())
		return
	}

	tokens, err := getTokens(group)
	if err != nil {
		glog.Error(err)
		util.WriteJsonString(w, object.NewServerErrReturnObj())
		return
	}
	tokens[token] = struct{}{}

	err = saveTokens(group, tokens)
	if err != nil {
		glog.Error(err)
		util.WriteJsonString(w, object.NewServerErrReturnObj())
	}

	if _, exist := gitlabClient.GitlabClients[token]; !exist {
		gitlabClient.PushGitlabClient(token)
	}

	util.WriteJsonString(w, object.NewSuccessReturnObj())
}

func GetTokens(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	group := vars["group"]

	if group == "" {
		glog.Error("group为空")
		util.WriteJsonString(w, object.NewParamErrReturnObj())
		return
	}
	tokens, err := getTokens(group)
	if err != nil {
		glog.Error(err)
		util.WriteJsonString(w, object.NewServerErrReturnObj())
		return
	}
	util.WriteJsonString(w, object.NewSuccessWithDataReturnObj(tokens))
}

func DelToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	group := vars["group"]
	token := vars["token"]

	if group == "" || token == "" {
		glog.Error("group or token 为空")
		util.WriteJsonString(w, object.NewParamErrReturnObj())
		return
	}

	tokens, err := getTokens(group)
	if err != nil {
		glog.Error(err)
		util.WriteJsonString(w, object.NewServerErrReturnObj())
		return
	}

	delete(tokens, token)
	err = saveTokens(group, tokens)
	if err != nil {
		glog.Error(err)
		util.WriteJsonString(w, object.NewServerErrReturnObj())
		return
	}
	util.WriteJsonString(w, object.NewSuccessWithDataReturnObj(tokens))
}

func SearchProject(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Error(err)
		util.WriteJsonString(w, object.NewServerErrReturnObj())
		return
	}

	// 获取group
	cache := make(map[string]interface{}, 16)
	err = util.UnMarshal(data, cache)
	if err != nil {
		glog.Error(err)
		util.WriteJsonString(w, object.NewServerErrReturnObj())
		return
	}

	var allTokens object.Set
	for g, _ := range cache {
		// 获取token
		tokens, err := getTokens(g)
		if err != nil {
			glog.Error(err)
			util.WriteJsonString(w, object.NewServerErrReturnObj())
			return
		}
		object.PushSet(allTokens, tokens)
	}

	// 获取所有项目
	for token, _ := range allTokens {
		gitlabClient.GitlabClients[token].ListOwnProject()
	}

}

func SelectBranch(w http.ResponseWriter, r *http.Request) {

}

func getTokens(group string) (object.Set, error) {
	data, err := store.View(group + object.GitClient)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return object.Set{}, nil
		}
		return nil, err;
	}
	var tokens object.Set
	err = util.UnMarshal(data, &tokens)
	if err != nil {
		return nil, err;
	}
	return tokens, nil
}

func saveTokens(group string, tokens object.Set) error {
	data, err := util.Marshal(tokens)
	if err != nil {
		return err
	}

	err = store.Store(group + object.GitClient, data)
	if err != nil {
		return err
	}
	return nil
}