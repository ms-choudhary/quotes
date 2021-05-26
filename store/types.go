package store

type Quote struct {
	ID     int
	Text   string
	Author string
	Tags   []string
}

type DB interface {
	Ping() error
	Create(q Quote) (Quote, error)
	Get(id int) (Quote, error)
	GetAll() ([]Quote, error)
	GetRandom() (Quote, error)
	Clean() error
}
