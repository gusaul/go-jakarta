package resource

import (
	"context"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

var RedisConn RedisCache

func InitRedisConn() {
	RedisConn = GetRedisConn()
}

//RedisCache implementation of redis
type RedisCache struct {
	pool *redis.Pool
}

func GetRedisConn() RedisCache {
	return RedisCache{
		pool: &redis.Pool{
			MaxIdle:         10,
			MaxActive:       15,
			IdleTimeout:     3 * time.Second,
			MaxConnLifetime: 10 * time.Second,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", "localhost:6379", redis.DialConnectTimeout(time.Second))
				if err != nil {
					log.Println("error getting connection:", err)
					return nil, err
				}
				return conn, err
			},
			Wait: true,
		},
	}
}

// getConn pool based on settings
func (r *RedisCache) getConn() (redis.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5000)
	defer cancel()
	return r.pool.GetContext(ctx)
}

type keyPair struct {
	key    string
	fields []string
}

func (r *RedisCache) MultiHashGetPipeline(keys map[string][]string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)

	conn, err := r.getConn()
	if err != nil {
		return result, err
	}
	defer conn.Close()

	pairs := make([]keyPair, len(keys))
	idx := 0
	for key, fields := range keys {
		cmd := make([]interface{}, len(fields)+1)
		cmd[0] = key
		for i := 1; i < len(cmd); i++ {
			cmd[i] = fields[i-1]
		}
		if err := conn.Send("HMGET", cmd...); err != nil {
			return result, err
		}
		pairs[idx] = keyPair{key, fields}
		idx++
	}
	if err := conn.Flush(); err != nil {
		return result, err
	}

	for _, v := range pairs {
		data, err := redis.Values(conn.Receive())
		if err != nil {
			return result, err
		}
		result[v.key] = make(map[string]string)
		for i, val := range data {
			if bytesVal, ok := val.([]byte); ok && len(v.fields) > i {
				result[v.key][v.fields[i]] = string(bytesVal)
			}
		}
	}

	return result, err
}

func (r *RedisCache) MultiHashSetPipeline(datas map[string][]string) (errs []error) {
	conn, err := r.getConn()
	if err != nil {
		return []error{err}
	}
	defer conn.Close()

	for key, data := range datas {
		cmd := make([]interface{}, len(data)+1)
		cmd[0] = key
		for i := 1; i < len(cmd); i++ {
			cmd[i] = data[i-1]
		}
		if err := conn.Send("HMSET", cmd...); err != nil {
			errs = append(errs, err)
		}
	}
	if err := conn.Flush(); err != nil {
		errs = append(errs, err)
	}

	for range datas {
		_, err := conn.Receive()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return
}
