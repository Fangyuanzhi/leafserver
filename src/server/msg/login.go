package msg

type AccountInfo struct {
	Account     string
	PassWD      string
	PassWDAgain string
	PassOld     string
	Op          string
}

type UserInfo struct {
	Id         uint64
	Name       string
	Level      int32
	Experience int32
	Token      string
}

type RetAccountInfo struct {
	Msg     string
	ErrCode int32
	User    *UserInfo
}
