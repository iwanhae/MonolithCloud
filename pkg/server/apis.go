package server

import "github.com/labstack/echo/v4"

type Server struct{}

var _ ServerInterface = &Server{}

// CreateVM implements ServerInterface
func (*Server) CreateVM(ctx echo.Context) error {
	panic("unimplemented")
}

// GetVM implements ServerInterface
func (*Server) GetVM(ctx echo.Context, vmId int) error {
	panic("unimplemented")
}

// ListVMs implements ServerInterface
func (*Server) ListVMs(ctx echo.Context) error {
	panic("unimplemented")
}
