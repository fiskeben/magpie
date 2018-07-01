package main

import (
	"database/sql"
	"log"
	"time"

	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

// EventsAccessor can interface with Docker.
type EventsAccessor interface {
	Events(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error)
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
}

// EventsClient is a wrapper around the Docker API client.
type EventsClient struct {
	client *client.Client
}

// Events wraps the events API.
func (d *EventsClient) Events(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error) {
	return d.client.Events(ctx, options)
}

// ContainerInspect wraps the container inspect API.
func (d *EventsClient) ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	return d.client.ContainerInspect(ctx, containerID)
}

func listen(ctx context.Context, c EventsAccessor, db *sql.DB, signals chan bool) error {
	events, errors := c.Events(ctx, types.EventsOptions{})
	for {
		select {
		case e := <-events:
			handleEvent(ctx, e, c, db)
		case err := <-errors:
			return err
		case _ = <-signals:
			return nil
		}
	}
}

func handleEvent(ctx context.Context, msg events.Message, c EventsAccessor, db *sql.DB) {
	if !shouldHandle(msg) {
		return
	}

	whitelist, err := GetWhitelist(db)
	if err != nil {
		log.Fatalf("failed to get whitelist: %v", err)
		return
	}

	blacklist, err := GetBlacklist(db)
	if err != nil {
		log.Fatalf("failed to get blacklist: %v", err)
		return
	}

	container, err := c.ContainerInspect(ctx, msg.ID)
	if err != nil {
		log.Fatalf("failed to inspect container: %v", err)
		return
	}

	data, err := dataFromContainer(ctx, container, whitelist, blacklist)
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := SaveConfiguration(db, *data); err != nil {
		log.Fatal(err.Error())
	}
}

func shouldHandle(msg events.Message) bool {
	return msg.Type == events.ContainerEventType && msg.Action == "start" && msg.Status == "start"
}

func dateToTimestamp(date string) (int64, error) {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return -1, err
	}
	return t.Unix(), nil
}
