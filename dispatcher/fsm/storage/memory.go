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
	record := d[ChatId][UserId]

	if *record == EmptyRecord {
		d[ChatId][UserId] = &MemoryStorageRecord{}
		record = d[ChatId][UserId]
	}

	return record
}

// SetData ...
func (ms *MemoryStorage) SetData(ChatID int64, UserID int64, data *PackType) {
	record := ms.ResolveData(ChatID, UserID)
	record.Data = data
}

// GetData ...
func (ms *MemoryStorage) GetData(ChatID int64, UserID int64) *PackType {
	return ms.ResolveData(ChatID, UserID).Data
}

// SetState ...
func (ms *MemoryStorage) SetState(ChatID int64, UserID int64, state string) {
	record := ms.ResolveData(ChatID, UserID)
	record.State = state
}

// GetState ...
func (ms *MemoryStorage) GetState(ChatID int64, UserID int64) string {
	return ms.ResolveData(ChatID, UserID).State
}

func (ms *MemoryStorage) Clear() {
	ms.Data = nil
}

var (
	EmptyRecord = MemoryStorageRecord{}
)
