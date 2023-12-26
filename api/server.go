package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/indramhrt/simplebank/db/sqlc"
	"github.com/indramhrt/simplebank/token"
	"github.com/indramhrt/simplebank/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config *util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSysmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{config: *config, store: store, tokenMaker: tokenMaker}
	server.setupRouter()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/account", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/account", server.listAccounts)
	router.DELETE("/account", server.deleteAccount)

	router.POST("/user", server.createUser)
	router.POST("/user/login", server.loginUser)

	router.POST("/transfer", server.transferMoney)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
