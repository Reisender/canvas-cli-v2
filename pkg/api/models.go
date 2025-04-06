package api

import "time"

// Course represents a Canvas course
type Course struct {
	ID                  int       `json:"id"`
	Name                string    `json:"name"`
	CourseCode          string    `json:"course_code"`
	StartAt             time.Time `json:"start_at"`
	EndAt               time.Time `json:"end_at"`
	Workflow            string    `json:"workflow_state"`
	AccountID           int       `json:"account_id"`
	EnrollmentTermID    int       `json:"enrollment_term_id"`
	GradingStandardID   int       `json:"grading_standard_id"`
	CreatedAt           time.Time `json:"created_at"`
	RestrictEnrollments bool      `json:"restrict_enrollments_to_course_dates"`
}

// Assignment represents a Canvas assignment
type Assignment struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	DueAt              time.Time `json:"due_at"`
	LockAt             time.Time `json:"lock_at"`
	UnlockAt           time.Time `json:"unlock_at"`
	CourseID           int       `json:"course_id"`
	PointsPossible     float64   `json:"points_possible"`
	GradingType        string    `json:"grading_type"`
	SubmissionTypes    []string  `json:"submission_types"`
	Published          bool      `json:"published"`
	HTMLURL            string    `json:"html_url"`
	SubmissionsURL     string    `json:"submissions_download_url"`
	GradeGroupStudents bool      `json:"grade_group_students_individually"`
}

// User represents a Canvas user
type User struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	SortableName  string `json:"sortable_name"`
	ShortName     string `json:"short_name"`
	SISUserID     string `json:"sis_user_id"`
	SISImportID   int    `json:"sis_import_id"`
	LoginID       string `json:"login_id"`
	IntegrationID string `json:"integration_id"`
	Email         string `json:"email"`
	Locale        string `json:"locale"`
	Avatar        string `json:"avatar_url"`
}

// Submission represents a Canvas assignment submission
type Submission struct {
	ID              int       `json:"id"`
	AssignmentID    int       `json:"assignment_id"`
	UserID          int       `json:"user_id"`
	SubmittedAt     time.Time `json:"submitted_at"`
	Score           float64   `json:"score"`
	Grade           string    `json:"grade"`
	AttemptNumber   int       `json:"attempt"`
	Body            string    `json:"body"`
	URL             string    `json:"url"`
	GradedAt        time.Time `json:"graded_at"`
	GraderID        int       `json:"grader_id"`
	Late            bool      `json:"late"`
	Missing         bool      `json:"missing"`
	SubmissionType  string    `json:"submission_type"`
	PreviewURL      string    `json:"preview_url"`
	GradeMatchesHub bool      `json:"grade_matches_current_submission"`
}

// Enrollment represents a Canvas enrollment (user enrollment in a course)
type Enrollment struct {
	ID                int       `json:"id"`
	UserID            int       `json:"user_id"`
	CourseID          int       `json:"course_id"`
	Type              string    `json:"type"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	StartAt           time.Time `json:"start_at"`
	EndAt             time.Time `json:"end_at"`
	LastActivityAt    time.Time `json:"last_activity_at"`
	TotalActivityTime int       `json:"total_activity_time"`
	HTMLURL           string    `json:"html_url"`
	Grades            struct {
		HTMLUrl      string  `json:"html_url"`
		CurrentScore float64 `json:"current_score"`
		FinalScore   float64 `json:"final_score"`
		CurrentGrade string  `json:"current_grade"`
		FinalGrade   string  `json:"final_grade"`
	} `json:"grades"`
	User            User   `json:"user"`
	CourseSectionID int    `json:"course_section_id"`
	EnrollmentState string `json:"enrollment_state"`
	LimitPrivileges bool   `json:"limit_privileges_to_course_section"`
	Role            string `json:"role"`
	RoleID          int    `json:"role_id"`
}
