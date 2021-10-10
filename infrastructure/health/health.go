package health

import (
	"github.com/apex/log"
	"github.com/dimiro1/health"
	"github.com/go-chi/chi/v5"
)

// HealthcheckController компонент добавляющий к серверу фукнционал урла отображающего его готовность принимать запросы извне
type HealthcheckController struct {
	readinessController *ReadinessController
	externalController  *ExternalController
}

// New конструктор менеджера healthcheck
func New(eventBus eventBusSubscriber, logger log.Interface) (*HealthcheckController, error) {
	rc, err := newReadinessController(eventBus, logger)
	if err != nil {
		return nil, err
	}

	ec := newExternalController()

	hc := &HealthcheckController{
		readinessController: rc,
		externalController:  ec,
	}

	return hc, nil
}

// Routes список роутеров, которые добавляются к черверу для интеграции функционала readiness проверок
func (h *HealthcheckController) Routes() chi.Router {
	router := chi.NewRouter()
	router.Get("/readiness", h.readinessController.readinessHandler)
	router.Get("/external", h.externalController.externalHandler)

	return router
}

// AddExternalChecker добавиляет проверку в external healthcheck
func (h *HealthcheckController) AddExternalChecker(name string, checker health.Checker) {
	h.externalController.addChecker(name, checker)
}
