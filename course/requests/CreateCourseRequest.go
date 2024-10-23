package requests

type CourseCreateRequest struct {
	Title           string `json:"title" gorm:"size:255;not null"`
	Description     string `json:"description" gorm:"size:1024"`
	Category        string `json:"category" gorm:"size:100;not null"`
	EnrollmentLimit int    `json:"enrollment_limit"`
}
