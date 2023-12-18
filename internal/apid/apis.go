package apid

import (
	"fmt"
	"net/http"

	commonv1 "master-otel/internal/proto/common/v1"

	"github.com/labstack/echo/v4"
)

func httpError(c echo.Context, code int, err error) error {
	if err := c.JSON(code, err.Error()); err != nil {
		return err
	}
	return err
}

type CreateUserPOST struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (s *Service) createUser(c echo.Context) error {
	var req CreateUserPOST
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("bind user: %w", err))
	}
	user, err := s.ctldClient.CreateUser(c.Request().Context(), &commonv1.User{
		Email:    req.Email,
		Username: req.Username,
	})
	if err != nil {
		return httpError(c, http.StatusInternalServerError, fmt.Errorf("create user: %w", err))
	}
	return c.JSON(http.StatusOK, user)
}
