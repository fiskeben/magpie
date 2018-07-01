package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var globalBlacklist = []string{
	"HOME",
	"PATH",
}

var globalWhitelist = []string{
	"NODE_VERSION",
}

// ContainerData holds information about a single container's configuration.
type ContainerData struct {
	name    string
	version string
	config  map[string]string
	created int64
}

type ContainerInquirer interface {
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
}

type ContainerClient struct {
	client *client.Client
}

func (c *ContainerClient) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	return c.client.ContainerList(ctx, options)
}

func (c *ContainerClient) ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	return c.client.ContainerInspect(ctx, containerID)
}

// GetConfigurations reads configurations from all running containers
// on the host. Configuration data will be filtered so that blacklisted
// values are removed and non-whitelisted values are masked.
func GetConfigurations(ctx context.Context, c ContainerInquirer, db *sql.DB) ([]ContainerData, error) {

	containers, err := c.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	whitelist, err := GetWhitelist(db)
	if err != nil {
		return nil, err
	}

	blacklist, err := GetBlacklist(db)
	if err != nil {
		return nil, err
	}

	result := make([]ContainerData, 0)

	for _, container := range containers {
		res, err := c.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			panic(err)
		}
		data, err := dataFromContainer(ctx, res, whitelist, blacklist)
		if err != nil {
			return nil, err
		}
		result = append(result, *data)
	}
	return result, nil
}

func dataFromContainer(ctx context.Context, container types.ContainerJSON, whitelist, blacklist []string) (*ContainerData, error) {
	env := listToMap(container.Config.Env)
	env = filter(env, whitelist, blacklist)

	image, version, err := splitImageName(container.Config.Image)
	if err != nil {
		return nil, err
	}

	timestamp, err := dateToTimestamp(container.State.StartedAt)
	if err != nil {
		return nil, err
	}

	data := ContainerData{
		name:    image,
		version: version,
		config:  env,
		created: timestamp,
	}

	return &data, nil
}

func listToMap(env []string) map[string]string {
	if env == nil || len(env) == 0 {
		return map[string]string{}
	}

	res := make(map[string]string)
	for _, item := range env {
		parts := strings.Split(item, "=")
		res[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return res
}

func filter(env map[string]string, whitelist, blacklist []string) map[string]string {
	res := make(map[string]string)

	for name, value := range env {
		if contains(blacklist, name) {
			continue
		}
		if contains(whitelist, name) {
			res[name] = value
			continue
		}
		res[name] = mask(value)
	}

	return res
}

func contains(list []string, item string) bool {
	for _, i := range list {
		if item == i {
			return true
		}
	}
	return false
}

func mask(value string) string {
	if len(value) < 6 {
		return "*****"
	}
	length := len(value) - 4
	m := strings.Repeat("*", length)
	return value[:2] + m + value[length+2:]
}

func splitImageName(name string) (string, string, error) {
	parts := strings.Split(name, ":")
	if len(parts) == 2 {
		return parts[0], parts[1], nil
	}
	if len(parts) == 1 {
		return parts[0], "", nil
	}
	return "", "", fmt.Errorf("image name '%s' has unexpected format", name)
}
