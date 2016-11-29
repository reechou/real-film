package controller

import (
	"net/http"
	
	"github.com/Sirupsen/logrus"
)

func (self *Logic) GetFilmPlayer(w http.ResponseWriter, r *http.Request) {
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
	
	WriteJSON(w, http.StatusOK, fp)
}
