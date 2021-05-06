package main

import (
	"authorizationservice/api"
	"authorizationservice/configs"
	handlers2 "authorizationservice/internal/app/profile/handlers"
	repository3 "authorizationservice/internal/app/profile/repository"
	usecase3 "authorizationservice/internal/app/profile/usecase"
	handlers3 "authorizationservice/internal/app/session/handlers"
	"authorizationservice/internal/app/session/repository"
	"authorizationservice/internal/app/session/usecase"
	"authorizationservice/pkg/logger"
	"authorizationservice/pkg/metric"
	"authorizationservice/pkg/middleware"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	traceutils "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/tarantool/go-tarantool"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/toml/server.toml", "path to config file")
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

	opts := tarantool.Opts{
		User: "guest",
	}
	conn, err := tarantool.Connect(config.DatabaseURL, opts)
	defer conn.Close()
	sessionRepository := repository.New(conn)

	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: "AuthService",
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
	grpcConn, err := grpc.Dial(
		":8081",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(traceutils.OpenTracingClientInterceptor(tracer)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer grpcConn.Close()
	client := api.NewUserClient(grpcConn)
	profileRepository := repository3.New(client)

	sessionUseCase := usecase.New(sessionRepository)
	profileUseCase := usecase3.New(profileRepository)

	sessionHandler := handlers3.New(sessionUseCase)
	profileHandler := handlers2.New(sessionUseCase, profileUseCase)

	csrfMiddleware := middleware.CSRFMiddleware(config.HTTPS)

	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/metrics").Handler(promhttp.Handler())

	apiRout := router.PathPrefix("/api").Subrouter()
	apiRout.Use(middleware.LoggingRequest)
	apiRout.HandleFunc("/logout", sessionHandler.LogOut).Methods(http.MethodDelete)
	apiRout.HandleFunc("/profile", profileHandler.RegistrationProfile).Methods(http.MethodPost)
	apiRout.HandleFunc("/login", profileHandler.AuthorisationProfile).Methods(http.MethodPost)

	profile := apiRout.PathPrefix("/profile").Subrouter()
	profile.Use(csrfMiddleware)
	profile.HandleFunc("/authorized", sessionHandler.CheckLogin).Methods(http.MethodGet)

	c := middleware.CorsMiddleware(config.Origin)
	server := &http.Server{
		Addr:    config.BindAddr,
		Handler: c.Handler(router),
	}

	go func() {
		if config.HTTPS {
			log.Println("TLS server starting at port: ", server.Addr)
			if err := server.ListenAndServeTLS(
				"/etc/letsencrypt/live/findfreelancer.ru/cert.pem",
				"/etc/letsencrypt/live/findfreelancer.ru/privkey.pem"); err != nil {
				log.Fatal(err)
			}
		}
		log.Println("HTTP server start on port", config.BindAddr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	metric.New()
	s := grpc.NewServer(grpc.UnaryInterceptor(traceutils.OpenTracingServerInterceptor(tracer)))
	srv := handlers3.NewGRPCServer(sessionUseCase)
	api.RegisterSessionServer(s, srv)
	l, err := net.Listen("tcp", ":8085")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("GRPC server start on port :8085")
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
