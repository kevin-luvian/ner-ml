package db

import (
	"testing"
	"time"

	"github.com/kevin-luvian/gomodify/pkg/assert"
	_ "github.com/proullon/ramsql/driver"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name string
		arg  Config
	}{{
		name: "success",
		arg: Config{
			SourceURL: "ramsql://dsn",
			Retries:   0,
		},
	}}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db, err := New(tt.arg)
			assert.NoError(t, err)
			assert.NotNil(t, db)

			err = db.Instance.Ping()
			assert.NoError(t, err)
		})
	}
}

func TestConnect(t *testing.T) {
	type args struct {
		sourceURL string
		retries   int
	}
	testCases := []struct {
		name string
		arg  args
	}{
		{
			name: "success",
			arg: args{
				sourceURL: "ramsql://dsn",
				retries:   3,
			},
		},
		{
			name: "success retries 1",
			arg: args{
				sourceURL: "ramsql://dsn",
				retries:   1,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Open(tt.arg.sourceURL, tt.arg.retries)
			assert.NoError(t, err)
			assert.NotNil(t, db)

			err = db.Ping()
			assert.NoError(t, err)
		})
	}
}

func TestSelect(t *testing.T) {
	cfg := Config{
		SourceURL:             "ramsql://dsn",
		Retries:               3,
		ConnectionMaxLifetime: 10 * time.Second,
		MaxIdleConnections:    2,
		MaxOpenConnections:    10,
	}

	// connect
	db, err := New(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = db.Instance.Ping()
	assert.NoError(t, err)

	// create table
	_, err = db.Instance.Exec(`CREATE TABLE address (id BIGSERIAL PRIMARY KEY, street TEXT, street_number INT);`)
	assert.NoError(t, err)

	type Address struct {
		Street       string `db:"street"`
		StreetNumber int    `db:"street_number"`
	}

	var returned Address
	expected := Address{
		Street:       "hugo",
		StreetNumber: 32,
	}

	// insert
	q := db.Instance.Rebind("INSERT INTO address (street, street_number) VALUES ($1, $2)")
	_, err = db.Instance.Exec(q, expected.Street, expected.StreetNumber)
	assert.NoError(t, err)

	// get
	err = db.Instance.Get(&returned, `select street, street_number from address where street='hugo'`)
	assert.NoError(t, err)
	assert.Equal(t, expected, returned)
}
