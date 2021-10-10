package health

import (
	"net/http"

	"github.com/dimiro1/health"
)

// ExternalController компонент добавляющий к серверу фукнционал урла отображающего статусы его внешних зависимостей
type ExternalController struct {
	external health.Handler
}

func newExternalController() *ExternalController {
	ec := ExternalController{
		external: health.NewHandler(),
	}
	return &ec
}

func (e *ExternalController) externalHandler(response http.ResponseWriter, request *http.Request) {
	e.external.ServeHTTP(response, request)
}

func (e *ExternalController) addChecker(name string, checker health.Checker) {
	e.external.AddChecker(name, checker)
}
