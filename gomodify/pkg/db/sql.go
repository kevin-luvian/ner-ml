package db

import (
	"context"
	"database/sql"
	"strings"
	"sync"
	"time"

	nurl "net/url"

	"github.com/jmoiron/sqlx"
	"github.com/kevin-luvian/gomodify/pkg/logging"
	_ "github.com/lib/pq"
)

type Config struct {
	SourceURL             string
	Retries               int
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
}

type DB struct {
	Instance *sqlx.DB
}

func New(cfg Config) (*DB, error) {
	if cfg.Retries < 1 || cfg.Retries > 10 {
		cfg.Retries = 3
	}

	db, err := Open(cfg.SourceURL, cfg.Retries)
	if err != nil {
		return nil, err
	}

	if cfg.MaxIdleConnections > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConnections)
	}

	if cfg.MaxOpenConnections > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConnections)
	}

	if cfg.ConnectionMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)
	}

	return &DB{Instance: db}, nil
}

func Open(sourceURL string, retries int) (*sqlx.DB, error) {
	var (
		db  *sqlx.DB
		err error
	)

	u, err := nurl.Parse(sourceURL)
	if err != nil {
		return nil, err
	}

	driver := u.Scheme

	for retries > 0 {
		retries--

		db, err = sqlx.Connect(driver, sourceURL)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			logging.Infoln("database connected")
			return db, nil
		}

		logging.Warnf("failed to connect to %s with error %s", sourceURL, err.Error())

		if retries > 0 {
			logging.Warnln("retrying to connect...")
			time.Sleep(time.Second * 3)
		}
	}

	return nil, err
}

func (db *DB) ExecContect(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.Instance.ExecContext(ctx, db.Instance.Rebind(query), args...)
}

func (db *DB) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	qValues := make([]string, len(args))
	for i := range args {
		qValues[i] = "?"
	}

	query = strings.ReplaceAll(query, ":values", strings.Join(qValues, ", "))

	return db.Instance.QueryRowxContext(ctx, db.Instance.Rebind(query), args...)
}

func (db *DB) Get(ctx context.Context, query string, param GetDBParam, target interface{}) error {
	query, args := param.WhereQuery().SortQuery().GetQuery(query)
	return db.Instance.SelectContext(ctx, target, db.Instance.Rebind(query), args...)
}

func (db *DB) GetWithCount(ctx context.Context, mainQuery, countQuery string, param GetDBParam, target interface{}, totalRecord *int) error {
	var (
		wg      sync.WaitGroup
		errChan = make(chan error)
	)

	if !param.DisableGet {
		wg.Add(1)
		go func() {
			defer wg.Done()

			errChan <- db.Get(ctx, mainQuery, param, target)
		}()
	}

	if !param.DisableCount {
		wg.Add(1)
		go func() {
			defer wg.Done()

			errChan <- db.getCount(ctx, countQuery, param, totalRecord)
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	// catch any error and blocking
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// getCount use GetNamedQueryCount with get context
func (db *DB) getCount(ctx context.Context, mainQuery string, param GetDBParam, total *int) error {
	query, args := param.WhereQuery().SortQuery().GetQuery(mainQuery)

	var mTotal int
	err := db.Instance.GetContext(ctx, &mTotal, db.Instance.Rebind(query), args...)
	if err != nil {
		return err
	}

	*total = mTotal

	return nil
}
