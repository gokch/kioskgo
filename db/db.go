package db

type DB interface {
	Set(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Exists(key []byte) (bool, error)
	Delete(key []byte) error
}

type TX interface {
	Set(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Commit() error
	Rollback() error
}

type Bulk interface {
	Set(key, value []byte) error
	Delete(key []byte) error

	Flush() error
	Revert() error
}
