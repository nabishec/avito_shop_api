package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nabishec/avito_shop_api/cmd/db_connection"
	"github.com/nabishec/avito_shop_api/internal/http_server/handlers/auth"
	"github.com/nabishec/avito_shop_api/internal/http_server/handlers/buy"
	"github.com/nabishec/avito_shop_api/internal/http_server/handlers/info"
	"github.com/nabishec/avito_shop_api/internal/http_server/handlers/send"
	"github.com/nabishec/avito_shop_api/internal/http_server/middlweare"
	"github.com/nabishec/avito_shop_api/internal/pkg"
	"github.com/nabishec/avito_shop_api/internal/storage/db"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/nabishec/avito_shop_api/docs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title API Avito shop
// @version 1.0.0
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	//TODO: init logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("d", false, "set log level to debug")
	easyReading := flag.Bool("r", false, "set console writer")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	//for easy reading
	if *easyReading {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	//TODO: init config
	err := LoadEnv()
	if err != nil {
		log.Error().Err(err).Msg("don't found configuration")
		os.Exit(1)
	}

	//TODO: init storage postgresql
	log.Info().Msg("Init storage")
	dbConnection, err := db_connection.NewDatabaseConnection()
	if err != nil {
		log.Error().AnErr(pkg.ErrReader(err)).Msg("Failed init database")
		os.Exit(1)
	}
	log.Info().Msg("Database init successful")

	//TODO: init middlewear
	s := CreateNewServer(dbConnection.DB)

	s.MountHandlers()
	//TODO: run server
	wrTime, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Error().Err(err).Msg("timeout not received from env")
		wrTime = 4 * time.Second
	}
	idleTime, err := time.ParseDuration(os.Getenv("IDLE_TIMEOUT"))
	if err != nil {
		log.Error().Err(err).Msg("idle timeout not received from env")
		idleTime = 60 * time.Second // CHECK THIS
	}

	srv := &http.Server{
		Addr:         ":" + os.Getenv("SERVER_PORT"),
		Handler:      s.Router,
		ReadTimeout:  wrTime,
		WriteTimeout: wrTime,
		IdleTimeout:  idleTime,
	}
	log.Info().Msgf("Starting server on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Error().Msg("failed to start server")
		os.Exit(1)
	}

	log.Error().Msg("Program ended")
}

type Server struct {
	Router  *chi.Mux
	Storage *db.Storage
}

func CreateNewServer(dbConnection *sqlx.DB) *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	s.Storage = db.NewStorage(dbConnection)
	return s
}

func (s *Server) MountHandlers() {

	sendCoin := send.NewSendingCoins(s.Storage)
	getInformation := info.NewUserInformation(s.Storage)
	buyItem := buy.NewBuying(s.Storage)
	authentication := auth.NewAuth(s.Storage)

	s.Router.Group(func(r chi.Router) {
		r.Get("/swagger/*", httpSwagger.WrapHandler)
		r.Post("/api/auth", authentication.ReturnAuthToken)
	})

	// Require Authentication
	s.Router.Group(func(r chi.Router) {
		r.Use(middlweare.Auth)
		r.Get("/api/buy/{item}", buyItem.BuyingItemByUser)
		r.Post("/api/sendCoin", sendCoin.SendCoins)
		r.Get("/api/info", getInformation.ReturnUserInfo)
	})

}

func LoadEnv() error {
	const op = "cmd.loadEnv()"
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("%s:%s", op, "failed load env file")
	}
	return nil
}
