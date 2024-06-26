package router

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	gorillaWs "github.com/gorilla/websocket"

	"github.com/semanser/ai-coder/graph"
	"github.com/semanser/ai-coder/websocket"
)

func New(db *gorm.DB) *gin.Engine {
	// Initialize Gin router
	r := gin.Default()

	// Configure CORS middleware
	config := cors.DefaultConfig()
	// TODO change to only allow specific origins
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))

	// GraphQL endpoint
	r.Any("/graphql", graphqlHandler(db))

	// GraphQL playground route
	r.GET("/playground", playgroundHandler())

	// WebSocket endpoint for Docker daemon
	r.GET("/terminal/:id", wsHandler())

	return r
}

func graphqlHandler(db *gorm.DB) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Db: db,
	}}))

	h.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		res := next(ctx)
		if res == nil {
			return res
		}

		err := res.Errors.Error()

		if err != "" {
			log.Printf("graphql error: %s", err)
		}

		return res
	})

	// We can't use the default error handler because it doesn't work with websockets
	// https://stackoverflow.com/a/75444816
	// So we add all the transports manually (see handler.NewDefaultServer in gqlgen for reference)
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})
	h.AddTransport(transport.MultipartForm{})

	h.SetQueryCache(lru.New(1000))

	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	// Add transport to support GraphQL subscriptions
	h.AddTransport(&transport.Websocket{
		Upgrader: gorillaWs.Upgrader{
			CheckOrigin: func(r *http.Request) (allowed bool) {
				// TODO change to only allow specific origins
				return true
			},
		},
		InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, *transport.InitPayload, error) {
			return ctx, &initPayload, nil
		},
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		playground.Handler("GraphQL", "/graphql").ServeHTTP(c.Writer, c.Request)
	}
}

func wsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		websocket.HandleWebsocket(c)
	}
}
