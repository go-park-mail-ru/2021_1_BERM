package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	traceutils "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"
	"imageservice/api"
	"imageservice/configs"
	imageHandlers "imageservice/internal/app/image/handlers"
	imageUCase "imageservice/internal/app/image/usecase"
	"imageservice/internal/app/imgcreator"
	"imageservice/internal/app/logger"
	"imageservice/internal/app/metric"
	"imageservice/internal/app/middleware"
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

	// grpc connect to UserService
	conn, err := grpc.Dial(
		config.Host+":8081",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(traceutils.OpenTracingClientInterceptor(tracer)))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	userRepo := api.NewUserClient(conn)

	imgCr := imgcreator.ImgCreator{}
	imageUseCase := imageUCase.NewUseCase(userRepo, &imgCr)

	imageHandler := imageHandlers.NewHandler(imageUseCase)

	csrfMiddleware := middleware.CSRFMiddleware(config.HTTPS)

	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/metrics").Handler(promhttp.Handler())

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.LoggingRequest)
	apiRouter.Use(csrfMiddleware)
	apiRouter.HandleFunc("/profile/avatar", imageHandler.PutAvatar).Methods(http.MethodPatch)
	metric.New()
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
