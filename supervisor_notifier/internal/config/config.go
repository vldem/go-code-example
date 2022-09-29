package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const configPath = "config/spv_notif_cfg.yml"

type Cfg struct {
	Email struct {
		From    string   `yaml:"from"`
		To      string   `yaml:"to"`
		Subject string   `yaml:"subject"`
		Prog    MailProg `yaml:"mail_prog"`
	}
	Telegram struct {
		BotKey string `yaml:"botkey"`
		ChatId int64  `yaml:"chatid"`
	}
}

var AppConfig Cfg

type MailProg struct {
	Cmd  string   `yaml:"cmd"`
	Args []string `yaml:"args"`
}

func ReadConfig() error {
	var f *os.File
	var err error

	exPath, err := GetExecutablePath()
	if err != nil {
		return errors.New("[readconfig] cannot define executable path")
	}

	f, err = os.Open(exPath + "/../" + configPath)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Open(exPath + "/../../" + configPath)
			if err != nil {
				return errors.New("[readconfig] cannot find config file")
			}
		} else {
			return errors.Wrap(err, "[readconfig] error during reading config file")
		}
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&AppConfig)

	if err != nil {
		return errors.Wrap(err, "error during config yaml decoding")
	}

	if AppConfig.Email.Prog.Cmd == "" && AppConfig.Telegram.ChatId == 0 {
		return errors.New("mail program and telegram bot parameters are not defined")
	}

	return nil
}

func GetExecutablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)
	return exPath, nil
}
