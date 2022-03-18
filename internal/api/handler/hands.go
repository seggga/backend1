package handler

// Hands ...
type Hands struct {
	dir string
}

// New ...
func New(dir string) *Hands {
	return &Hands{
		dir: dir,
	}
}
