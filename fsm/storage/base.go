package storage

type PackType map[string]interface{}

func (pt PackType) Pop(key string) interface{} {
	v, ok := pt[key]
	if ok {
		delete(pt, key)
		return v
	}
	return nil
}

// Simple storage interface for saving data,
// and uses for save FSM data
type Storage interface {
	SetData(cid int64, uid int64, data *PackType)
	GetData(cid int64, uid int64) *PackType
	SetState(cid int64, uid int64, state string)
	GetState(cid int64, uid int64) string
	Clear()
}
