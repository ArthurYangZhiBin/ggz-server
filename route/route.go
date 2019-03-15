package route

import (
	"ggz-server/handler"
	"ggz-server/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

var R *mux.Router

func init() {
	R = mux.NewRouter()

	c := cors.AllowAll()

	R.Handle("/config/global", negroni.New(c, middlewares.ParseFormMiddlerware, negroni.WrapFunc(handler.CreateGitlab))).Methods("POST", "OPTIONS")
	R.Handle("/config/global", negroni.New(c, negroni.WrapFunc(handler.GetGitlab))).Methods("GET", "OPTIONS")
	R.Handle("/config/project/setting/{group}/{token}", negroni.New(c, middlewares.ParseFormMiddlerware, negroni.WrapFunc(handler.CreateGitlabClient))).Methods("POST", "OPTIONS")
	R.Handle("/config/project/setting/{group}", negroni.New(c, middlewares.ParseFormMiddlerware, negroni.WrapFunc(handler.GetTokens))).Methods("GET", "OPTIONS")
	R.Handle("/config/project/setting/{group}/{token}", negroni.New(c, middlewares.ParseFormMiddlerware, negroni.WrapFunc(handler.DelToken))).Methods("DELETE", "OPTIONS")
	R.Handle("/build/projects", negroni.New(c, middlewares.ParseFormMiddlerware, negroni.WrapFunc(handler.SearchProject))).Methods("POST", "OPTIONS")
	R.Handle("/build/project/{id}/{token}", negroni.New(c, middlewares.ParseFormMiddlerware, negroni.WrapFunc(handler.SelectBranch))).Methods("POST", "OPTIONS")

	R.Handle("/beat", negroni.New(c, middlewares.ParseFormMiddlerware, negroni.WrapFunc(handler.Beat))).Methods("GET", "OPTIONS")

	R.Handle("/audit/add", negroni.New(c, middlewares.ParseFormMiddlerware, negroni.WrapFunc(handler.Add))).Methods("POST", "OPTIONS")
	R.Handle("/audit/list/{nodeId}/{username}", negroni.New(c, middlewares.ParseFormMiddlerware, negroni.WrapFunc(handler.Find))).Methods("GET", "OPTIONS")
	R.Handle("/audit/remove/{nodeId}/{aprovUserId}", negroni.New(c, middlewares.ParseFormMiddlerware, negroni.WrapFunc(handler.Remove))).Methods("GET", "OPTIONS")
	R.Handle("/audit/update", negroni.New(c, middlewares.ParseFormMiddlerware, negroni.WrapFunc(handler.Update))).Methods("POST", "OPTIONS")
	R.Handle("/audit/detail/{nodeId}/{aprovUserId}", negroni.New(c, middlewares.ParseFormMiddlerware, negroni.WrapFunc(handler.Detail))).Methods("GET", "OPTIONS")
	R.Handle("/audit/list/{nodeId}", negroni.New(c, middlewares.ParseFormMiddlerware, negroni.WrapFunc(handler.Find))).Methods("GET", "OPTIONS")

	R.Handle("/operating/result", negroni.New(c, middlewares.ParseFormMiddlerware, negroni.WrapFunc(handler.OperatingInfoFind))).Methods("GET", "OPTIONS")
	R.Handle("/operating/result/{flowID}", negroni.New(c, middlewares.ParseFormMiddlerware, negroni.WrapFunc(handler.OperatingInfoFind))).Methods("GET", "OPTIONS")

	R.Handle("/userTask/find/{type}/{aprovUserId}", negroni.New(c, middlewares.ParseFormMiddlerware, negroni.WrapFunc(handler.UserTaskFindFin))).Methods("GET", "OPTIONS")
}
