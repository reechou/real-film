package controller

import (
	"net/http"
	"encoding/json"
	
	"github.com/reechou/real-film/config"
	"github.com/reechou/real-film/ext"
	"github.com/Sirupsen/logrus"
)

type Logic struct {
	cfg *config.Config
	film *ext.FilmHttpClient
}

func NewLogic(cfg *config.Config) *Logic {
	l := &Logic{
		cfg: cfg,
		film: ext.NewFilmHttpClient(),
	}
	l.init()
	
	return l
}

func (self *Logic) init() {
	http.HandleFunc("/diediao/getPlayer", self.GetFilmPlayer)
}

func (self *Logic) Run() {
	if self.cfg.Debug {
		EnableDebug()
	}
	
	logrus.Infof("film server starting on[%s]..", self.cfg.Host)
	logrus.Infoln(http.ListenAndServe(self.cfg.Host, nil))
}

func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func EnableDebug() {
	logrus.SetLevel(logrus.DebugLevel)
}
