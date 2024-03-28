package server

import (
	"github.com/Edouard127/lambda-rpc/pkg/setting"
	"github.com/Edouard127/lambda-rpc/router"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
)

type Server struct {
	log      *slog.Logger
	engine   *gin.Engine
	settings *setting.Server
}

func New(settings *setting.Server) *Server {
	if settings == nil {
		settings = setting.DefaultServer
	}

	return &Server{
		log:      slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		engine:   router.CreateEngine(),
		settings: settings,
	}
}

func (s *Server) Start() {
	s.log.Info("Starting server")
	s.engine.Run(":8080") // Will panic if the server can't start, no need to handle the error
}
