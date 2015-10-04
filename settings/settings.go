package settings

import (
	"encoding/json"
	"os"
)

type settings struct {
	PlayerCount int
	PlayerNames [4]string
	PlayerTypes [4]PlayerType
	IPs         [4]string
	Ports       [4]string
	Language    int
}

var Settings = &settings{
	3,
	[4]string{"1", "2", "3", "4"},
	[4]PlayerType{Human, Human, Human, Human},
	[4]string{"127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1"},
	[4]string{"5555", "5555", "5555", "5555"},
	0,
}

const settingsPath = "./settings.txt"

func (s *settings) Load() error {
	file, err := os.Open(settingsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var loaded settings
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&loaded)
	if err != nil {
		return err
	}
	*s = loaded

	return nil
}

func (s *settings) Save() error {
	file, err := os.Create(settingsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(s)
}

type PlayerType int

const (
	Human         PlayerType = iota
	AI            PlayerType = iota
	NetworkPlayer PlayerType = iota
)
