package logic

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

var URL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v"
var ShortData = "{\"action\":\"long2short\",\"long_url\":\"%v\"}"
var ShortURL = "https://api.weixin.qq.com/cgi-bin/shorturl?access_token=%v"

type ShortUrl struct {
	Errcode   int64  `json:"errorcode"`
	Errmsg    string `json:"errmsg"`
	Short_url string `json:"short_url"`
}
type Token struct {
	Access_token string `json:"access_token"`
	Expires_in   int64  `json:"expires_in"`
}

type Urlcheck struct {
	URL         string
	SortUrl     string
	Inavailable bool
}

var tokens []*TokenAccount

type TokenAccount struct {
	Expires      int64
	Access_token string
	WxAppId      string
	WxSecret     string
}

func getToken() string {
	t := tokens[rand.Intn(len(tokens))]
	if t.Access_token == "" || t.Expires-time.Now().Unix() < 60 {
		t.getToken()
	}
	return t.Access_token
}

func (self *TokenAccount) httpGet(url string) (string) {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

func (self *TokenAccount) getToken() {
	body := self.httpGet(fmt.Sprintf(URL, self.WxAppId, self.WxSecret))
	token := &Token{}
	err := json.Unmarshal([]byte(body), &token)
	if err != nil {
		glog.Error(err)
		return
	}
	self.Access_token = token.Access_token
	self.Expires = token.Expires_in + time.Now().Unix()
}

func InitAccount(list [][]string) {
	for _, v := range list {
		tokens = append(tokens, &TokenAccount{WxAppId: v[0], WxSecret: v[1]})
	}
}

type HashUrlcheck struct {
	sync.RWMutex
	hash map[string]*Urlcheck
}

func (self HashUrlcheck) get(url string) *Urlcheck {
	self.RLock()
	d := self.hash[url]
	self.RUnlock()
	return d
}

func (self HashUrlcheck) set(check *Urlcheck) {
	self.Lock()
	self.hash[check.URL] = check
	self.Unlock()
}

var hash = &HashUrlcheck{hash: make(map[string]*Urlcheck)}

func Check(url string) bool {
	cheker := hash.get(url)

	if cheker == nil {
		cheker = &Urlcheck{URL: url}
		hash.set(cheker)
	}
	return cheker.check()
}

func (self *Urlcheck) check() bool {
	if self.Inavailable {
		return false
	}
	if self.SortUrl == "" {
		self.SortUrl = self.shortUrl(self.URL)
		//glog.Error("self.SortUrl == ''1")
	}

	if (self.SortUrl == "") {
		glog.Error("self.SortUrl == ''2")
		return true
	}
	content := self.httpGet(self.SortUrl)
	if content == "" {
		glog.Error("content == ''")
		return true
	}
	//glog.Info(content)

	self.Inavailable = strings.Contains(content, "\"title\":\"已停止访问该网页\"")
	return !self.Inavailable
}
func (self *Urlcheck) shortUrl(long_url string) string {
	token := getToken()
	post_data := fmt.Sprintf(ShortData, long_url)
	url := fmt.Sprintf(ShortURL, token)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(post_data))

	if err != nil {
		glog.Error(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	short := &ShortUrl{}

	err = json.Unmarshal(body, short)
	if err != nil {
		glog.Error(err)
		return ""
	}
	if (short.Errcode != 0) {
		glog.Error(short.Errcode)
		return ""
	}
	return short.Short_url
}

func (self *Urlcheck) httpGet(url string) (string) {
	glog.Info(url)
	resp, err := http.Get(url)
	if err != nil {
		glog.Error(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return string(body)
}
