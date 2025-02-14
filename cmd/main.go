package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/nabishec/avito_shop_api/internal/http_server/handlers/auth"
	"github.com/nabishec/avito_shop_api/internal/http_server/handlers/shop"
	"github.com/nabishec/avito_shop_api/internal/http_server/middlweare"
	"github.com/nabishec/avito_shop_api/internal/pkg"
	"github.com/nabishec/avito_shop_api/internal/storage/db"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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
	err := loadEnv()
	if err != nil {
		log.Error().Err(err).Msg("don't found configuration")
		os.Exit(1)
	}

	//TODO: init storage postgresql
	log.Info().Msg("Init storage")
	storage, err := db.NewDatabase()
	if err != nil {
		log.Error().AnErr(pkg.ErrReader(err)).Msg("Failed init storage")
		os.Exit(1)
	}
	log.Info().Msg("Storage init successful")

	//TODO: init middlewear
	router := chi.NewRouter()

	//TODO: getInformation :=
	//TODO: sendCoin :=
	buyItem := shop.NewBuyItem(storage)
	authentication := auth.NewAuth(storage)

	router.Group(func(r chi.Router) {
		router.Post("/api/auth", authentication.ReturnAuthToken)
	})

	// Private Routes
	// Require Authentication
	router.Group(func(r chi.Router) {
		r.Use(middlweare.Auth)
		router.Get("/api/buy", buyItem.BuyingItemByUser)
	})

	//TODO: run server
	wrTime, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Error().Err(err).Msg("timeout not received from env")
		wrTime = 4 * time.Second
	}
	idleTime, err := time.ParseDuration(os.Getenv("IDLE_TIMEOUT"))
	if err != nil {
		log.Error().Err(err).Msg("idle timeout not received from env")
		idleTime = 60 * time.Second //CHECK THIS
	}

	srv := &http.Server{
		Addr:         os.Getenv("ADDRESS"),
		Handler:      router,
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

func loadEnv() error {
	const op = "cmd.loadEnv()"
	err := godotenv.Load("./configs/configuration.env")
	if err != nil {
		return fmt.Errorf("%s:%s", op, "failed load env file")
	}
	return nil
}
