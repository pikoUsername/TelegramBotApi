package storage

type PackType map[string]interface{}

// Simple storage interface for saving data,
// and uses for save FSM data
type Storage interface {
	SetData(cid int64, uid int64, data *PackType)
	GetData(cid int64, uid int64) *PackType
	SetState(cid int64, uid int64, state string)
	GetState(cid int64, uid int64) string
	Clear()
}
