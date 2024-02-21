package repo

type InMemoryStorageI interface {
	Set(string, string) error
	Get(string) (interface{}, error)
	SetWithTTL(string, string, int) error
}
