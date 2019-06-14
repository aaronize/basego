package db

type DB interface {
	InitDB() error
	Close() error
}

func NewDB(driver string) DB {
	switch driver {
	case "mysql":
		return &mysql{

		}
	case "mongodb":

		return nil
	case "leveldb":

		return nil
	case "sqlite":

		return nil
	case "redis":

		return nil
	default:

		return nil
	}
}

// mysql
type mysql struct {

}

func (m *mysql) InitDB() error {

	return nil
}

func (m *mysql) Close() error {

	return nil
}

// mongodb
type mongodb struct {

}

func (m *mongodb) InitDB() error {
	return nil
}

func (m *mongodb) Close() error {
	return nil
}

// leveldb
type leveldb struct {

}

func (l *leveldb) InitDB() error {
	return nil
}

func (l *leveldb) Close() error {
	return nil
}

// sqlite
type sqlite struct {

}

func (s *sqlite) InitDB() error {
	return nil
}

func (s *sqlite) Close() error {
	return nil
}

// redis
type redis struct {

}

func (r *redis) InitDB() error {
	return nil
}

func (r *redis) Close() error {
	return nil
}
