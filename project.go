package togo

type Project struct {
	Name        string
	Description string
	store       ProjectStore
}

func NewProject(name, description string, s ProjectStore) Project {
	return Project{Name: name, Description: description, store: s}
}

func (p Project) TasksByPriority() ([]Task, error) {
	ts, err := p.store.TasksByPriority(p.Name)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (p Project) AddTask(t Task) error {
	return p.store.AddTask(p, t)
}
