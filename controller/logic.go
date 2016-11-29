package controller

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/real-film/config"
	"github.com/reechou/real-film/ext"
)

type Logic struct {
	sync.Mutex
	
	cfg  *config.Config
	film *ext.FilmHttpClient
	
	playerMap map[string]*FilmPlayerInfo
}

func NewLogic(cfg *config.Config) *Logic {
	l := &Logic{
		cfg:  cfg,
		film: ext.NewFilmHttpClient(),
		playerMap: make(map[string]*FilmPlayerInfo),
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

func (self *Logic) GetPlayerCache(key string) []byte {
	self.Lock()
	defer self.Unlock()
	
	v, ok := self.playerMap[key]
	if ok {
		if time.Now().Unix() - v.UpdateTime < int64(self.cfg.FilmPlayerExpired) {
			result, err := json.Marshal(v.Player)
			if err != nil {
				return nil
			}
			return result
		}
	}
	
	return nil
}

func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func WriteBytes(w http.ResponseWriter, code int, v []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	w.Header().Set("Content-Type", "application/json")
	
	w.WriteHeader(code)
	w.Write(v)
}

func EnableDebug() {
	logrus.SetLevel(logrus.DebugLevel)
}
