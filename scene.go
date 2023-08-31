package decorator

type Scene interface {
	sealed()

	D() any
	GetCall(string) (any, bool)
}

func NewScene[D any](d D, opt Option[D], opts ...Option[D]) (Scene, error) {
	registry := make(map[string]any)
	if err := opt(registry); err != nil {
		return nil, err
	}

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
