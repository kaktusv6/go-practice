package config

import (
	"os"
)

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func Init[Storage any](pathToConfig string, storage *Storage) error {
	rawYAML, err := os.ReadFile(pathToConfig)
	if err != nil {
		return errors.WithMessage(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, storage)
	if err != nil {
		return errors.WithMessage(err, "parsing yaml")
	}

	return nil
}
