package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"context"

	"github.com/docker/docker/client"
	"github.com/fiskeben/magpie/migrate"
)

var Version string
var BuildDate string

// ConfigMap stores data read from the configuration file.
type ConfigMap map[string]string

// Get reads the value from the map and returns it.
// If the key doesn't exist a fallback value will be returned instead.
func (c ConfigMap) Get(key, fallback string) string {
	val, ok := c[key]
	if ok {
		return val
	}
	return fallback
}

func main() {
	ctx := context.Background()

	init := flag.Bool("init", false, "Scrape the currently running Docker containers at startup")
	configPath := flag.String("config", "~/.magpie.conf", "Path to the magpie configuration file")
	migrateFlag := flag.Bool("migrate", false, "Migrate database and exit")
	migrateVersion := flag.Int("force", -1, "Force database migrations to this version")
	versionFlag := flag.Bool("version", false, "Print version and build info")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("magpie version %s (built %s)\n", Version, BuildDate)
		os.Exit(0)
	}

	config, err := readConfig(*configPath)
	if err != nil {
		logAndExit(err)
	}

	connectionString := MakeConnectionString(*config)

	if *migrateFlag {
		if err := migrate.Migrate(connectionString, *migrateVersion); err != nil {
			logAndExit(err)
		}
		os.Exit(0)
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		logAndExit(err)
	}

	c, err := client.NewEnvClient()
	if err != nil {
		logAndExit(err)
	}

	EventsClient := &EventsClient{client: c}
	containerClient := &ContainerClient{client: c}

	if *init {
		data, err := GetConfigurations(ctx, containerClient, db)
		if err != nil {
			logAndExit(err)
		}
		if err = SaveConfigurations(db, data); err != nil {
			logAndExit(err)
		}
	}

	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		_ = <-signals
		done <- true
	}()

	listen(ctx, EventsClient, db, done)
}

func logAndExit(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
