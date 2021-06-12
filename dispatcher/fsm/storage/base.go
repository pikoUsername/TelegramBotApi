package storage

type PackType map[string]interface{}

// Simple storage interface for saving data,
// and uses for save FSM data
type Storage interface {
	SetData(int64, int64, *PackType)
	GetData(int64, int64) *PackType
	SetState(int64, int64, string)
	GetState(int64, int64) string
	Clear()
}
