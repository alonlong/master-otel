package apid

import (
	"fmt"
	"master-otel/internal/proto/common/v1"
	"master-otel/pkg/log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/opentelemetry-go-extra/otelplay"
	"go.uber.org/zap"
)

type CreateUserPOST struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (s *Service) createUser(c echo.Context) error {
	ctx := c.Request().Context()

	log.WithContext(ctx).Info("create user", zap.String("remote-addr", c.Request().RemoteAddr))
	otelplay.PrintTraceID(ctx)
	// log.WithContext(c.Request().Context()).Info("create user", zap.String("remote-addr", c.Request().RemoteAddr))

	var req CreateUserPOST
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("bind user: %w", err))
	}
	user, err := s.ctldClient.CreateUser(c.Request().Context(), &common.User{
		Email:    req.Email,
		Username: req.Username,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("create user: %v", err))
	}
	return c.JSON(http.StatusOK, user)
}
