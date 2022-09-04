package models

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v2"
)

// RegAddr возвращает адрес для хендлера регистрации
func (c Config) RegAddr() string {
	return c.HostAddr + "/api/v1/registration"
}

// AuthAddr возвращает адрес для хендлера авторизации
func (c Config) AuthAddr() string {
	return c.HostAddr + "/api/v1/auth"
}

// ActionAddr возвращает адрес для хендлера выполнения действий(обновление, добавление...)
func (c Config) ActionAddr() string {
	return c.HostAddr + "/api/v1/items"
}

func (c *Config) Init(filename string) error {
	f, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(f)
	decoder := yaml.NewDecoder(reader)

	return decoder.Decode(c)
}

// Valid проверяет заполнение полей и валидность структуры для обработки
func (r UserRequest) Valid() bool {
	return r.Login != "" && r.Password != ""
}

// Valid проверяет заполнение полей и валидность структуры для обработки
func (r UserData) Valid() bool {
	return r.Data != ""
}

// Valid проверяет заполнение полей и валидность структуры для обработки
func (r DeleteRequest) Valid() bool {
	return r.ID != ""
}
