package storage

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-redis/redis/v8"
)

var (
	// redis_addr = os.Getenv("redis_addr")
	redis_addr     = "127.0.0.1:6379"
	redis_password = os.Getenv("redis_password")
)

func TestRedisSetData(t *testing.T) {
	cl := redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: redis_password,
	})
	st := NewRedisStorage(cl)
	ptt := PackType{}
	text := "LOL"
	ptt["1"] = text
	err := st.SetData(0, 0, ptt)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	pt, err := st.GetData(0, 0)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if pt["1"] != text {
		t.Error("Cant set a key: value, or cant get from redis correct value, value =", pt)
		t.Fail()
	}
}

func TestRedisResolveKey(t *testing.T) {
	rs := &RedisStorage{}
	var cid, uid int64 = 0, 0
	key := rs.resolveKey(cid, uid)
	key2 := strings.Join([]string{fmt.Sprintln(cid), fmt.Sprintln(uid)}, ":")
	if key != key2 {
		t.Error("Key is not same Key2!")
	}
}
