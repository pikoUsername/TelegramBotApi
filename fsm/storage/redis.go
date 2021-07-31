package storage

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/go-redis/redis/v8"
)

// RedisStorage uses as connector to redis server
type RedisStorage struct {
	client  *redis.Client
	Context context.Context
}

// NewRedisStorage ...
func NewRedisStorage(cl *redis.Client) *RedisStorage {
	return &RedisStorage{
		client:  cl,
		Context: context.Background(),
	}
}

// Close ...
func (rs *RedisStorage) Close() {
	rs.client.Close()
}

// ResolveKey ...
func (rs *RedisStorage) ResolveKey(raw_parts ...interface{}) string {
	var parts []string

	for _, r_part := range raw_parts {
		b, _ := json.Marshal(r_part)
		parts = append(parts, (string)(b))
	}

	return strings.Join(parts, ":")
}

// SetData ...
func (rs *RedisStorage) SetData(cid, uid int64, pt PackType) error {
	state, err := rs.GetState(cid, uid)
	if err != nil {
		return err
	}
	v := &StorageRecord{
		Data:  pt,
		State: state,
	}
	return rs.client.Set(rs.Context, rs.ResolveKey(cid, uid), v, 0).Err()
}

// SetState ...
func (rs *RedisStorage) SetState(cid, uid int64, state string) error {
	data, err := rs.GetData(cid, uid)
	if err != nil {
		return err
	}
	v := &StorageRecord{
		Data:  data,
		State: state,
	}
	return rs.client.Set(rs.Context, rs.ResolveKey(cid, uid), v, 0).Err()
}

// GetData ...
func (rs *RedisStorage) GetData(cid, uid int64) (PackType, error) {
	sr, err := rs.GetValue(cid, uid)
	if err != nil {
		return (PackType)(nil), err
	}
	return sr.Data, nil
}

// GetState
func (rs *RedisStorage) GetState(cid, uid int64) (string, error) {
	sr, err := rs.GetValue(cid, uid)
	if err != nil {
		return "", err
	}
	return sr.State, nil
}

// GetValue ...
func (rs *RedisStorage) GetValue(cid, uid int64) (*StorageRecord, error) {
	val := rs.client.Get(rs.Context, rs.ResolveKey(cid, uid))
	v, err := val.Result()
	if err != nil {
		return (*StorageRecord)(nil), err
	}
	var res *StorageRecord
	err = json.Unmarshal(([]byte)(v), res)
	if err != nil {
		return res, err
	}
	return res, nil
}
