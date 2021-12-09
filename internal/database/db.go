package database

type Class struct {
	ID int

	ClassIdentifier string
	TermID          string
	Students        int
	Credits         int
	ClassGPA        float64

	A      int
	AMinus int
	B      int
	BMinus int
	BPlus  int
	C      int
	CMinus int
	CPlus  int
	D      int
	DMinus int
	DPlus  int
	F      int
	S      int
	U      int
	W      int
	I      int // All incompletes regardless of grade (for now)
}

type Professor struct {
	ID         int
	Name       string
	Department string
}

//struct to tie a professor with a certian term of a classs
type ProfessorTerm struct {
	ID              int
	ProfessorID     int
	TermID          string
	ClassIdentifier string
}
