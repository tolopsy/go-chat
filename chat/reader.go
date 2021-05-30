package main

/* Reads from config.json */

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type JsonReader struct {
	ClientSecret string
	ClientId     string
}

func (configReader *JsonReader) reads() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	file, err := os.Open(filepath.Join(workingDir, "config.json"))
	defer file.Close()

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(configReader)
	if err != nil {
		return err
	}

	return nil
}
