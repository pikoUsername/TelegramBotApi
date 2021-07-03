package storage

// 2 dimensional dictionary(mapping(assaciative mapping))
// First key is Chat id, second is User id
type DataType map[int64]map[int64]*MemoryStorageRecord

type MemoryStorageRecord struct {
	Data  *PackType
	State string
}

// Simple memory storage
type MemoryStorage struct {
	Data *DataType
}

// ResolveData ...
func (ms *MemoryStorage) ResolveData(ChatId int64, UserId int64) *MemoryStorageRecord {
	d := *ms.Data
	_, ok := d[ChatId][UserId]
	if !ok {
		d[ChatId] = map[int64]*MemoryStorageRecord{}
	}
	record, ok := d[ChatId][UserId]

	if !ok {
		d[ChatId][UserId] = &MemoryStorageRecord{}
		record = d[ChatId][UserId]
	}

	return record
}

// SetData ...
func (ms *MemoryStorage) SetData(cid int64, uid int64, data *PackType) {
	record := ms.ResolveData(cid, uid)
	record.Data = data
}

// GetData ...
func (ms *MemoryStorage) GetData(cid int64, uid int64) *PackType {
	return ms.ResolveData(cid, uid).Data
}

// SetState ...
func (ms *MemoryStorage) SetState(cid int64, uid int64, state string) {
	record := ms.ResolveData(cid, uid)
	record.State = state
}

// GetState ...
func (ms *MemoryStorage) GetState(cid int64, uid int64) string {
	return ms.ResolveData(cid, uid).State
}

func (ms *MemoryStorage) Clear() {
	ms.Data = &DataType{}
}

var (
	EmptyRecord = MemoryStorageRecord{}
)

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		Data: &DataType{},
	}
}
