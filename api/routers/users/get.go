package users

import (
	"net/http"

	"github.com/artem-malko/auth-and-go/api/middleware"

	"github.com/artem-malko/auth-and-go/api/response"
	userManager "github.com/artem-malko/auth-and-go/managers/user"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (h *handlers) getUserResponse(
	w http.ResponseWriter,
	r *http.Request,
	user *models.User,
	err error,
	method string,
) {
	if err != nil {
		switch errors.Cause(err) {
		case userManager.ErrUserNotFound:
			response.NotFound(w)
			return
		default:
			h.logger(r).
				WithField("method", method).
				WithField("code", http.StatusInternalServerError).
				Error(errors.Wrap(err, "Error during "+method).Error())
			response.InternalServerError(w)
			return
		}
	}

	if user.AccountStatus == models.AccountStatusDeleted ||
		user.AccountStatus == models.AccountStatusBanned {
		response.NotFound(w)
		return
	}

	response.OK(w, user)
}

func (h *handlers) CreateGetUserByID(contextKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		accountID := ctx.Value(contextKey).(uuid.UUID)
		user, err := h.userManager.GetUserByID(accountID)

		h.getUserResponse(w, r, user, err, "GetUserByID")
	}
}

func (h *handlers) CreateGetUserByName(queryParamName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userName := r.URL.Query().Get(queryParamName)

		if userName == "" {
			response.Error(w, http.StatusBadRequest, "username can not be empty")
			return
		}

		user, err := h.userManager.GetUserByName(userName)

		h.getUserResponse(w, r, user, err, "GetUserByName")
	}
}

func (h *handlers) GetFullUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessTokenInfo := middleware.GetAccessTokenInfo(ctx)

	if accessTokenInfo == nil {
		response.Error(w, http.StatusUnauthorized, "You are not authed")
		return
	}

	u, err := h.userManager.GetFullUser(accessTokenInfo.AccountID)

	if err != nil {
		switch errors.Cause(err) {
		case userManager.ErrUserNotFound:
			response.Error(w, http.StatusNotFound, "There is no user with id "+accessTokenInfo.AccountID.String())
		default:
			h.logger(r).
				WithField("method", "GetFullUser").
				WithField("code", http.StatusInternalServerError).
				Error(err.Error())
			response.InternalServerError(w)
		}
		return
	}

	response.OK(w, u)
}

// func (h *handlers) GetAcceptedChallengesForOwnUser(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	accessTokenInfo := middleware.GetAccessTokenInfo(ctx)

// 	if accessTokenInfo == nil {
// 		response.Error(w, http.StatusUnauthorized, "You are not authed")
// 		return
// 	}

// 	formatArg := r.URL.Query().Get("format")

// 	if formatArg == "short" {
// 		h.getShortFormat(w, r)
// 		return
// 	}

// 	response.OKWithoutContent(w)
// }

// type Counter struct {
// 	ChallengeID         uuid.UUID `json:"challenge_id"`
// 	AcceptedChallengeID uuid.UUID `json:"accepted_challenge_id"`
// }

// func (h *handlers) getShortFormat(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	accessTokenInfo := middleware.GetAccessTokenInfo(ctx)

// 	if accessTokenInfo == nil {
// 		response.Error(w, http.StatusUnauthorized, "You are not authed")
// 		return
// 	}

// 	acceptedChallenges, err := h.activityManager.GetAcceptedChallengesByUserID(
// 		accessTokenInfo.AccountID,
// 		models.AcceptedChallengesFilters{
// 			Status: []models.AcceptedChallengeStatus{
// 				models.AcceptedChallengeStatusInProgress,
// 				models.AcceptedChallengeStatusPaused,
// 			},
// 		},
// 	)

// 	if err != nil {
// 		switch errors.Cause(err) {
// 		case challenges.ErrAcceptedChallengesNotFound:
// 			response.NotFound(w)
// 		default:
// 			h.logger(r).
// 				WithField("method", "GetCounters").
// 				WithField("code", http.StatusInternalServerError).
// 				Error(errors.Wrap(err, "Error during GetCounters").Error())
// 			response.InternalServerError(w)
// 		}
// 		return
// 	}

// 	res := []Counter{}

// 	for _, c := range acceptedChallenges {
// 		res = append(res, Counter{
// 			ChallengeID:         c.ChallengeID,
// 			AcceptedChallengeID: c.ID,
// 		})
// 	}

// 	response.OK(w, res)
// }
