package internal

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"reflect"
	com "server/common"
	"server/msg"
	redisAgent "server/redis"
)

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	RedisConn = redisAgent.GetRedisConn()
	handleMsg(&msg.AccountInfo{}, dealAccount)
}

func dealAccount(args []interface{}) {
	// 收到的 Hello 消息
	m := args[0].(*msg.AccountInfo)
	// 消息的发送者
	a := args[1].(gate.Agent)

	// 输出收到的消息的内容
	log.Debug("hello %v\\nAccount:%v"+
		"PassWD:%v\\nPassWDAgain:%v\\nPassOld:%v",
		m.Op, m.Account, m.PassWD, m.PassWDAgain, m.PassOld, m.PassOld)

	var data interface{}

	switch m.Op {
	case "login":
		userinfo, ok := handleLogin(m)
		if ok == 0 {
			data = userinfo
		}
	case "register":
		userinfo, ok := handleCreateAccount(m)
		if ok == 0 {
			data = userinfo
		}
	//case "changePassword":
	//	userinfo, ok := handleLogin(m)
	//	if ok == 0 {
	//		data = userinfo
	//	}
	default:
		log.Debug("unknown op: %v", m.Op)
	}
	// 给发送者回应
	a.WriteMsg(&data)
}

// 登录
func handleLogin(request *msg.AccountInfo) (*msg.RetAccountInfo, int32) {
	var (
		name    = request.Account
		passwd  = request.PassWD
		retData = &msg.RetAccountInfo{
			Msg:     name,
			ErrCode: 0,
		}
		user      = &msg.UserInfo{}
		retInt    int
		retInts   []int
		retString string
		token     string
		id        uint64
	)
	_, code := check(name, passwd, "")
	if !code {
		return retData, 0
	}

	c, err := RedisConn.Dial()
	if err != nil {
		log.Error(err.Error())
		retData.Msg = err.Error()
		return retData, 0
	}

	for {
		retInt, err = redis.Int(c.Do("HEXISTS", name, "passwd"))
		if err != nil || retInt < 1 {
			retData.Msg = "redis 用户不存在"
			break
		}

		retString, err = redis.String(c.Do("HGET", name, "passwd"))
		if err != nil || retString == "" {
			retData.Msg = "redis 用户密码获取失败" + err.Error()
			break
		}
		if retString != passwd {
			retData.Msg = "用户密码错误：" + name
			break
		}
		// 设置状态
		_, err = redis.Int(c.Do("HSET", name, "status", 1))
		if err != nil {
			retData.Msg = "用户状态错误"
			break
		}

		// 设置token
		token = com.GetToken(16)
		_, err = redis.Int(c.Do("HSET", name, "token", token))
		if err != nil {
			retData.Msg = "用户设置token失败"
			break
		}
		id, err = redis.Uint64(c.Do("HGET", name, "id"))
		if err != nil {
			retData.Msg = "用户获取唯一用户uid失败"
			break
		}
		// 获取等级，经验
		retInts, err = redis.Ints(c.Do("HMGET", id, "level", "experience"))
		if err != nil {
			retData.Msg = "用户获取等级，经验失败"
			break
		}

		user.Id = id
		user.Level = int32(retInts[0])
		user.Experience = int32(retInts[1])
	}

	if err != nil {
		return retData, 1
	}
	retData.User = user
	return retData, 0
}

// 创建账号处理函数
func handleCreateAccount(request *msg.AccountInfo) (*msg.RetAccountInfo, int32) {

	var (
		name        = request.Account
		passwd      = request.PassWD
		passwdAgain = request.PassWDAgain
		retData     = &msg.RetAccountInfo{
			Msg:     name,
			ErrCode: 0,
		}
		user   = &msg.UserInfo{}
		retInt int
		token  string
		id     uint64
	)
	retData.User = user

	_, code := check(name, passwd, passwdAgain)
	if !code {
		return retData, 0
	}

	c, err := RedisConn.Dial()
	if err != nil {
		log.Error(err.Error())
		retData.Msg = err.Error()
		return retData, 0
	}

	for {
		retInt, err = redis.Int(c.Do("HEXISTS", name, "passwd"))
		if err != nil {
			break
		}
		if retInt > 0 {
			retData.Msg = "用户已经存在：" + name
			break
		}
		id, err = redis.Uint64(c.Do("INCR", "userid"))
		if err != nil {
			break
		}
		user.Id = id
		token = com.GetToken(16)
		_, err = redis.String(c.Do("HMSET", name, "passwd", passwdAgain, "status", 1, "token", token, "id", id))
		_, err = redis.String(c.Do("HMSET", id, "name", name, "level", 1, "experience", 0))
		if err != nil {
			break
		}
		user.Token = token
	}
	if err != nil {
		retData.Msg = err.Error()
	} else {
		_, err = redis.Int(c.Do("SADD", "playerList", id))
	}
	if err != nil {
		retData.Msg = err.Error()
		return retData, 1
	}

	return retData, 0
}

// 检查用户名、密码是否合理
func check(name, passwd1, passwd2 string) (string, bool) {
	var (
		ret  string
		code = true
	)
	if len(name) > 14 || len(name) < 3 {
		ret = fmt.Sprintf("用户名长度必须在3-14之间,您的输入的用户名为:%v", name)
		log.Error(ret)
		code = false
	} else if len(passwd1) > 24 || len(passwd1) < 6 {
		ret = "密码长度必须在6-24之间"
		log.Error(ret)
		code = false
	} else if len(passwd2) > 0 && passwd1 != passwd2 {
		ret = "两次输入的密码不一致"
		log.Error(ret)
		code = false
	}

	return ret, code
}
