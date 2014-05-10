package contribot

type Submission struct {
	Name    string
	Address string
	Email   string
	Size    string
}

type Backend func(*Submission)
