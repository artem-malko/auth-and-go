package models

import "github.com/google/uuid"

type MonthlyStats struct {
	ID                uuid.UUID
	ChallengerID      uuid.UUID
	StartedChallenges int
	AchivedChallenges int
	CreatedChallenges int
	FinishedTasks     int
	Posts             int
	Feedbacks         int
	Month             int
	Year              int
}

type MonthlyStatsType string

const (
	MonthlyStatsTypeAcceptedChallenges MonthlyStatsType = "acceptedChallenges"
	MonthlyStatsTypeAchivedChallenges  MonthlyStatsType = "achivedChallenges"
	MonthlyStatsTypeCreatedChallenges  MonthlyStatsType = "createdChallenges"
	MonthlyStatsTypeFinishedTasks      MonthlyStatsType = "finishedTasks"
	MonthlyStatsTypePosts              MonthlyStatsType = "posts"
	MonthlyStatsTypeFeedbacks          MonthlyStatsType = "feedbacks"
)
