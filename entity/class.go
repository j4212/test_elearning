package entity

// This folder is handling related into data layer.
// For example: In GORM if you want to specify column to select, you have to create your own struct.

// Ex: There's table Courses that concist of column 'id', 'title', 'description'.
// And then you want to select only 'title' & 'description'. Then create struct like below:

/*
type CourseEntity struct {
	Title string
	Description string
}
*/

// Next, in your repository logic

/*
(Wrong)
func (repos *reposImpl) FindCourse(slug string) {
	var course model.Course
	tx := repos.DB.First(&course, "slug = ?", slug)
}
*/

/*
(Correct)
func (repos *reposImpl) FindCourse(slug string) {
	var course entity.CourseEntity

	tx := repos.DB.First(&course, "slug = ?", slug)
}
*/

type ClassEntity struct {
	Class string
}
