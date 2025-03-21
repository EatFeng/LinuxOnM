package dto

import "time"

type HostConnTest struct {
	Addr       string `json:"addr" validate:"required"`
	Port       uint   `json:"port" validate:"required,number,max=65535,min=1"`
	User       string `json:"user" validate:"required"`
	AuthMode   string `json:"authMode" validate:"oneof=password key"`
	Password   string `json:"password"`
	PrivateKey string `json:"privateKey"`
	PassPhrase string `json:"passPhrase"`
}

type HostOperate struct {
	ID               uint   `json:"id"`
	GroupID          uint   `json:"group_id"`
	Name             string `json:"name"`
	Addr             string `json:"addr" validate:"required"`
	Port             uint   `json:"port" validate:"required,number,max=65535,min=1"`
	User             string `json:"user" validate:"required"`
	AuthMode         string `json:"authMode" validate:"oneof=password key"`
	Password         string `json:"password"`
	PrivateKey       string `json:"privateKey"`
	PassPhrase       string `json:"passPhrase"`
	RememberPassword bool   `json:"rememberPassword"`

	Description string `json:"description"`
}

type HostInfo struct {
	ID               uint      `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	GroupID          uint      `json:"group_id"`
	GroupBelong      string    `json:"groupBelong"`
	Name             string    `json:"name"`
	Addr             string    `json:"addr"`
	Port             uint      `json:"port"`
	User             string    `json:"user"`
	AuthMode         string    `json:"authMode"`
	Password         string    `json:"password"`
	PrivateKey       string    `json:"privateKey"`
	PassPhrase       string    `json:"passPhrase"`
	RememberPassword bool      `json:"rememberPassword"`

	Description string `json:"description"`
}

type SearchForTree struct {
	Info string `json:"info"`
}

type HostTree struct {
	ID       uint        `json:"id"`
	Label    string      `json:"label"`
	Children []TreeChild `json:"children"`
}

type TreeChild struct {
	ID    uint   `json:"id"`
	Label string `json:"label"`
}

type SearchHostWithPage struct {
	PageInfo
	GroupID uint   `json:"group_id"`
	Info    string `json:"info"`
}

type ChangeHostGroup struct {
	ID      uint `json:"id" validate:"required"`
	GroupID uint `json:"group_id" validate:"required"`
}
