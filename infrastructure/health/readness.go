package health

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/apex/log"
	"github.com/artem-malko/auth-and-go/infrastructure/server"
	"github.com/pkg/errors"
)

type eventBusSubscriber interface {
	Subscribe(topic string, fn interface{}) error
}

// ReadinessController компонент добавляющий к серверу фукнционал урла отображающего его готовность принимать запросы извне
type ReadinessController struct {
	readinessStatusValue int
	logger               log.Interface
	mutex                sync.RWMutex
	subscriber           eventBusSubscriber
}

// Используется позднее связывание компонентов через систему сообщений
func newReadinessController(eventBus eventBusSubscriber, logger log.Interface) (*ReadinessController, error) {
	rm := &ReadinessController{subscriber: eventBus, logger: logger}

	err := rm.subscriber.Subscribe(server.EventBeforeStart, func() {
		rm.setReadinessStatus(http.StatusOK)
	})
	if err != nil {
		return nil, errors.Wrap(err, "Cant't attach event BeforeStart web server")
	}

	err = rm.subscriber.Subscribe(server.EventBeforeStop, func() {
		rm.setReadinessStatus(http.StatusServiceUnavailable)
	})
	if err != nil {
		return nil, errors.Wrap(err, "Cant't attach event BeforeStop web server")
	}

	return rm, nil
}

func (r *ReadinessController) readinessStatus() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.readinessStatusValue
}

func (r *ReadinessController) setReadinessStatus(status int) {
	r.mutex.Lock()
	r.readinessStatusValue = status
	r.mutex.Unlock()
}

type healthcheckResponse struct {
	AppStatus string `json:"app_status"`
}

func (r *ReadinessController) readinessHandler(w http.ResponseWriter, request *http.Request) {

	status := r.readinessStatus()
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(&healthcheckResponse{
		AppStatus: "OK",
	})

	if err != nil {
		r.logger.WithError(err).Error("")
	}
}
