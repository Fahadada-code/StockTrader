package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	Conn *sql.DB
}

func NewPostgresDB(url string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	pg := &PostgresDB{Conn: db}
	if err := pg.initSchema(); err != nil {
		return nil, err
	}

	return pg, nil
}

func (pg *PostgresDB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS market_data (
		id SERIAL PRIMARY KEY,
		symbol VARCHAR(10) NOT NULL,
		price DECIMAL(18, 4) NOT NULL,
		volume BIGINT NOT NULL,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_symbol_timestamp ON market_data (symbol, timestamp DESC);
	
	CREATE TABLE IF NOT EXISTS anomalies (
		id SERIAL PRIMARY KEY,
		symbol VARCHAR(10) NOT NULL,
		type VARCHAR(50) NOT NULL,
		confidence DECIMAL(5, 4) NOT NULL,
		description TEXT,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := pg.Conn.Exec(schema)
	return err
}

func (pg *PostgresDB) SaveQuote(symbol string, price float64, volume int64) error {
	_, err := pg.Conn.Exec(
		"INSERT INTO market_data (symbol, price, volume) VALUES ($1, $2, $3)",
		symbol, price, volume,
	)
	return err
}

func (pg *PostgresDB) GetHistoricalData(symbol string, limit int) (*sql.Rows, error) {
	return pg.Conn.Query(
		"SELECT price, volume, timestamp FROM market_data WHERE symbol = $1 ORDER BY timestamp DESC LIMIT $2",
		symbol, limit,
	)
}
