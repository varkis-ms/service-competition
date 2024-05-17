package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	grpcapp "github.com/varkis-ms/service-competition/internal/app/grpc"
	"github.com/varkis-ms/service-competition/internal/config"
	"github.com/varkis-ms/service-competition/internal/pkg/database/postgresdb"
	"github.com/varkis-ms/service-competition/internal/pkg/logger/handlers/slogpretty"
	"github.com/varkis-ms/service-competition/internal/pkg/logger/sl"
	"github.com/varkis-ms/service-competition/internal/rpc/competition_list"
	"github.com/varkis-ms/service-competition/internal/rpc/create_competition"
	"github.com/varkis-ms/service-competition/internal/rpc/edit_competition"
	"github.com/varkis-ms/service-competition/internal/rpc/get_competition_info"
	"github.com/varkis-ms/service-competition/internal/rpc/get_leaderboard"
	"github.com/varkis-ms/service-competition/internal/rpc/get_next_solution"
	"github.com/varkis-ms/service-competition/internal/rpc/save_solution"
	"github.com/varkis-ms/service-competition/internal/rpc/save_solution_result"
	"github.com/varkis-ms/service-competition/internal/rpc/user_activity_full"
	"github.com/varkis-ms/service-competition/internal/rpc/user_activity_total"
	"github.com/varkis-ms/service-competition/internal/storage"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type App struct {
	GRPCServer *grpcapp.App
}

func Run(configPath string) {
	// Config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		panic("config.LoadConfig failed" + err.Error())
	}

	logger := setupLogger(cfg.Env)

	// Repository
	logger.Info("Initializing postgres...")
	db, err := postgresdb.New(&cfg)
	if err != nil {
		logger.Error("postgresdb.New failed", sl.Err(err))
	}
	defer db.Close()
	repositories := storage.New(db)

	// Handlers
	createCompetitionHandler := create_competition.New(repositories, logger)
	competitionListHandler := competition_list.New(repositories, logger)
	getCompetitionInfoHandler := get_competition_info.New(repositories, logger)
	getLeaderboardHandler := get_leaderboard.New(repositories, logger)
	competitionEditHandler := edit_competition.New(repositories, logger)
	userActivityFullHandler := user_activity_full.New(repositories, logger)
	userActivityTotalHandler := user_activity_total.New(repositories, logger)
	saveSolutionResultHandler := save_solution_result.New(repositories, logger)
	getNextSolution := get_next_solution.New(repositories, logger)
	saveSolution := save_solution.New(repositories, logger)

	// gRPC server
	app := grpcapp.New(
		logger,
		cfg.Port,
		createCompetitionHandler,
		competitionEditHandler,
		competitionListHandler,
		getCompetitionInfoHandler,
		getLeaderboardHandler,
		userActivityFullHandler,
		userActivityTotalHandler,
		saveSolutionResultHandler,
		getNextSolution,
		saveSolution,
	)

	go func() {
		app.MustRun()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.Shutdown()
	logger.Info("Gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
