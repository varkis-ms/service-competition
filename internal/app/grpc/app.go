package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/varkis-ms/service-competition/internal/rpc/competition_list"
	"github.com/varkis-ms/service-competition/internal/rpc/edit_competition"
	"github.com/varkis-ms/service-competition/internal/rpc/get_competition_info"
	"github.com/varkis-ms/service-competition/internal/rpc/get_leaderboard"
	"github.com/varkis-ms/service-competition/internal/rpc/get_next_solution"
	"github.com/varkis-ms/service-competition/internal/rpc/save_solution"
	"github.com/varkis-ms/service-competition/internal/rpc/save_solution_result"
	"github.com/varkis-ms/service-competition/internal/rpc/user_activity_full"
	"github.com/varkis-ms/service-competition/internal/rpc/user_activity_total"

	competitiongrpc "github.com/varkis-ms/service-competition/internal/grpc/competition"
	"github.com/varkis-ms/service-competition/internal/rpc/create_competition"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int64
}

// New creates new gRPC server app.
func New(
	log *slog.Logger,
	port int64,
	competitionCreateHandler *create_competition.Handler,
	competitionEditHandler *edit_competition.Handler,
	competitionListHandler *competition_list.Handler,
	getCompetitionInfoHandler *get_competition_info.Handler,
	getLeaderboardHandler *get_leaderboard.Handler,
	userActivityFullHandler *user_activity_full.Handler,
	userActivityTotalHandler *user_activity_total.Handler,
	saveSolutionResultHandler *save_solution_result.Handler,
	getNextSolution *get_next_solution.Handler,
	saveSolution *save_solution.Handler,
) *App {
	loggingOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(
			grpclog.StartCall, grpclog.FinishCall,
		),
		// Add any other option (check functions starting with logging.With).
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpclog.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
			recovery.UnaryServerInterceptor(recoveryOpts...),
		),
	)

	competitiongrpc.Register(
		gRPCServer,
		competitionCreateHandler,
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

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (a *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Shutdown stops gRPC server.
func (a *App) Shutdown() {
	const op = "grpcapp.Shutdown"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int64("port", a.port))

	a.gRPCServer.GracefulStop()
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
