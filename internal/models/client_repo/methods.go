package client_repo

import (
	"github.com/azazel3ooo/keeper/internal/models"
)

// NewClient возвращает Client с переданными параметрами
func NewClient(opts ...func(client *Client)) *Client {
	c := &Client{}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithStorage добавляет переданный models.ClientStorable для клиента
func WithStorage(store models.ClientStorable) func(*Client) {
	return func(c *Client) {
		c.store = store
	}
}

// WithConfig добавляет переданный models.Config для клиента
func WithConfig(cfg models.Config) func(*Client) {
	return func(c *Client) {
		c.cfg = cfg
	}
}

// WithClient добавляет переданный models.ClientHttpInterface для клиента
func WithClient(cl models.ClientHttpInterface) func(*Client) {
	return func(c *Client) {
		c.cl = cl
	}
}

// UpdateToken обновляет токен клиента
func (c *Client) UpdateToken(newToken string) {
	c.token = newToken
}

// ReadyForActions проверяет, что клиент готов к работе (токен не пустой)
func (c Client) ReadyForActions() bool {
	if c.token == "" {
		return false
	}

	return true
}

// RegistrationAddress возвращает адрес для метода регистрации
func (c Client) RegistrationAddress() string {
	return c.cfg.RegAddr()
}

// AuthorizationAddress возвращает адрес для метода авторизации
func (c Client) AuthorizationAddress() string {
	return c.cfg.AuthAddr()
}

// ActionAddr возвращает адрес для методов действий(добавление, удаление...)
func (c Client) ActionAddr() string {
	return c.cfg.ActionAddr()
}

// SetLocal производит cast интерфейса в необходимый тип и выполняет метод хранилища с преобразованной структурой
func (c Client) SetLocal(request models.Validatable) error {
	r, ok := request.(models.UserData)
	if !ok {
		return models.ErrUncastable
	}

	return c.store.Set(r)
}

// UpdateLocal производит cast интерфейса в необходимый тип и выполняет метод хранилища с преобразованной структурой
func (c Client) UpdateLocal(request models.Validatable) error {
	r, ok := request.(models.UserData)
	if !ok {
		return models.ErrUncastable
	}

	return c.store.Update(r)
}

// DeleteLocal производит cast интерфейса в необходимый тип и выполняет метод хранилища с преобразованной структурой
func (c Client) DeleteLocal(request models.Validatable) error {
	r, ok := request.(models.DeleteRequest)
	if !ok {
		return models.ErrUncastable
	}

	return c.store.Delete(r)
}

// GetAll обертка над методом хранилища для получения полного списка данных
func (c Client) GetAll() ([]models.UserData, error) {
	return c.store.GetAll()
}

// ActualizeStorage вызывает метод получения актуальных данных с сервера и производит обновление хранилища
func (c Client) ActualizeStorage() error {
	data, err := c.GetActualData()
	if err != nil {
		return err
	}

	for _, el := range data {
		err = c.store.Set(el)
		if err != nil {
			return err
		}
	}

	return nil
}
