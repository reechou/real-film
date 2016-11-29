package ext

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	
	"github.com/Sirupsen/logrus"
)

type FilmHttpClient struct {
	client *http.Client
}

func NewFilmHttpClient() *FilmHttpClient {
	return &FilmHttpClient{
		client: &http.Client{},
	}
}

func (self *FilmHttpClient) GetPlayer(playName, vid string) (*FilmPlayer, error) {
	p, err := self.GetHostType(playName, vid)
	if err != nil {
		logrus.Errorf("get host type error: %v", err)
		return nil, err
	}
	fp := &FilmPlayer{Status: 1, PlayName: playName}
	switch p.Host {
	case "myunbo.duapp.com":
		err = self.GetPlayerMyunbo(p, fp)
	case "v.pptvyun.com", "iqiyi_vip", "le_vip":
		fp.PlayName = playName
		fp.PlayerUrl = p.Path
	}
	
	return fp, err
}

func (self *FilmHttpClient) GetHostType(playName, vid string) (*FilmHostType, error) {
	p := &FilmHostType{
		Path: "/dy2.php?vid=" + vid + ".haikqq",
		Host: "myun.kb20.cc",
		Port: 80,
	}
	
	var err error
	switch playName {
	case "acfun", "letv", "bibibi", "youku", "qq", "hunantv", "acku":
		err = self.GetMyunbo(vid, p)
	case "lev":
		err = self.GetLeVip(vid, p)
	case "ppyun":
		p.Host = "v.pptvyun.com"
		p.Path = vid
	case "qiyi":
		err = self.GetQiyi(vid, p)
	default:
		err = self.GetYundied(vid, p)
	}
	
	return p, err
}

func (self *FilmHttpClient) GetMyunbo(vid string, p *FilmHostType) error {
	p.Host = "myunbo.duapp.com"
	
	u := fmt.Sprintf("https://myunbo.duapp.com/player.php?vid=%s", vid)
	httpReq, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	httpReq.Header.Add("User-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 5_1 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Mobile/9B176 MicroMessenger/4.3.2")
	httpReq.Header.Add("Referer", "http://m.kb20.cc/Public/player2.9/youku.html?vid=XMTUwOTcxMDc3Ng==")
	
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
	
	reg := regexp.MustCompile(`var(.*?)urlplay(.*?)=(.*?)\'(.*?)\'`)
	urlPlayVar := reg.FindString(string(rspBody))
	if urlPlayVar == "" {
		logrus.Errorf("cannot found urlplayer string")
		return fmt.Errorf("cannot found urlplay")
	}
	urlPlayVar = strings.Replace(urlPlayVar, "'", "", -1)
	varSlice := strings.Split(urlPlayVar, " ")
	if len(varSlice) < 4 {
		logrus.Errorf("urlplay len(varSlice) < 4")
		return fmt.Errorf("cannot found urlplay")
	}
	p.UrlPlay = varSlice[3]
	
	reg = regexp.MustCompile(`var(.*?)tm(.*?)=(.*?)\'(.*?)\'`)
	tmVar := reg.FindString(string(rspBody))
	if tmVar == "" {
		logrus.Errorf("cannot found tm string")
		return fmt.Errorf("cannot found tm")
	}
	tmVar = strings.Replace(tmVar, "'", "", -1)
	varSlice = strings.Split(tmVar, " ")
	if len(varSlice) < 4 {
		logrus.Errorf("tm len(varSlice) < 4")
		return fmt.Errorf("cannot found tm")
	}
	p.Tm = varSlice[3]
	
	reg = regexp.MustCompile(`var(.*?)sign(.*?)=(.*?)\'(.*?)\'`)
	signVar := reg.FindString(string(rspBody))
	if signVar == "" {
		logrus.Errorf("cannot found sign string")
		return fmt.Errorf("cannot found sign")
	}
	signVar = strings.Replace(signVar, "'", "", -1)
	varSlice = strings.Split(signVar, " ")
	if len(varSlice) < 4 {
		logrus.Errorf("sign len(varSlice) < 4")
		return fmt.Errorf("cannot found sign")
	}
	p.Sign = varSlice[3]
	
	reg = regexp.MustCompile(`var(.*?)refer(.*?)=(.*?)\'(.*?)\'`)
	referVar := reg.FindString(string(rspBody))
	if referVar == "" {
		logrus.Errorf("cannot found refer string")
		return fmt.Errorf("cannot found refer")
	}
	referVar = strings.Replace(referVar, "'", "", -1)
	varSlice = strings.Split(referVar, " ")
	if len(varSlice) < 4 {
		logrus.Errorf("refer len(varSlice) < 4")
		return fmt.Errorf("cannot found refer")
	}
	p.UserLink = varSlice[3]
	
	p.Path = "/parse.php?h5url=" + p.UrlPlay + "&tm=" + p.Tm + "&sign=" + p.Sign + "&ajax=1&userlink=" + p.UserLink;
	
	return nil
}

