package cmd

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"user/internal/app/user/handlers"
	"user/pkg/midlewhare"

	"log"
	"net/http"
	"user/configs"
	specializeRepository "user/internal/app/specialize/repository/postgresql"
	userRepository  "user/internal/app/user/repository/postgresql"
	"user/internal/app/user/usecase/impl"
	"user/pkg/database/postgresql"
	"user/pkg/logger"
)
var (
	configPath string
)
func init() {
	flag.StringVar(&configPath, "config-path", "config/server.toml", "path to config file")
}

func main() {
	flag.Parse()
	config := configs.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err = logger.InitLogger("stdout"); err != nil {
		log.Fatal(err)
	}


	postgres, err := postgresql.NewPostgres(config.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := postgres.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	userRepository := &userRepository.Repository{
		Db: postgres.GetPostgres(),
	}
	specializeRepository :=  &specializeRepository.Repository{
		Db: postgres.GetPostgres(),
	}

	userUseCase := impl.New(userRepository, specializeRepository)

	userHandler := handlers.New(userUseCase)
	router := mux.NewRouter()
	router.Use(midlewhare.LoggingRequest)
	router.HandleFunc("", userHandler.GetUserInfo).Methods(http.MethodGet)
	router.HandleFunc("", userHandler.ChangeProfile).Methods(http.MethodPut)


	c := midlewhare.CorsMiddleware(config.Origin)
	server := &http.Server{
		Addr:    config.BindAddr,
		Handler: c.Handler(router),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
