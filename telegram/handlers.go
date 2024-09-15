package telegram

type Handler struct {
	updates  <-chan Update
	handlers []Handle
}

type Handle struct {
	handler    func(Update)
	predicates []Predicate
}

func NewHandler(updates <-chan Update) *Handler {
	return &Handler{updates: updates}
}

func (h *Handler) Handle(handler func(Update), predicates ...Predicate) {
	h.handlers = append(h.handlers, Handle{handler: handler, predicates: predicates})
}

func (h *Handler) Start() {
	for update := range h.updates {
		for _, handle := range h.handlers {
			if len(handle.predicates) == 0 {
				handle.handler(update)
				break
			}

			needHandle := true

			for _, predicate := range handle.predicates {
				if !predicate(update) {
					needHandle = false
					break
				}
			}

			if needHandle {
				handle.handler(update)
				break
			}
		}
	}
}
