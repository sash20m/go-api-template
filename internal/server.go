package internal

import (
	"context"
	"errors"
	"fmt"
	"go-api-template/config"
	"go-api-template/internal/libs/database"
	"go-api-template/internal/libs/queue"
	"go-api-template/internal/libs/renderer"
	"go-api-template/internal/service"
	httpTransport "go-api-template/internal/transport/http"
	queueTransport "go-api-template/internal/transport/queue"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
)

type Server struct {
	HTTP             *httpTransport.HTTPTransport
	Queue            *queueTransport.QueueTransport
	ResponseRenderer *renderer.ResponseRenderer

	Services *service.Services

	// Used for migrations and connection closing on shutdown
	PostgresDB *database.PostgresDB
	RabbitMQ   *queue.RabbitMQ
}

func NewServer() *Server {
	postgresDB, err := database.NewPostgresDB()
	if err != nil {
		panic(err)
	}

	var rabbit *queue.RabbitMQ
	var publisher queue.Publisher = queue.NoopPublisher{}
	if config.CONFIG.RabbitMQEnabled {
		rabbit = queue.NewRabbitMQ(config.CONFIG.RabbitMQURL, queue.RabbitMQOptions{
			Prefetch: config.CONFIG.RabbitMQPrefetch,
		})
		if err := rabbit.Connect(); err != nil {
			panic(err)
		}
		if err := rabbit.EnsureTopology(queue.EventsExchangeName, queue.OwnedQueues()); err != nil {
			panic(err)
		}
		publisher = rabbit
	}

	services := service.NewServices(postgresDB.Database, publisher)

	responseRenderer := renderer.NewResponseRenderer()

	httpTransport := httpTransport.NewHTTPTransport(services, responseRenderer)
	queueTransport := queueTransport.NewQueueTransport(services, rabbit)

	return &Server{
		Services:         services,
		HTTP:             httpTransport,
		Queue:            queueTransport,
		ResponseRenderer: responseRenderer,
		PostgresDB:       postgresDB,
		RabbitMQ:         rabbit,
	}
}

func (s *Server) Run(ctx context.Context) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.CONFIG.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	secureMiddleware := secure.New(secure.Options{
		IsDevelopment:      config.CONFIG.Env == config.DEV_ENV,
		ContentTypeNosniff: true,
		// SSLRedirect:        config.CONFIG.Env == config.PROD_ENV,
		// When the API is behind nginx
		// SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
	r.Use(secureMiddleware.Handler)

	r.NotFound(s.NotFoundHandler)
	r.MethodNotAllowed(s.NotAllowedHandler)

	s.RegisterRoutes(r)

	server := http.Server{
		Addr:         ":" + config.CONFIG.Port,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 900 * time.Second,
		IdleTimeout:  1200 * time.Second,
		Handler:      r,
	}

	errCh := make(chan error, 1)

	// Start HTTP server
	go func() {
		startupMessage := "Starting HTTP server (v" + config.CONFIG.Version + ")"
		startupMessage = startupMessage + " on port " + config.CONFIG.Port
		startupMessage = startupMessage + " in " + string(config.CONFIG.Env) + " mode."
		logrus.Info(startupMessage)

		logrus.Info("HTTP Server Listening...")
		errCh <- server.ListenAndServe()

	}()

	// Start RabbitMQ listener
	if config.CONFIG.RabbitMQEnabled && s.Queue != nil {
		go func() {
			if err := s.Queue.StartConsumers(ctx); err != nil && !errors.Is(err, context.Canceled) {
				errCh <- err
			}
		}()
	}

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		server.Shutdown(shutdownCtx)

		s.OnShutdown()
		return nil

	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			s.OnShutdown()
			return nil
		}
		logrus.WithError(err).Error("Server stopped with error")
		s.OnShutdown()
		return err
	}
}

func (s *Server) RegisterRoutes(r chi.Router) {
	// All top level routes should be registered here.
	r.Route("/api", func(r chi.Router) {
		r.Route("/users", s.HTTP.Users.RegisterRoutes)
	})
}

// OnShutdown is called when the server has a panic.
// It is used to cleanup the server resources.
func (s *Server) OnShutdown() {
	s.PostgresDB.Close()
	logrus.Info("PostgresDB closed")

	if s.RabbitMQ != nil {
		_ = s.RabbitMQ.Close()
		logrus.Info("RabbitMQ closed")
	}

	logrus.WithError(fmt.Errorf("Executed OnShutdown")).Error("Executed OnShutdown")
}

func (s *Server) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	s.ResponseRenderer.JSON(w, http.StatusNotFound, fmt.Sprint("Not Found ", r.URL))
}

func (s *Server) NotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	s.ResponseRenderer.JSON(w, http.StatusMethodNotAllowed, fmt.Sprint(r.Method, " method not allowed"))

}
