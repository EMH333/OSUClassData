package database

// Class ClassIdentifier and TermID as primary key
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

// ClassInfo info about a class not specific to a term
type ClassInfo struct {
	ClassIdentifier     string
	Credits             int
	ClassName           string
	RetrievedClassName  bool //has the name been retrieved from the OSU API?
	NormalizedClassName bool //has the name been normalized?
	ClassDescription    string
}

type Professor struct {
	ID         int
	Name       string
	Department string
}

// ProfessorTerm struct to tie a professor with a certain term of a class
//All 3 values should be used as primary key
type ProfessorTerm struct {
	ProfessorID     int
	TermID          string
	ClassIdentifier string
}
