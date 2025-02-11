package main

import (
	"Commentary/app"
	"Commentary/internal/config"
	"Commentary/internal/graph"
	"Commentary/internal/logger"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/ast"
	"log"
	"net/http"
	"strconv"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		logrus.Fatal(err)
	}
	logger.InitLogger(cfg)
	appObj := app.InitApp(cfg)

	port := strconv.Itoa(cfg.Server.Port)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: appObj.Resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	srv.AddTransport(transport.Websocket{
		Upgrader: upgrader,
	})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
