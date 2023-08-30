package decorator

type Scene[D any] interface {
	sealed()

	D() D
	GetCall(string) (any, bool)
}

func NewScene[D any](d D, opt Option[D], opts ...Option[D]) (Scene[D], error) {
	registry := make(map[string]any)
	if err := opt(registry); err != nil {
		return nil, err
	}

	for _, o := range opts {
		if err := o(registry); err != nil {
			return nil, err
		}
	}

	return &scene[D]{d: d, registry: registry}, nil
}

type scene[D any] struct {
	d        D
	registry map[string]any
}

func (d *scene[D]) sealed() {}

func (d *scene[D]) D() D {
	return d.d
}

func (d *scene[D]) GetCall(name string) (any, bool) {
	m, ok := d.registry[name]
	return m, ok
}
