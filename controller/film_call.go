package controller

import (
	"net/http"
	"time"
	
	"github.com/Sirupsen/logrus"
)

func (self *Logic) GetFilmPlayer(w http.ResponseWriter, r *http.Request) {
	urlKey := r.URL.String()
	start := time.Now()
	defer func() {
		logrus.Debugf("[film] http: request url[%s] use_time[%v]", urlKey, time.Now().Sub(start))
	}()
	
	result := self.GetPlayerCache(urlKey)
	if result != nil {
		WriteBytes(w, http.StatusOK, result)
		return
	}
	
	r.ParseForm()
	
	if len(r.Form["playname"]) == 0 || len(r.Form["vid"]) == 0 {
		logrus.Errorf("req has no playname or vid")
		return
	}
	
	logrus.Debugf("get player req: %s %s", r.Form["playname"][0], r.Form["vid"][0])
	fp, err := self.film.GetPlayer(r.Form["playname"][0], r.Form["vid"][0])
	if err != nil {
		logrus.Errorf("get film player error: %v", err)
	}
	
	self.Lock()
	fv, ok := self.playerMap[urlKey]
	if ok {
		fv.UpdateTime = time.Now().Unix()
		fv.Player = fp
	} else {
		self.playerMap[urlKey] = &FilmPlayerInfo{
			Player: fp,
			UpdateTime: time.Now().Unix(),
		}
	}
	self.Unlock()
	
	WriteJSON(w, http.StatusOK, fp)
}
