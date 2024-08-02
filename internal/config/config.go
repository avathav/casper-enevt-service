package config

import (
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

var config *koanf.Koanf

var requiredParams = []string{
	"MYSQL.HOST",
	"MYSQL.DATABASE",
	"MYSQL.USERNAME",
	"MYSQL.PASSWORD",
}

func loadConfig() {
	if config == nil {
		config = koanf.New(".")
		if err := config.Load(file.Provider("config/config.yml"), yaml.Parser()); err != nil {
			log.Panicf("could not load configuration")
		}

		keyMap := config.KeyMap()
		for _, k := range requiredParams {
			if _, ok := keyMap[k]; !ok {
				log.Fatalf("missing configuration parameter %s", k)
			}
		}
	}
}

func GetParam(key string) any {
	loadConfig()

	return config.Get(key)
}

func GetString(key string) string {
	loadConfig()

	return config.String(key)
}

func GetInt(key string) int {
	loadConfig()

	return config.Int(key)
}

func GetBool(key string) bool {
	loadConfig()

	return config.Bool(key)
}

func GetStringOrFallback(key string, fallback string) string {
	v := GetString(key)

	if v == "" {
		return fallback
	}

	return v
}

func GetIntOrFallback(key string, fallback int) int {
	v := GetInt(key)

	if v == 0 {
		return fallback
	}

	return v
}
