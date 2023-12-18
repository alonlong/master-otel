package apid

import (
	"fmt"
	"net/http"
	"strconv"

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

func (s *Service) deleteUser(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return httpError(c, http.StatusBadRequest, fmt.Errorf("parse id: %w", err))
	}
	if _, err := s.ctldClient.DeleteUser(c.Request().Context(), &commonv1.Identity{
		Id: id,
	}); err != nil {
		return httpError(c, http.StatusInternalServerError, fmt.Errorf("delete user: %w", err))
	}
	return c.JSON(http.StatusOK, "success")
}
