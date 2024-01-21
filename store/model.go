package store

type Day struct {
	OpenHour  int
	CloseHour int
	Workout   *Workout
}

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

type HourArrivals struct {
	Hour          int
	ArrivingUsers []string
}
