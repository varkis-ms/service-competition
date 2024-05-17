package competition

import (
	"context"
	"errors"

	"github.com/varkis-ms/service-competition/internal/model"
	"github.com/varkis-ms/service-competition/internal/pkg/pb"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	competition.UnimplementedCompetitionServer
	competitionCreateHandler  *create_competition.Handler
	competitionEditHandler    *edit_competition.Handler
	competitionListHandler    *competition_list.Handler
	getCompetitionInfoHandler *get_competition_info.Handler
	getLeaderboardHandler     *get_leaderboard.Handler
	userActivityFullHandler   *user_activity_full.Handler
	userActivityTotalHandler  *user_activity_total.Handler
	saveSolutionResultHandler *save_solution_result.Handler
	getNextSolution           *get_next_solution.Handler
	saveSolution              *save_solution.Handler
}

func Register(
	gRPCServer *grpc.Server,
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
) {
	competition.RegisterCompetitionServer(gRPCServer, &server{
		competitionCreateHandler:  competitionCreateHandler,
		competitionEditHandler:    competitionEditHandler,
		competitionListHandler:    competitionListHandler,
		getCompetitionInfoHandler: getCompetitionInfoHandler,
		getLeaderboardHandler:     getLeaderboardHandler,
		userActivityFullHandler:   userActivityFullHandler,
		userActivityTotalHandler:  userActivityTotalHandler,
		saveSolutionResultHandler: saveSolutionResultHandler,
		getNextSolution:           getNextSolution,
		saveSolution:              saveSolution,
	})
	reflection.Register(gRPCServer)
}

func (s *server) CompetitionCreate(
	ctx context.Context,
	in *competition.CompetitionCreateRequest,
) (*competition.CompetitionCreateResponse, error) {
	if in.GetUserID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "userID is required")
	}

	if in.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	if in.GetDescription() == "" {
		return nil, status.Error(codes.InvalidArgument, "description is required")
	}

	if in.GetDatasetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "dataset title is required")
	}

	if in.GetDatasetDescription() == "" {
		return nil, status.Error(codes.InvalidArgument, "dataset description is required")
	}

	out := &competition.CompetitionCreateResponse{}
	if err := s.competitionCreateHandler.Handle(ctx, in, out); err != nil {
		if errors.Is(err, model.ErrCompExists) {
			return nil, status.Error(codes.AlreadyExists, model.ErrCompExists.Error())
		}

		return nil, status.Error(codes.Internal, "failed to create competition")
	}

	return out, nil
}

func (s *server) CompetitionEdit(
	ctx context.Context,
	in *competition.CompetitionEditRequest,
) (*emptypb.Empty, error) {
	if in.GetUserID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "userID is required")
	}

	if in.GetCompetitionID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "userID is required")
	}

	// TODO: задать ограничения
	//if in.GetTitle() == "" {
	//	return nil, status.Error(codes.InvalidArgument, "title is required")
	//}
	//
	//if in.GetDescription() == "" {
	//	return nil, status.Error(codes.InvalidArgument, "description is required")
	//}
	//
	//if in.GetDatasetTitle() == "" {
	//	return nil, status.Error(codes.InvalidArgument, "dataset title is required")
	//}
	//
	//if in.GetDatasetDescription() == "" {
	//	return nil, status.Error(codes.InvalidArgument, "dataset description is required")
	//}

	if err := s.competitionEditHandler.Handle(ctx, in); err != nil {
		if errors.Is(err, model.ErrCompExists) {
			return nil, status.Error(codes.AlreadyExists, model.ErrCompExists.Error())
		}
		if errors.Is(err, model.ErrCompNotFound) {
			return nil, status.Error(codes.NotFound, model.ErrCompNotFound.Error())
		}
		if errors.Is(err, model.ErrNoAccessToComp) {
			return nil, status.Error(codes.InvalidArgument, model.ErrNoAccessToComp.Error())
		}

		return nil, status.Error(codes.Internal, "failed to edit competition")
	}

	return &emptypb.Empty{}, nil
}

