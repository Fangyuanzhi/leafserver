package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

type AgentRedis struct {
	RedisConn *redis.Pool
	Addr      string
}

var (
	GRedisAgent = &AgentRedis{}
)

func GetRedisConn() *redis.Pool {
	return GRedisAgent.GetRedisConn()
}

func SetRedisConnAddr(addr string) {
	GRedisAgent.Setup(addr)
}

// Setup RedisConn redis连接池
func (r *AgentRedis) Setup(addr string) {
	r.Addr = addr
}

func (r *AgentRedis) GetRedisConn() *redis.Pool {
	if r.Addr == "" {
		return r.RedisConn
	}

	if r.RedisConn == nil {
		r.RedisConn = &redis.Pool{
			//最大空闲连接数
			MaxIdle: 30,
			//在给定时间内，允许分配的最大连接数（当为零时，没有限制）
			MaxActive: 30,
			//在给定时间内，保持空闲状态的时间，若到达时间限制则关闭连接（当为零时，没有限制）
			IdleTimeout: 200,
			//提供创建和配置应用程序连接的一个函数
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", r.Addr)
				if err != nil {
					return nil, err
				}
				//如果redis设置了用户密码，使用auth指令
				//if _, err := c.Do("AUTH", password); err != nil {
				//	c.Close()
				//	return nil, err
				//}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}
	}

	return r.RedisConn
}
