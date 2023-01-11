package generate

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type KittenList struct {
	Kittens []Kittens `yaml:"Kittens"`
}
type Kittens struct {
	Name    string `yaml:"name"`
	Path    string `yaml:"path"`
	Compile string `yaml:"compile"`
}

func NewKittenList() (*KittenList, error) {
	var list KittenList
	filename, err := filepath.Abs("cmd\\generate\\kittens.yml")
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &list)
	if err != nil {
		return nil, err
	}
	return &list, error(nil)
}

func (k *KittenList) GetKittenNames() []string {
	var names []string
	for _, v := range k.Kittens {
		names = append(names, v.Name)
	}
	return names
}

func (k *KittenList) GetKittenPaths() []string {
	var paths []string
	for _, v := range k.Kittens {
		paths = append(paths, v.Path)
	}
	return paths
}

func (k *KittenList) GetKittensPath(name string) (string, error) {
	for _, v := range k.Kittens {
		if v.Name == name {
			return v.Path, error(nil)
		}
	}
	return "", errors.New("kitten not found")
}

func (k *KittenList) GetKittensCompile(name string) (string, error) {
	for _, v := range k.Kittens {
		if v.Name == name {
			return v.Compile, error(nil)
		}
	}
	return "", errors.New("kitten not found")
}
