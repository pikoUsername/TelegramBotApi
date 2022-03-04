package storage

type PackType map[string]interface{}

// Simple storage interface for saving data,
// and uses for save FSM data
type Storage interface {
	SetData(cid, uid int64, data PackType) error
	GetData(cid, uid int64) (PackType, error)
	SetState(cid, uid int64, state string) error
	GetState(cid, uid int64) (string, error)
	Clear(cid, uid int64) error
	Close()
}

// StorageRecord uses for input, and output value type
type StorageRecord struct {
	Data  PackType
	State string
}

var (
	EmptyRecord = StorageRecord{}
)
