package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	urlcnt "rest_url_shortener/internal/app/v1/url"
	pb "rest_url_shortener/internal/pb"
	rep "rest_url_shortener/internal/repository"
	"rest_url_shortener/internal/repository/cache"
	pgdb "rest_url_shortener/internal/repository/postgres"
	"rest_url_shortener/internal/service/url"
	"rest_url_shortener/internal/utils"
	"time"
)

func main() {
	var database string
	// Проверка флага для выбора хранилища
	flag.StringVar(&database, "database", "in-memory", "select database: postgres | in-memory")
	flag.Parse()
	log.Printf("✓	database selected: %s", database)
	// Загрузка конфигурационного файла .env
	cfg, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("✘	cannot load config:", err)
	}
	log.Print("✓	configuration successfully set")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var repo rep.Repository
	/*
		Создание хранилища в зависимости от выбранного флага (значение по умолчанию in-memory)
		Но если будет передан флаг нереализованного хранилища, то сервис остановится
	*/
	switch database {
	case "postgres":
		// Создание подключения к базе данных

		pool, err := utils.GetConnectToPg(ctx, &cfg)
		if err != nil {
			log.Fatalf("✘	unable to connection to database: %v\n", err)
		}
		// Создание репозитория Postgres
		repo = pgdb.NewRepository(pool)
	case "in-memory":
		// Создание репозитория In-Memory
		repo = cache.NewInMemoryRepository()
	default:
		log.Fatalf("✘	non-existent repository specified")
	}
	log.Print("✓	Successful database connection")
	// Используем конструктор для создания сервиса и передаём выбранный репозиторий
	service := url.NewService(repo)
	// Используем конструктор для создания контроллера
	controller := urlcnt.NewUrlController(service)
	// Создание Listener для прослушивания нужных портов gRPC и Proxy (gRPC gateway).
	lisGrpc, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.PortGrpc))
	if err != nil {
		log.Fatalf("✘	failed to listen: %v ", err)
	}
	lisHttp, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.PortHttp))
	if err != nil {
		log.Fatalf("✘	failed to listen: %v", err)
	}
	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() (err error) {
		return runGRPCServer(lisGrpc, controller)
	})
	log.Printf("✓	start gRPC server at %s", lisGrpc.Addr().String())
	g.Go(func() (err error) {
		return runGatewayServer(lisHttp, controller)
	})
	log.Printf("✓	start HTTP (gRPC gateway) server at %s", lisHttp.Addr().String())

	err = g.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func runGRPCServer(lis net.Listener, server *urlcnt.UrlController) error {
	// Создание нового gRPC сервера
	grpcServer := grpc.NewServer()
	// Регистрация сервера и реализованных эндпоинтов (связь с сгенерированным proto)
	pb.RegisterUrlShortenerServer(grpcServer, server)
	reflection.Register(grpcServer)
	return grpcServer.Serve(lis)
}

func runGatewayServer(lis net.Listener, server *urlcnt.UrlController) error {
	/*
		Создаём ServeMux, который будет сопоставлять HTTP запросы с шаблонами и вызывать
		соответствующие обработчики
	*/
	mux := runtime.NewServeMux()
	// Регистрация HTTP обработчика для сервиса
	err := pb.RegisterUrlShortenerHandlerServer(context.Background(), mux, server)
	if err != nil {
		log.Fatal(err)
	}
	// Добавлении нового обработчкика из файла
	fs := http.FileServer(http.Dir("./doc/swagger"))
	// Создание нового ServeMux, но уже HTTP и прописывание роутов
	restMux := http.NewServeMux()
	restMux.Handle("/", mux)
	restMux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	return http.Serve(lis, restMux)
}
