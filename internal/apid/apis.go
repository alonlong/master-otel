package apid

import (
	"fmt"
	"master-otel/internal/proto/common/v1"
	"master-otel/pkg/log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateUserPOST struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (s *Service) createUser(c echo.Context) error {
	ctx := c.Request().Context()
	log.WithCtx(ctx).Info("create user")
	var req CreateUserPOST
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("bind user: %w", err))
	}
	user, err := s.ctldClient.CreateUser(ctx, &common.User{
		Email:    req.Email,
		Username: req.Username,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("create user: %v", err))
	}
	return c.JSON(http.StatusOK, user)
}
