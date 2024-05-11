package model

import "time"

type CompetitionList []Competition

type Competition struct {
	UserID             int64
	CompetitionID      int64
	Title              string
	Description        string
	DatasetTitle       string
	DatasetDescription string
}

type CompetitionEdit struct {
	UserID             int64
	CompetitionID      int64
	Title              *string
	Description        *string
	DatasetTitle       *string
	DatasetDescription *string
}

type LeaderBoardList []LeaderBoard

type LeaderBoard struct {
	UserID        int64
	CompetitionID int64
	Score         float32
	AddedAt       time.Time
}

type UserActivityTotal struct {
	TotalTime              string
	TotalAttempts          int64
	TotalCompetitions      int64
	TotalOwnerCompetitions int64
}

type UserActivityFull struct {
	Owner  []CompetitionInfoFullOwner
	Member []CompetitionInfoFull
}

type CompetitionInfoFull struct {
	CompetitionID int64
	Title         string
	DatasetTitle  string
	Score         float32
	AddedAt       time.Time
}

type CompetitionInfoFullOwner struct {
	CompetitionID int64
	Title         string
	DatasetTitle  string
	AmountUsers   int64
	AddedAt       time.Time
}
