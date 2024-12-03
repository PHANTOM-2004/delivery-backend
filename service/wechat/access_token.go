package wechat_service

import (
	"delivery-backend/internal/setting"
	"encoding/json"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

var WXClient = &http.Client{
	Timeout: 5 * time.Second, // 设置超时时间
}

func Setup() {
	WXTokenHandler = newWXtoken()
	WXTokenHandler.Refresh()
	go WXTokenHandler.Serve()
}

type WXtoken struct {
	access_token string
	lock         sync.RWMutex
	ticker       time.Ticker
	signal       chan struct{}
	refreshing   atomic.Bool
}

var WXTokenHandler *WXtoken

func newWXtoken() *WXtoken {
	return &WXtoken{
		// 7200s refresh
		ticker: *time.NewTicker(time.Duration(setting.WechatSetting.TokenRefreshInterval) * time.Second),
		signal: make(chan struct{}),
	}
}

func (w *WXtoken) Serve() {
	for range w.ticker.C {
		log.Info("timer tickes[internal:7200s] refreshing access_token")
		w.Refresh()
	}
}

func (w *WXtoken) Get() string {
	w.lock.RLock()
	defer w.lock.RUnlock()
	return w.access_token
}

func (w *WXtoken) Refresh() {
	// 如果正在刷新, 那么不调用
	if w.refreshing.Load() {
		return
	}
	w.refreshing.Store(true)
	log.Info("refreshing state set to true")
	// NOTE:第一次直接锁定refresh函数，时长两秒，防止多次refresh
	defer time.AfterFunc(time.Second*2,
		func() {
			w.refreshing.Store(false)
			log.Info("refreshing state set to false")
		})

	w.lock.Lock()
	defer w.lock.Unlock()
	log.Info("refreshing token")

	// NOTE: 进行refresh 请求
	// https://api.weixin.qq.com/cgi-bin/token
	url := setting.WechatSetting.GetAccessTokenURL()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Panic(err)
	}
	resp, err := WXClient.Do(req)
	if err != nil {
		log.Panic(err)
	}
	wxserverResp := map[string]any{}
	err = json.NewDecoder(resp.Body).Decode(&wxserverResp)
	if err != nil {
		log.Panic(err)
	}
	log.Trace("wxserver getAccessToken response: ", wxserverResp)

	ok := false
	w.access_token, ok = wxserverResp["access_token"].(string)
	if !ok {
		log.Error("Failed to get access_token")
		return
	}
	log.Info("token refreshed, access_token: ", w.access_token)
}
