package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type config struct {
	ZipFile          string
	Format           string
	Output           string
	TableName        string
	LocalCodeField   string
	ZipCodeField     string
	PrefecturesField string
	CityField        string
	TownField        string
}

type configJSON struct {
	Output string
}

type configSQL struct {
	Output string
}

func getDefaultConfig() config {
	return config{
		Format:           "json",
		Output:           "jzip.json",
		TableName:        "zipcode",
		LocalCodeField:   "local",
		ZipCodeField:     "zip",
		PrefecturesField: "pref",
		CityField:        "city",
		TownField:        "town",
	}
}

func initConfig(path string) config {
	def := getDefaultConfig()
	if path == "" || path[0] != '/' {
		wd, err := os.Getwd()
		if err != nil {
			return def
		}
		path = filepath.Join(wd, path)
	}
	inf, err := os.Stat(path)
	if err != nil {
		return def
	}
	if inf.IsDir() {
		path = filepath.Join(path, "jzip.toml")
		inf, err = os.Stat(path)
		if err != nil {
			return def
		}
	}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return def
	}
	c := config{}
	if err := toml.Unmarshal(buf, &c); err != nil {
		return def
	}
	if c.Format == "" {
		c.Format = def.Format
	}
	if c.Output == "" {
		if c.Format == "sql" {
			c.Output = "jzip.sql"
		} else {
			c.Output = "jzip.json"
		}
	}
	if c.TableName == "" {
		c.TableName = def.TableName
	}
	if c.LocalCodeField == "" {
		c.LocalCodeField = def.LocalCodeField
	}
	if c.ZipCodeField == "" {
		c.ZipCodeField = def.ZipCodeField
	}
	if c.PrefecturesField == "" {
		c.PrefecturesField = def.PrefecturesField
	}
	if c.CityField == "" {
		c.CityField = def.CityField
	}
	if c.TownField == "" {
		c.TownField = def.TownField
	}
	return c
}
