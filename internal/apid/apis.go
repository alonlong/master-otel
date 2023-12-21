package apid

import (
	"github.com/labstack/echo/v4"
)

func httpError(c echo.Context, code int, err error) error {
	if err := c.JSON(code, err.Error()); err != nil {
		return err
	}
	return err
}

type CreateUserPOST struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Service) createUser(c echo.Context) error {
	return nil
}

func (s *Service) deleteUser(c echo.Context) error {
	return nil
}
