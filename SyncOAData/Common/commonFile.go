package Common

import (
	"net/http"
)

type (
	Syncer struct {
		Transport *http.Transport
		Client    *http.Client
	}
	ResBody2 struct {
		Code int      `json:"code"`
		Msg  string   `json:"msg"`
		Data []Config `json:"data"`
	}
	Config struct {
		Key      string     `json:"key"`
		Env      string     `json:"env"`
		AppName  string     `json:"appName"`
		Version  string     `json:"version"`
		Config   string     `json:"config"`
		Remark   string     `json:"remark"`
		Versions []Versions `json:"versions"`
	}
	Versions struct {
		Version string `json:"version"`
		Config  string `json:"config"`
		Remark  string `json:"remark"`
	}
)
