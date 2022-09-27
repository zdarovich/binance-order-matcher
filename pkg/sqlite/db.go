package sqlite

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"net/url"
	"os"
)

func Create(dbURL string) error {
	u, err := url.Parse(dbURL)
	if err != nil {
		return errors.Wrap(err, "malformed database URL")
	}
	_, err = os.Create(u.Path)
	if os.IsExist(err) {
		return nil
	}
	return err
}
func NewConnection(dbURL string) (*sqlx.DB, error) {
	u, err := url.Parse(dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "malformed database URL")
	}

	switch u.Scheme {
	case "sqlite3":
		db, err := sqlx.Connect("sqlite3", u.Host)
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		return db, nil
	default:
		return nil, errors.New(fmt.Sprintf("unsupported dabase url scheme %s", u.Scheme))
	}
}
