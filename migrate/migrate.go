package migrate

import (
	"fmt"

	"github.com/fiskeben/magpie/migrate/steps"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/golang-migrate/migrate/source/go_bindata"
)

// Migrate migrates the database
func Migrate(connection string, version int) error {
	s := bindata.Resource(steps.AssetNames(),
		func(name string) ([]byte, error) {
			return steps.Asset(name)
		})

	driver, err := bindata.WithInstance(s)
	if err != nil {
		return fmt.Errorf("failed to create migration data driver: %v", err)
	}

	m, err := migrate.NewWithSourceInstance("go-bindata", driver, connection)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	if version > -1 {
		return m.Force(version)
	}

	if err = m.Up(); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	return nil
}
