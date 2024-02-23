package endpoint

type Handler interface {
	State() State
	Handle() error
}

type State int

const (
	Unknown State = iota
	ServerStart
	ServerStop
	ClientStart
	ClientStop
	ExtractorChanges
	InjectorChanges
)

func (e *Endpoint) serverStarted() {
	if h, ok := e.handlers[ExtractorChanges]; ok {
		e.server.AddExtHandler(h)
	}
	if h, ok := e.handlers[InjectorChanges]; ok {
		e.server.AddInjHandler(h)
	}

	e.logger.Info("syncer endpoint server started")
	if h, ok := e.handlers[ServerStart]; ok {
		err := h()
		if err != nil {
			e.logger.Error(err.Error())
		}
	}
}
func (e *Endpoint) serverStopped() {
	e.logger.Info("syncer endpoint server stopped")
	if h, ok := e.handlers[ServerStop]; ok {
		err := h()
		if err != nil {
			e.logger.Error(err.Error())
		}
	}
}

func (e *Endpoint) clientStarted() {
	if h, ok := e.handlers[ExtractorChanges]; ok {
		e.client.AddExtHandler(h)
	}
	if h, ok := e.handlers[InjectorChanges]; ok {
		e.client.AddInjHandler(h)
	}
	e.logger.Info("syncer endpoint client started")
	if h, ok := e.handlers[ClientStart]; ok {
		err := h()
		if err != nil {
			e.logger.Error(err.Error())
		}
	}
}

func (e *Endpoint) clientStopped() {
	e.logger.Info("syncer endpoint client stopped")
	if h, ok := e.handlers[ClientStop]; ok {
		err := h()
		if err != nil {
			e.logger.Error(err.Error())
		}
	}
}
