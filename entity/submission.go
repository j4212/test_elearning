package entity

type TotalSubmissionStatus struct {
	Pending   int
	Rejected  int
	Approved int
}

type PlaceholderFilterSubmission struct {
	Course TCoursePlaceholder
	ListMaterial []TMaterialPlaceholder
}

type TCoursePlaceholder struct {
	ID string
	Title string 
}

type TMaterialPlaceholder struct {
	ID string
	Title string  
}