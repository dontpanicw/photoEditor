package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	http2 "homework-dontpanicw/app/api/http"
	_ "homework-dontpanicw/app/docs"
	"homework-dontpanicw/app/pkg/db"
	pkgHttp "homework-dontpanicw/app/pkg/http"
	"homework-dontpanicw/app/pkg/swagger"
	"homework-dontpanicw/app/repository/postgres"
	"homework-dontpanicw/app/repository/rabbitmq"
	"homework-dontpanicw/app/repository/redis"
	"homework-dontpanicw/app/usecases/auth"
	service2 "homework-dontpanicw/app/usecases/service"
	"log"
)

// @title My API
// @version 1.0
// @description This is a sample server.

// @host localhost:8080
// @BasePath /
func main() {

	addr := flag.String("addr", ":8080", "server service address")

	flag.Parse()

	//SESSION
	SessionRepo, err := redis.NewSessionCashStorage("redis:6379", "", 0)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	//USERS
	UserRepo, err := postgres.NewUserPostgresStorage("postgres://postgres:postgres@db:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	if err := db.EnsureTasksTable(UserRepo.GetDb()); err != nil {
		log.Fatalf("failed to ensure tasks table: %v", err)
	}
	UserService := service2.NewUser(UserRepo, SessionRepo)
	AuthHandlers := http2.NewUserHandler(UserService)

	//TASKS
	TaskRepo, err := postgres.NewTaskPostgresStorage("postgres://postgres:postgres@db:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	if err := db.EnsureTasksTable(TaskRepo.GetDb()); err != nil {
		log.Fatalf("failed to ensure tasks table: %v", err)
	}

	TaskSender, err := rabbitmq.NewRabbitMQSender("amqp://guest:guest@broker:5672/", "photos", "idResponses")
	if err != nil {
		log.Fatal("failed creating rabbitmq, %w", err)
	}

	TaskService := service2.NewTask(TaskRepo, TaskSender)
	TaskHandlers := http2.NewTaskHandler(TaskService)

	go TaskSender.ListenForResponses(TaskService)

	r := chi.NewRouter()
	swagger.CreateSwaggerRouter(r)

	authMiddleware := auth.NewAuthMiddleware(SessionRepo)

	TaskHandlers.WithTaskHandlers(r, authMiddleware)

	AuthHandlers.WithAuthHandlers(r)

	log.Printf("Listening on %s", *addr)
	if err := pkgHttp.CreateAndRunServer(r, *addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
