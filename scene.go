package decorator

// Represents registry of rules how to decorate functions or method calls using single decorator service.
type Scene interface {
	sealed()

	// Returns decorating service.
	D() any
	//Returns registered function or method call.
	GetCall(string) (any, bool)
}

// Creates new Scene using provided options.
// Needs options to know how to decorate functions or method calls.
func NewScene[D any](d D, opts ...Option[D]) (Scene, error) {
	registry := make(map[string]any)

	for _, o := range opts {
		if err := o(registry); err != nil {
			return nil, err
		}
	}

	return &scene{d: d, registry: registry}, nil
}

type scene struct {
	d        any
	registry map[string]any
}

func (d *scene) sealed() {}

func (d *scene) D() any {
	return d.d
}

func (d *scene) GetCall(name string) (any, bool) {
	m, ok := d.registry[name]
	return m, ok
}
