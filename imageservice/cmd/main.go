package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"
	"image/api"
	"image/configs"
	imageHandlers "image/internal/app/image/handlers"
	imageUCase "image/internal/app/image/usecase"
	"image/internal/app/logger"
	"image/internal/app/middleware"
	"log"
	"net/http"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/toml/server.toml", "path to config file")
}

func main() {
	config := configs.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := logger.InitLogger("stdout"); err != nil {
		log.Fatal(err)
	}

	// grpc connect to UserService
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: "ImageService",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831",
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)

	if err != nil {
		log.Fatal("cannot create tracer", err)
	}

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	userRepo := api.NewUserClient(conn)

	imageUseCase := imageUCase.NewUseCase(userRepo)

	imageHandler := imageHandlers.NewHandler(*imageUseCase)

	csrfMiddleware := middleware.CSRFMiddleware(config.HTTPS)

	router := mux.NewRouter()
	router.Use(middleware.LoggingRequest)
	router.Use(csrfMiddleware)

	router.HandleFunc("/api/profile/avatar", imageHandler.PutAvatar).Methods(http.MethodPatch)

	c := middleware.CorsMiddleware(config.Origin)

	server := &http.Server{
		Addr:    config.BindAddr,
		Handler: c.Handler(router),
	}

	if config.HTTPS {
		log.Println("TLS server starting at port: ", server.Addr)
		if err := server.ListenAndServeTLS(
			"/etc/letsencrypt/live/findfreelancer.ru/cert.pem",
			"/etc/letsencrypt/live/findfreelancer.ru/privkey.pem"); err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Server starting at port", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
