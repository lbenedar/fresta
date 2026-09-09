package config

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	_ "embed"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Commands []*discordgo.ApplicationCommand

	GuildId string
	Errs    chan error
}

const RelPath = "cmd/api/discord_bot/config/"

type FileSystemPath struct {
	path string
	stat os.FileInfo
	file embed.FS
}

//go:embed commands.json
var commandFile embed.FS

var configFiles = map[string]*FileSystemPath{
	"Commands": &FileSystemPath{path: "commands.json", file: commandFile},
}

func ParseConfigs() (*Config, error) {
	cfg := &Config{
		Commands: make([]*discordgo.ApplicationCommand, 0),

		Errs: make(chan error, 5),
	}

	for name, fs := range configFiles {
		path := RelPath + fs.path
		fs.stat, _ = os.Stat(path)

		buffer, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		err = insertByName(cfg, buffer, name)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func insertByName(cfg *Config, jsData []byte, name string) error {
	cfgElem := reflect.ValueOf(cfg).Elem()

	fieldValue := cfgElem.FieldByName(name)
	if fieldValue.Kind() != reflect.Slice {
		return ErrWrongKindPassed
	}

	if fieldValue.IsValid() && fieldValue.CanSet() {
		err := json.Unmarshal(jsData, fieldValue.Addr().Interface())
		if err != nil {
			return err
		}
	}

	return nil
}

func CheckIfFilesChanged() ([]string, error) {
	changedFiles := make([]string, 0)

	for name, fs := range configFiles {
		path := RelPath + fs.path
		stat, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		if fs.stat == nil || stat.ModTime() != fs.stat.ModTime() || stat.Size() != fs.stat.Size() {
			changedFiles = append(changedFiles, name)
			fs.stat = stat
		}
	}
	return changedFiles, nil
}

func UpdateConfigsByNames(cfg *Config, changedFiles []string) bool {
	length := len(changedFiles)
	changed := false

	for i := range changedFiles {
		confName := changedFiles[length-i-1]
		fs := configFiles[confName]
		path := RelPath + fs.path

		buffer, err := os.ReadFile(path)
		if err != nil {
			cfg.Errs <- ConfigError{err: err, configName: confName}
			fs.stat = nil
			continue
		}

		fmt.Println(buffer)

		err = insertByName(cfg, buffer, confName)
		if err != nil {
			cfg.Errs <- ConfigError{err: err, configName: confName}
			fs.stat = nil
			continue
		}

		changed = true
	}
	return changed
}

func InitCommands(session *discordgo.Session, cfg *Config) error {
	_, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, cfg.GuildId, []*discordgo.ApplicationCommand{})
	if err != nil {
		return err
	}
	_, err = session.ApplicationCommandBulkOverwrite(session.State.User.ID, "", []*discordgo.ApplicationCommand{})
	if err != nil {
		return err
	}

	return createCommands(session, cfg)
}

func createCommands(session *discordgo.Session, cfg *Config) error {
	commands := cfg.Commands
	for _, v := range commands {
		session.ApplicationCommandCreate(session.State.User.ID, cfg.GuildId, v)
	}
	return nil
}
