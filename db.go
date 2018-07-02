package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
)

const connectionTemplate = "postgres://%s:%s@%s:%s/%s?sslmode=disable"

const saveSQL = `INSERT INTO configurations (service_name, version, created_at, configuration) 
SELECT 
	CAST($1 AS VARCHAR), 
	CAST($2 AS VARCHAR), 
	to_timestamp($3), 
	$4
WHERE NOT EXISTS (SELECT id FROM configurations
	WHERE service_name = $1 AND version = $2 AND created_at = to_timestamp($3))`

const whitelistSQL = `SELECT name FROM whitelists`

const blacklistSQL = `SELECT name FROM blacklists`

type valueMap map[string]string

func (v valueMap) Value() (driver.Value, error) {
	return json.Marshal(v)
}

// MakeConnectionString creates a database connection string from a configuration.
func MakeConnectionString(config ConfigMap) string {
	username := config.Get("username", "postgres")
	password := config.Get("password", "postgres")
	host := config.Get("host", "localhost")
	port := config.Get("port", "5432")
	dbname := config.Get("database", "postgres")

	return fmt.Sprintf(connectionTemplate, username, password, host, port, dbname)
}

// SaveConfiguration saves a single configuration to the database.
func SaveConfiguration(db *sql.DB, config ContainerData) error {
	return SaveConfigurations(db, []ContainerData{config})
}

// SaveConfigurations takes a list of configurations and store them in the database.
func SaveConfigurations(db *sql.DB, configs []ContainerData) error {
	for _, c := range configs {
		if _, err := db.Exec(saveSQL, c.name, c.version, c.created, valueMap(c.config)); err != nil {
			return fmt.Errorf("failed to insert config: %v", err)
		}
	}

	return nil
}

// GetWhitelist gets all whitelisted vars.
func GetWhitelist(db *sql.DB) ([]string, error) {
	return getList(db, whitelistSQL)
}

// GetBlacklist gets all blacklisted vars.
func GetBlacklist(db *sql.DB) ([]string, error) {
	return getList(db, blacklistSQL)
}

func getList(db *sql.DB, sql string) ([]string, error) {
	rows, err := db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("failed to get whitelist: %v", err)
	}
	defer rows.Close()

	res := []string{}

	for rows.Next() {
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("failed to iterate result set: %v", err)
		}

		var item string
		if err = rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("failed to read row: %v", err)
		}
		res = append(res, item)
	}
	return res, nil
}