func (self *FilmHttpClient) GetLeVip(vid string, p *FilmHostType) error {
	p.Host = "le_vip"
	p.Port = 6677
	
	u := fmt.Sprintf("http://vipm.kb20.cc:6677/playm3u8/?type=le_vip&vid=%s", vid)
	httpReq, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	httpReq.Header.Add("User-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 5_1 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Mobile/9B176 MicroMessenger/4.3.2")
	httpReq.Header.Add("Referer", "http://m.kb20.cc/Public/player2.9/youku.html?vid=XMTUwOTcxMDc3Ng==")
	
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
	
	reg := regexp.MustCompile(`var(.*?)play_url(.*?)=(.*?)\"(.*?)\"\;`)
	urlPlayVar := reg.FindString(string(rspBody))
	if urlPlayVar == "" {
		logrus.Errorf("cannot found urlplayer string")
		return fmt.Errorf("cannot found urlplay")
	}
	urlPlayVar = strings.Replace(urlPlayVar, "'", "", -1)
	varSlice := strings.Split(urlPlayVar, " ")
	if len(varSlice) < 4 {
		logrus.Errorf("urlplay len(varSlice) < 4")
		return fmt.Errorf("cannot found urlplay")
	}
	p.Path = varSlice[3]
	
	return nil
}

func (self *FilmHttpClient) GetQiyi(vid string, p *FilmHostType) error {
	vidSlice := strings.Split(vid, ",")
	if len(vidSlice) < 2 {
		logrus.Errorf("qiyi len(vidSlice) < 2")
		return fmt.Errorf("qiyi len(vidSlice) < 2")
	}
	
	p.Host = "iqiyi_vip"
	p.Port = 6677
	
	u := fmt.Sprintf("http://vipm.kb20.cc:6677/playm3u8/?type=iqiyi_vip&vid=%s", vidSlice[1])
	httpReq, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	httpReq.Header.Add("User-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 5_1 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Mobile/9B176 MicroMessenger/4.3.2")
	httpReq.Header.Add("Referer", "http://m.kb20.cc/Public/player2.9/youku.html?vid=XMTUwOTcxMDc3Ng==")
	
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
	
	reg := regexp.MustCompile(`var(.*?)play_url(.*?)=(.*?)\"(.*?)\"\;`)
	urlPlayVar := reg.FindString(string(rspBody))
	if urlPlayVar == "" {
		logrus.Errorf("cannot found urlplayer string")
		return fmt.Errorf("cannot found urlplay")
	}
	urlPlayVar = strings.Replace(urlPlayVar, "'", "", -1)
	varSlice := strings.Split(urlPlayVar, " ")
	if len(varSlice) < 4 {
		logrus.Errorf("urlplay len(varSlice) < 4")
		return fmt.Errorf("cannot found urlplay")
	}
	p.Path = varSlice[3]
	
	return nil
}

func (self *FilmHttpClient) GetYundied(vid string, p *FilmHostType) error {
	p.Host = "yundied.duapp.com"
	
	u := fmt.Sprintf("https://yundied.duapp.com/player.php?vid=%s", vid)
	httpReq, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	httpReq.Header.Add("User-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 5_1 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Mobile/9B176 MicroMessenger/4.3.2")
	httpReq.Header.Add("Referer", "http://m.kb20.cc/Public/player2.9/youku.html?vid=XMTUwOTcxMDc3Ng==")
	
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
	
	reg := regexp.MustCompile(`var(.*?)urlplay(.*?)=(.*?)\'(.*?)\'`)
	urlPlayVar := reg.FindString(string(rspBody))
	if urlPlayVar == "" {
		logrus.Errorf("cannot found urlplayer string")
		return fmt.Errorf("cannot found urlplay")
	}
	urlPlayVar = strings.Replace(urlPlayVar, "'", "", -1)
	varSlice := strings.Split(urlPlayVar, " ")
	if len(varSlice) < 4 {
		logrus.Errorf("urlplay len(varSlice) < 4")
		return fmt.Errorf("cannot found urlplay")
	}
	p.UrlPlay = varSlice[3]
	
	reg = regexp.MustCompile(`var(.*?)tm(.*?)=(.*?)\'(.*?)\'`)
	tmVar := reg.FindString(string(rspBody))
	if tmVar == "" {
		logrus.Errorf("cannot found tm string")
		return fmt.Errorf("cannot found tm")
	}
	tmVar = strings.Replace(tmVar, "'", "", -1)
	varSlice = strings.Split(tmVar, " ")
	if len(varSlice) < 4 {
		logrus.Errorf("tm len(varSlice) < 4")
		return fmt.Errorf("cannot found tm")
	}
	p.Tm = varSlice[3]
	
	reg = regexp.MustCompile(`var(.*?)sign(.*?)=(.*?)\'(.*?)\'`)
	signVar := reg.FindString(string(rspBody))
	if signVar == "" {
		logrus.Errorf("cannot found sign string")
		return fmt.Errorf("cannot found sign")
	}
	signVar = strings.Replace(signVar, "'", "", -1)
	varSlice = strings.Split(signVar, " ")
	if len(varSlice) < 4 {
		logrus.Errorf("sign len(varSlice) < 4")
		return fmt.Errorf("cannot found sign")
	}
	p.Sign = varSlice[3]
	
	reg = regexp.MustCompile(`var(.*?)refer(.*?)=(.*?)\'(.*?)\'`)
	referVar := reg.FindString(string(rspBody))
	if referVar == "" {
		logrus.Errorf("cannot found refer string")
		return fmt.Errorf("cannot found refer")
	}
	referVar = strings.Replace(referVar, "'", "", -1)
	varSlice = strings.Split(referVar, " ")
	if len(varSlice) < 4 {
		logrus.Errorf("refer len(varSlice) < 4")
		return fmt.Errorf("cannot found refer")
	}
	p.UserLink = varSlice[3]
	
	p.Path = "parse.php?h5url=" + p.UrlPlay + "&tm=" + p.Tm + "&sign=" + p.Sign + "&ajax=1&userlink=" + p.UserLink;
	
	return nil
}