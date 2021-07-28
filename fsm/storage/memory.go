package storage

// 2 dimensional dictionary(mapping(assaciative mapping))
// First key is Chat id, second is User id
type DataType map[int64]map[int64]*StorageRecord

// Simple memory storage
type MemoryStorage struct {
	Data *DataType
}

// ResolveData ...
func (ms *MemoryStorage) ResolveData(ChatId int64, UserId int64) *StorageRecord {
	d := *ms.Data
	_, ok := d[ChatId][UserId]
	if !ok {
		d[ChatId] = map[int64]*StorageRecord{}
	}
	record, ok := d[ChatId][UserId]

	if !ok {
		d[ChatId][UserId] = &StorageRecord{}
		record = d[ChatId][UserId]
	}

	return record
}

// SetData ...
func (ms *MemoryStorage) SetData(cid int64, uid int64, data PackType) error {
	record := ms.ResolveData(cid, uid)
	record.Data = data
	return nil
}

// GetData ...
func (ms *MemoryStorage) GetData(cid int64, uid int64) (PackType, error) {
	return ms.ResolveData(cid, uid).Data, nil
}

// SetState ...
func (ms *MemoryStorage) SetState(cid int64, uid int64, state string) error {
	record := ms.ResolveData(cid, uid)
	record.State = state
	return nil
}

// GetState ...
func (ms *MemoryStorage) GetState(cid int64, uid int64) (string, error) {
	return ms.ResolveData(cid, uid).State, nil
}

// Deletes all stored data
func (ms *MemoryStorage) Close() {
	d := *ms.Data
	for key, value := range d {
		for key := range value {
			delete(value, key)
		}
		delete(d, key)
	}
	ms.Data = &DataType{}
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		Data: &DataType{},
	}
}
