package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultConfigFile = "~/.magpie.conf"

var errFileNotFound = errors.New("configuration file not found")

func readConfig(path string) (*ConfigMap, error) {
	f, err := readFile(path)
	if err != nil {
		if err == errFileNotFound {
			log.Println("no configuration file, using defaults")
			return &ConfigMap{}, nil
		}
		return nil, err
	}
	defer f.Close()

	return parseConfig(f)
}

func readFile(path string) (*os.File, error) {
	if strings.TrimSpace(path) == "" {
		path = defaultConfigFile
	}

	f, err := os.Open(path)
	if err != nil {
		if path == defaultConfigFile {
			return nil, errFileNotFound
		}
		return nil, fmt.Errorf("config file '%s' not found", path)
	}
	return f, nil
}

func parseConfig(r io.Reader) (*ConfigMap, error) {
	res := ConfigMap{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("illegal config format '%s'", line)
		}
		res[parts[0]] = parts[1]
	}

	return &res, nil
}
