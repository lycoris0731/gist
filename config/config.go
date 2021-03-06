package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Core   Core
	Gist   Gist
	Flag   Flag
	Screen Screen
}

type Core struct {
	Editor        string `toml:"editor"`
	SelectCmd     string `toml:"selectcmd"`
	TomlFile      string `toml:"tomlfile"`
	User          string `toml:"user"`
	ShowIndicator bool   `toml:"show_indicator"`
	BaseURL       string `toml:"base_url"`
}

type Gist struct {
	Token    string `toml:"token"`
	Dir      string `toml:"dir"`
	StdinExt string `toml:"default_stdin_ext"`
}

type Flag struct {
	ShowSpinner bool `toml:"show_spinner"`
	Verbose     bool `toml:"verbose"`
	OpenURL     bool `toml:"open_url"`
	NewPrivate  bool `toml:"new_private"`
	OpenBaseURL bool `toml:"open_base_url"`

	// TODO
	Sort string `toml:"sort"`
	Only string `toml:"only"`

	EditDesc         bool
	OpenStarredItems bool
	FromClipboard    bool
}

type Screen struct {
	ShowPrivateSymbol bool `toml:"show_private_symbol"`
}

var Conf Config

func GetDefaultDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	default:
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	case "windows":
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data")
		}
	}
	dir = filepath.Join(dir, "gist")

	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return dir, fmt.Errorf("cannot create directory: %v", err)
	}

	return dir, nil
}

func (cfg *Config) LoadFile(file string) error {
	_, err := os.Stat(file)
	if err == nil {
		_, err := toml.DecodeFile(file, cfg)
		if err != nil {
			return err
		}
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	cfg.Gist.Token = os.Getenv("GITHUB_TOKEN")
	cfg.Core.Editor = os.Getenv("EDITOR")
	if cfg.Core.Editor == "" {
		cfg.Core.Editor = "vim"
	}
	cfg.Core.SelectCmd = "fzf-tmux --multi:fzf --multi:peco"
	cfg.Core.TomlFile = file
	cfg.Core.User = os.Getenv("USER")
	cfg.Core.ShowIndicator = true
	cfg.Core.BaseURL = "https://gist.github.com"
	dir := filepath.Join(filepath.Dir(file), "files")
	os.MkdirAll(dir, 0700)
	cfg.Gist.Dir = dir
	cfg.Flag.ShowSpinner = true
	cfg.Flag.Verbose = true
	cfg.Flag.OpenURL = false
	cfg.Flag.NewPrivate = false
	cfg.Flag.OpenBaseURL = false
	cfg.Screen.ShowPrivateSymbol = false

	return toml.NewEncoder(f).Encode(cfg)
}
