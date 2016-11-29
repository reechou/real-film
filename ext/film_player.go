package ext

import (
	"net/http"
	"io/ioutil"
)

func (self *FilmHttpClient) GetPlayerMyunbo(p *FilmHostType, fp *FilmPlayer) error {
	u := "http://" + p.Host + p.Path
	httpReq, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	httpReq.Header.Add("User-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 5_1 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Mobile/9B176 MicroMessenger/4.3.2")
	httpReq.Header.Add("Referer", "https://myunbo.duapp.com/player.php?vid=CODUzODAxMg==~8c2ad552.acku")
	
	rsp, err := self.client.Do(httpReq)
	defer func() {
		if rsp != nil {
			rsp.Body.Close()
		}
	}()
	if err != nil {
		return err
	}
	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	fp.PlayerUrl = string(rspBody)
	
	return nil
}
