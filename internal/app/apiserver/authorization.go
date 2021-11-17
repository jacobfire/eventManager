package apiserver

import (
	"calendar/internal/app/model"
	"errors"
)

func Authorize(c model.Credentials) (*model.User, error) {
	if c.Username == "jack@example.com" && c.Password== "Qwerty" {
		User := &model.User {
			Username: c.Username,
			Password: c.Password,
		}

		return User, nil
	}

	return nil, errors.New("not authorized user")
}
