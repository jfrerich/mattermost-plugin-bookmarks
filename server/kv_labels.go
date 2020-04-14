package main

// Labels contains a map of labels with the label name as the key
type Labels struct {
	ByName map[string]*Label
}

// Label defines the parameters of a label
type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (l *Labels) add(label *Label) {
	l.ByName[label.Name] = label
}

func (l *Labels) get(ID string) *Label {
	return l.ByName[ID]
}

func (l *Labels) delete(ID string) {
	delete(l.ByName, ID)
}

func (l *Labels) exists(ID string) (*Label, bool) {
	if label, ok := l.ByName[ID]; ok {
		return label, true
	}
	return nil, false
}
