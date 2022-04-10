package storage

// 2 dimensional dictionary(mapping(assaciative mapping))
// First key is Chat id, second is User id
type DataType map[int64]map[int64]*StorageRecord

// Simple memory storage
type MemoryStorage struct {
	Data DataType
}

// ResolveData ...
func (ms *MemoryStorage) ResolveData(ChatId int64, UserId int64) *StorageRecord {
	_, ok := ms.Data[ChatId][UserId]
	if !ok {
		ms.Data[ChatId] = map[int64]*StorageRecord{}
	}
	record, ok := ms.Data[ChatId][UserId]

	if !ok {
		ms.Data[ChatId][UserId] = &StorageRecord{}
		record = ms.Data[ChatId][UserId]
	}

	return record
}

// SetData ...
func (ms *MemoryStorage) SetData(cid, uid int64, data PackType) error {
	record := ms.ResolveData(cid, uid)
	record.Data = data
	return nil
}

// GetData ...
func (ms *MemoryStorage) GetData(cid, uid int64) (PackType, error) {
	return ms.ResolveData(cid, uid).Data, nil
}

// SetState ...
func (ms *MemoryStorage) SetState(cid, uid int64, state string) error {
	record := ms.ResolveData(cid, uid)
	record.State = state
	return nil
}

// GetState ...
func (ms *MemoryStorage) GetState(cid, uid int64) (string, error) {
	return ms.ResolveData(cid, uid).State, nil
}

func (ms *MemoryStorage) Clear(cid, uid int64) error {
	ms.Data[cid][uid] = nil
	return nil
}

// Deletes all stored data
func (ms *MemoryStorage) Close() {
	for key, value := range ms.Data {
		for key := range value {
			delete(value, key)
		}
		delete(ms.Data, key)
	}
	ms.Data = DataType{}
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		Data: DataType{},
	}
}
