package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Icon       bool `json:"enable_icon"` // 是否启用图标
	ByteOutput bool `json:"byteOutput"`  // 是否启用字节输出
	// SortType   string `json:"sortType"` 暂时不需要
}

var Global Config

func defaultConfig() Config {
	return Config{
		Icon:       true,
		ByteOutput: false,
		// SortType:   "name",
	}
}

func configFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "nls", "config.json"), nil
}

func Load() {
	Global = defaultConfig()

	path, err := configFilePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR:unable to locate config dir, using defaults: %s\n", err)
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "ERR:failed to read config file %s, using defaults: %s\n", path, err)
		}
		return
	}

	if err := json.Unmarshal(data, &Global); err != nil {
		fmt.Fprintf(os.Stderr, "ERR:failed to parse config file %s, using defaults: %s\n", path, err)
		Global = defaultConfig()
	}
}
