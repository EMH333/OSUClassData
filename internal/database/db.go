package database

//ClassIdentifier and TermID as primary key
type Class struct {
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

	Visible bool
}

//info about a class not specific to a term
type ClassInfo struct {
	ClassIdentifier string
	ClassName       string
	Credits         int
}

type Professor struct {
	ID         int
	Name       string
	Department string
}

//struct to tie a professor with a certian term of a classs
//All 3 values should be used as primary key
type ProfessorTerm struct {
	ProfessorID     int
	TermID          string
	ClassIdentifier string
}