func (s *server) CompetitionList(
	ctx context.Context,
	in *competition.CompetitionListRequest,
) (*competition.CompetitionListResponse, error) {
	if in.GetUserID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "userID is required")
	}

	out := &competition.CompetitionListResponse{}
	if err := s.competitionListHandler.Handle(ctx, in, out); err != nil {
		return nil, status.Error(codes.Internal, "failed to get competition list")
	}

	return out, nil
}

func (s *server) GetCompetitionInfo(
	ctx context.Context,
	in *competition.CompetitionInfoRequest,
) (*competition.CompetitionInfoResponse, error) {
	if in.GetUserID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "userID is required")
	}

	if in.GetCompetitionID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "competitionID is required")
	}

	out := &competition.CompetitionInfoResponse{}
	if err := s.getCompetitionInfoHandler.Handle(ctx, in, out); err != nil {
		if errors.Is(err, model.ErrCompNotFound) {
			return nil, status.Error(codes.NotFound, model.ErrCompNotFound.Error())
		}

		return nil, status.Error(codes.Internal, "failed to get competition list")
	}

	return out, nil
}

func (s *server) LeaderBoard(
	ctx context.Context,
	in *competition.LeaderBoardRequest,
) (*competition.LeaderBoardResponse, error) {
	if in.GetCompetitionID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "competitionID is required")
	}

	out := &competition.LeaderBoardResponse{}
	if err := s.getLeaderboardHandler.Handle(ctx, in, out); err != nil {
		return nil, status.Error(codes.Internal, "failed to get leaderboard")
	}

	return out, nil
}

func (s *server) UserActivityTotal(
	ctx context.Context,
	in *competition.UserActivityTotalRequest,
) (*competition.UserActivityTotalResponse, error) {
	if in.GetUserID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	out := &competition.UserActivityTotalResponse{}
	if err := s.userActivityTotalHandler.Handle(ctx, in, out); err != nil {
		return nil, status.Error(codes.Internal, "failed to get userActivityTotal")
	}

	return out, nil
}
func (s *server) UserActivityFull(
	ctx context.Context,
	in *competition.UserActivityFullRequest,
) (*competition.UserActivityFullResponse, error) {
	if in.GetUserID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	out := &competition.UserActivityFullResponse{}
	if err := s.userActivityFullHandler.Handle(ctx, in, out); err != nil {
		return nil, status.Error(codes.Internal, "failed to get UserActivityFull")
	}

	return out, nil
}

func (s *server) SaveSolutionResult(
	ctx context.Context,
	in *competition.SaveSolutionResultRequest,
) (*emptypb.Empty, error) {
	if in.GetSolutionID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "solution id is required")
	}

	if in.GetScore() == 0 {
		return nil, status.Error(codes.InvalidArgument, "score is required")
	}

	if in.GetRunTime() == "" {
		return nil, status.Error(codes.InvalidArgument, "run time is required")
	}

	if err := s.saveSolutionResultHandler.Handle(ctx, in); err != nil {
		return nil, status.Error(codes.Internal, "failed to save solution result")
	}

	return &emptypb.Empty{}, nil
}

func (s *server) GetNextSolution(
	ctx context.Context,
	_ *emptypb.Empty,
) (*competition.GetNextSolutionResponse, error) {
	out := &competition.GetNextSolutionResponse{}
	if err := s.getNextSolution.Handle(ctx, out); err != nil {
		return nil, status.Error(codes.Internal, "failed to get next solution")
	}

	return out, nil
}

func (s *server) SaveSolution(
	ctx context.Context,
	in *competition.SaveSolutionRequest,
) (*emptypb.Empty, error) {
	if in.GetUserID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	if in.GetCompetitionID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "competition id is required")
	}

	if err := s.saveSolution.Handle(ctx, in); err != nil {
		return nil, status.Error(codes.Internal, "failed to save solution")
	}

	return &emptypb.Empty{}, nil
}
