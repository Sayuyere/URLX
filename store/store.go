package store

type Store interface {
	Set(short, long string)
	Get(short string) (string, bool)
	Delete(short string)
}
