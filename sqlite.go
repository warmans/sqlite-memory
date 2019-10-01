package sqlite

import (
	"database/sql"

	"github.com/go-joe/joe"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const KeysTable = "keys"

// Memory is a joe.Option which is supposed to be passed to joe.New(…) to
// configure a new bot. The dsn is passed to Sqlite to open the DB. This
// could be a path to a file or in-memory DB:
//
// https://godoc.org/github.com/mattn/go-sqlite3#SQLiteDriver.Open
//
// If the DB already exists it will be opened, otherwise it will be
// created and initialised. If the file exists but cannot be opened
// its error will be deferred until the bot is actually started via
// its Run() function.
//
// Example usage:
//
//     b := joe.New("example",
//         sqlite.Memory(":memory:"),
//         …
//     )
//
func Memory(dsn string) joe.Module {
	return joe.ModuleFunc(func(conf *joe.Config) error {
		memory, err := NewMemory(dsn, WithLogger(conf.Logger("memory")))
		if err != nil {
			return err
		}
		conf.SetMemory(memory)
		return nil
	})

}

// NewMemory will create a new sqlite Memory instance and set
// up the DB to be used with a joe bot.
func NewMemory(dsn string, opts ...Option) (joe.Memory, error) {

	mem := &memory{logger: zap.NewNop()}
	for _, opt := range opts {
		err := opt(mem)
		if err != nil {
			return nil, err
		}
	}

	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open DB for given DSN")
	}
	mem.conn = conn

	if err := mem.init(); err != nil {
		return nil, err
	}

	return mem, nil
}

// memory is an implementation of a joe.Memory which stores all values in an embedded sqlite DB.
type memory struct {
	conn   *sql.DB
	logger *zap.Logger
}

func (s *memory) Set(key string, value []byte) error {
	_, err := s.conn.Exec(`
		INSERT INTO `+KeysTable+` ("key", "value") VALUES (?, ?)
		ON CONFLICT("key") DO UPDATE SET "value"=? WHERE "key" = ?
	`, key, value, value, key)

	return errors.Wrap(err, "unable to set key")
}

func (s *memory) Get(key string) ([]byte, bool, error) {
	var value []byte
	var count int

	err := s.conn.
		QueryRow(`SELECT "value", COUNT(1) FROM `+KeysTable+` WHERE "key" = ? `, key).
		Scan(&value, &count)

	return value, count == 1, errors.Wrap(err, "unable to get key")
}

func (s *memory) Delete(key string) (bool, error) {
	res, err := s.conn.Exec(`DELETE FROM `+KeysTable+` WHERE "key" = ?`, key)
	if err != nil {
		return false, errors.Wrap(err, "failed to delete key")
	}
	affected, err := res.RowsAffected()
	return affected > 0, errors.Wrap(err, "unable to delete key")
}

func (s *memory) Keys() ([]string, error) {
	res, err := s.conn.Query(`SELECT "key" FROM ` + KeysTable)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get key list")
	}
	defer res.Close()

	keys := make([]string, 0)
	for res.Next() {
		var key string
		if err := res.Scan(&key); err != nil {
			return nil, errors.Wrap(err, "failed to scan key")
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func (s *memory) Close() error {
	if s.conn == nil {
		return nil
	}
	return s.conn.Close()
}

func (s *memory) init() error {

	var count int
	if err := s.conn.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type ='table' AND name = ?`, KeysTable).Scan(&count); err != nil {
		return errors.Wrap(err, "failed to find existing tables")
	}
	if count == 1 {
		s.logger.Debug("DB is already initialised", zap.String("table", KeysTable))
		return nil
	}

	s.logger.Debug("Initialising DB", zap.String("table", KeysTable))
	_, err := s.conn.Exec(`CREATE TABLE "` + KeysTable + `" (key TEXT NOT NULL PRIMARY KEY, value BLOB)`)

	return errors.Wrap(err, "failed to initialise")
}
