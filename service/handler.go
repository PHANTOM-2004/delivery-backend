package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminInfoHandler struct {
	session sessions.Session
}

func NewAdminInfoHanlder(c *gin.Context) *AdminInfoHandler {
	s := sessions.Default(c)
	return &AdminInfoHandler{session: s}
}

func (a *AdminInfoHandler) SetAccount(account string) {
	a.session.Set("admin_account", account)
}

func (a *AdminInfoHandler) SetID(id uint) {
	a.session.Set("admin_id", id)
}

func (a *AdminInfoHandler) GetAccount() string {
	account := a.session.Get("admin_account")
	if account == nil {
		return ""
	}
	return account.(string)
}

func (a *AdminInfoHandler) GetID() uint {
	id := a.session.Get("admin_id")
	if id == nil {
		return 0
	}
	return id.(uint)
}

// 删除这个session, 注意必须有key才能够成功工作
func (a *AdminInfoHandler) Delete() error {
	// 注意必须含有key的情况下才能够成功clear
	a.session.Clear()
	a.session.Options(sessions.Options{MaxAge: -1})
	err := a.session.Save()
	return err
}

func (a *AdminInfoHandler) Save() error {
	return a.session.Save()
}

type MerchInfoHanlder struct {
	session sessions.Session
}

func NewMerchInfoHanlder(c *gin.Context) *MerchInfoHanlder {
	s := sessions.Default(c)
	return &MerchInfoHanlder{session: s}
}

func (a *MerchInfoHanlder) SetAccount(account string) {
	a.session.Set("merch_account", account)
}

func (a *MerchInfoHanlder) SetID(id uint) {
	a.session.Set("merch_id", id)
}

func (a *MerchInfoHanlder) GetAccount() string {
	account := a.session.Get("merch_account")
	if account == nil {
		return ""
	}
	return account.(string)
}

func (a *MerchInfoHanlder) GetID() uint {
	id := a.session.Get("merch_id")
	if id == nil {
		return 0
	}
	return id.(uint)
}

// 删除这个session, 注意必须有key才能够成功工作
func (a *MerchInfoHanlder) Delete() error {
	// 注意必须含有key的情况下才能够成功clear
	a.session.Clear()
	a.session.Options(sessions.Options{MaxAge: -1})
	err := a.session.Save()
	return err
}

func (a *MerchInfoHanlder) Save() error {
	return a.session.Save()
}
