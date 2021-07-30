package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
)

// RedisStorage uses as connector to redis server
type RedisStorage struct {
	client *redis.Client
}

// NewRedisStorage ...
func NewRedisStorage(cl *redis.Client) *RedisStorage {
	return &RedisStorage{
		client: cl,
	}
}

func (rs *RedisStorage) Close() {
	rs.client.Close()
}

func (rs *RedisStorage) ResolveKey(parts ...interface{}) string {
	map_f := func(f func(interface{}) string, val []interface{}) []string {
		s := []string{}
		for _, v := range val {
			b := f(v)
			s = append(s, b)
		}
		return s
	}
	part := map_f(func(i interface{}) string { return fmt.Sprintln(i) }, parts)
	s := strings.Join(part, ":")
	return s
}

func (rs *RedisStorage) SetData(cid, uid int64, pt PackType) error {
	state, err := rs.GetState(cid, uid)
	if err != nil {
		return err
	}
	v := &StorageRecord{
		Data:  pt,
		State: state,
	}
	return rs.client.Set(context.Background(), rs.ResolveKey(cid, uid), v, 0).Err()
}

func (rs *RedisStorage) SetState(cid, uid int64, state string) error {
	data, err := rs.GetData(cid, uid)
	if err != nil {
		return err
	}
	v := &StorageRecord{
		Data:  data,
		State: state,
	}
	return rs.client.Set(context.Background(), rs.ResolveKey(cid, uid), v, 0).Err()
}

func (rs *RedisStorage) GetData(cid, uid int64) (PackType, error) {
	sr, err := rs.GetValue(cid, uid)
	if err != nil {
		return PackType{}, err
	}
	return sr.Data, nil
}

func (rs *RedisStorage) GetState(cid, uid int64) (string, error) {
	sr, err := rs.GetValue(cid, uid)
	if err != nil {
		return "", err
	}
	return sr.State, nil
}

func (rs *RedisStorage) GetValue(cid, uid int64) (*StorageRecord, error) {
	val := rs.client.Get(context.Background(), rs.ResolveKey(cid, uid))
	v, err := val.Bytes()
	if err != nil {
		return nil, err
	}
	err = val.Err()
	if err != nil {
		return nil, err
	}
	var res *StorageRecord
	err = json.Unmarshal(v, res)
	if err != nil {
		return res, err
	}
	return res, nil
}
