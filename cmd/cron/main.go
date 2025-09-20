package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fikryfahrezy/let-it-go/config"
	"github.com/fikryfahrezy/let-it-go/pkg/database"
	"github.com/fikryfahrezy/let-it-go/pkg/http_server"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/go-co-op/gocron/v2"

	blogRepository "github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	blogService "github.com/fikryfahrezy/let-it-go/feature/blog/service"
	userRepository "github.com/fikryfahrezy/let-it-go/feature/user/repository"
	userService "github.com/fikryfahrezy/let-it-go/feature/user/service"
)

type Job struct {
	name    string
	crontab string
	task    func()
}

func sampleTask(log *slog.Logger, userSrv userService.UserService, blogSrv blogService.BlogService) func() {
	return func() {
		ctx := context.Background()
		log.Info("Running sample task")
		// Example task: You can replace this with actual logic
		users, userCount, err := userSrv.ListUsers(ctx, userService.ListUsersRequest{
			PaginationRequest: http_server.PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
		})
		if err != nil {
			log.Error("Failed to fetch users",
				slog.String("error", err.Error()),
			)
			return
		}
		log.Info("Fetched users", slog.Int("count", len(users)), slog.Int64("total", userCount))

		blogs, blogCount, err := blogSrv.ListBlogs(ctx, blogService.ListBlogsRequest{
			PaginationRequest: http_server.PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
		})
		if err != nil {
			log.Error("Failed to fetch blogs",
				slog.String("error", err.Error()),
			)
			return
		}
		log.Info("Fetched blogs", slog.Int("count", len(blogs)), slog.Int64("total", blogCount))
	}
}

func main() {
	cfg := config.Load()

	log := logger.NewLogger(cfg.Logger)
	log.Info("Starting cron")

	db, err := database.NewDB(cfg.Database)
	if err != nil {
		log.Error("Failed to connect to database",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Error("Failed to create cron scheduler",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	userRepo := userRepository.NewUserRepository(log, db)
	userService := userService.NewUserService(log, userRepo)

	blogRepo := blogRepository.NewBlogRepository(log, db)
	blogService := blogService.NewBlogService(log, blogRepo)

	jobs := []Job{
		{
			name:    "sample_task",
			crontab: cfg.Crontab["sample_task"],
			task:    sampleTask(log, userService, blogService),
		},
	}

	for _, job := range jobs {
		_, err := s.NewJob(
			gocron.CronJob(
				job.crontab,
				false,
			),
			gocron.NewTask(
				job.task,
			),
		)
		if err != nil {
			log.Error("Failed to create cron job "+job.name,
				slog.String("error", err.Error()),
			)
			os.Exit(1)
		}
		log.Info("Scheduled job: " + job.name + " with crontab: " + job.crontab)
	}

	// start the scheduler
	s.Start()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down gracefully...")

	if err := s.Shutdown(); err != nil {
		log.Error("Failed to shutdown cron scheduler",
			slog.String("error", err.Error()),
		)
	}

	log.Info("Cron shutdown complete")
}
