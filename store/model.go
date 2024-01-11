package store

type Workout struct {
	Title       string
	Description string
	Color       int
	Routines    []Routine
}

type Routine struct {
	Title       string
	Description string
}
