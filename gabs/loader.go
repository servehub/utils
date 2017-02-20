package gabs

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/fatih/color"
	"github.com/ghodss/yaml"

	"github.com/servehub/utils/mergemap"
)

func LoadYamlFile(path string) (*Container, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New(color.RedString("Manifest file `%s` not found: %v", path, err))
	}

	if result, err := ParseYaml(data); err != nil {
		return nil, fmt.Errorf("Error on load file %s: %v", path, err)
	} else {
		return result, nil
	}
}

func ParseYaml(data []byte) (*Container, error) {
	if jsonData, err := yaml.YAMLToJSON(data); err != nil {
		return nil, errors.New(color.RedString("Error on parse yaml: %v", err))
	} else {
		return ParseJSON(jsonData)
	}
}

func (g *Container) WithFallback(original *Container) error {
	merged, err := mergemap.Merge(original.Data().(map[string]interface{}), g.Data().(map[string]interface{}))
	if err != nil {
		return err
	}

	_, err = g.Set(merged)
	return err
}

func (g *Container) WithFallbackYaml(data []byte) error {
	original, err := ParseYaml(data)
	if err != nil {
		return err
	}

	return g.WithFallback(original)
}

func (g *Container) WithFallbackYamlFile(path string) error {
	original, err := LoadYamlFile(path)
	if err != nil {
		return err
	}

	return g.WithFallback(original)
}
