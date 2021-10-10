package caller

import (
	"runtime"
	"strings"
)

func GetCaller(depth int) (string, bool) {
	// Remove Get Caller from stack
	pc, _, _, ok := runtime.Caller(depth + 1)
	details := runtime.FuncForPC(pc)

	if ok && details != nil {
		removedPathPrefix := strings.
			Replace(details.Name(), "github.com/artem-malko/auth-and-go", "", 1)

		removedChallengesManagerPrefix := strings.
			Replace(removedPathPrefix, "/managers/challenges/manager.(*challengesManager).", "", 1)

		removedUsersManagerPrefix := strings.
			Replace(removedChallengesManagerPrefix, "/managers/user/manager.(*userManager).", "", 1)

		removeChallengesRepositoryPrefix := strings.
			Replace(removedUsersManagerPrefix, "/repositories/challenges/repository.(*repository).", "", 1)

		removeChallengesRouterPrefix := strings.
			Replace(removeChallengesRepositoryPrefix, "/api/routers/challenges.(*handlers).", "", 1)

		return removeChallengesRouterPrefix, true
	}

	return "unknown", false
}
