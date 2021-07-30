package storage_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/pikoUsername/tgp/fsm/storage"
)

var (
	// redis_addr = os.Getenv("redis_addr")
	redis_addr = "127.0.0.1:6379"
)

func assert(b bool, t testing.T) {
	if !b {
		t.Error("Assert failed")
	}
}

func TestRedisSetData(t *testing.T) {
	cl := redis.NewClient(&redis.Options{
		Addr: redis_addr,
	})
	st := storage.NewRedisStorage(cl)
	ptt := storage.PackType{}
	text := "LOL"
	ptt["1"] = text
	err := st.SetData(0, 0, ptt)
	if err != nil {
		t.Error(err)
	}
	pt, err := st.GetData(0, 0)
	if err != nil {
		t.Error(err)
	}
	if pt["1"] != text {
		t.Error("Cant set a key: value, or cant get from redis correct value, value =", pt)
	}
}

func TestRedisResolveKey(t *testing.T) {
	rs := &storage.RedisStorage{}
	cid, uid := 0, 0
	key := rs.ResolveKey(cid, uid)
	key2 := strings.Join([]string{fmt.Sprintln(cid), fmt.Sprintln(uid)}, ":")
	if key != key2 {
		t.Error("Key is not same Key2!")
	}
}
