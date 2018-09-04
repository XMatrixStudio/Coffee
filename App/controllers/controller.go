package controllers

// CommonRes 返回值
type CommonRes struct {
	State string
	Data  string
}

const (
	StatusSuccess  = "success"
	StatusBadReq   = "bad_req"
	StatusNotLogin = "not_login"
	StatusNotAllow = "not_allow"
	StatusExist    = "had_exist"
)
