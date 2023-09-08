package server

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/satvikprasad/vikingx/db"
	"github.com/satvikprasad/vikingx/trader"
)

type Context struct {
	Trader   trader.Trader
	Database db.Database
	Context  *gin.Context
}

type Server struct {
	t  trader.Trader
	db db.Database
	r  *gin.Engine

	port string
}

func CreateServer(db db.Database, t trader.Trader, port string) *Server {
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	r.Use(cors.New(config))
	r.Static("/js/", "./templates/js/")
	r.Static("/assets/", "./templates/assets/")

	return &Server{
		r:    r,
		t:    t,
		db:   db,
		port: port,
	}
}

func (s *Server) GET(route string, fn ApiFunc) {
	s.r.GET(route, makeAPIFunc(s.t, s.db, fn))
}

func (s *Server) POST(route string, fn ApiFunc) {
	s.r.POST(route, makeAPIFunc(s.t, s.db, fn))
}

func (s *Server) Listen() {
	s.r.Run(":" + s.port)
}

func makeAPIFunc(t trader.Trader, db db.Database, fn ApiFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if db == nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not initialise database"})
			return
		}

		ctx := &Context{
			Trader:   t,
			Database: db,
			Context:  c,
		}

		if err := fn(ctx); err != nil {
			fmt.Println(err)
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}

type ApiFunc func(c *Context) error
