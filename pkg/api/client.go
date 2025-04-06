package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Reisender/canvas-cli-v2/pkg/config"
)

// Client represents a Canvas API client
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new Canvas API client
func NewClient() *Client {
	cfg := config.GetConfig()

	return &Client{
		BaseURL:    cfg.BaseURL,
		APIKey:     cfg.APIKey,
		HTTPClient: &http.Client{},
	}
}

// Request makes an API request to Canvas
func (c *Client) Request(method, path string, query url.Values) ([]byte, error) {
	// Build the URL
	endpoint, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	endpoint.Path += path

	if query != nil {
		endpoint.RawQuery = query.Encode()
	}

	// Create the request
	req, err := http.NewRequest(method, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add auth header
	req.Header.Add("Authorization", "Bearer "+c.APIKey)

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return body, nil
}

// RequestWithBody makes an API request with a JSON body
func (c *Client) RequestWithBody(method, path string, query url.Values, body interface{}) ([]byte, error) {
	// Build the URL
	endpoint, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	endpoint.Path += path

	if query != nil {
		endpoint.RawQuery = query.Encode()
	}

	// Marshal the body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	// Create the request
	req, err := http.NewRequest(method, endpoint.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add headers
	req.Header.Add("Authorization", "Bearer "+c.APIKey)
	req.Header.Add("Content-Type", "application/json")

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		responseBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(responseBody))
	}

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return responseBody, nil
}

// GetCourses retrieves courses from Canvas
func (c *Client) GetCourses() ([]Course, error) {
	data, err := c.Request("GET", "/courses", nil)
	if err != nil {
		return nil, err
	}

	var courses []Course
	if err := json.Unmarshal(data, &courses); err != nil {
		return nil, fmt.Errorf("error parsing courses: %w", err)
	}

	return courses, nil
}

// GetAssignments retrieves assignments for a course
func (c *Client) GetAssignments(courseID string) ([]Assignment, error) {
	path := fmt.Sprintf("/courses/%s/assignments", courseID)
	data, err := c.Request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var assignments []Assignment
	if err := json.Unmarshal(data, &assignments); err != nil {
		return nil, fmt.Errorf("error parsing assignments: %w", err)
	}

	return assignments, nil
}

// GetUsers retrieves users for a course with pagination support
func (c *Client) GetUsers(courseID string, page int, perPage int) ([]User, error) {
	path := fmt.Sprintf("/courses/%s/users", courseID)
	query := url.Values{}
	query.Add("include[]", "email") // Include email addresses

	// Add pagination parameters
	if page > 0 {
		query.Add("page", strconv.Itoa(page))
	}
	if perPage > 0 {
		query.Add("per_page", strconv.Itoa(perPage))
	} else {
		// Default to 50 per page if not specified
		query.Add("per_page", "50")
	}

	data, err := c.Request("GET", path, query)
	if err != nil {
		return nil, err
	}

	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("error parsing users: %w", err)
	}

	return users, nil
}

// GetUserDetails retrieves detailed information about a user
func (c *Client) GetUserDetails(userID string) (*User, error) {
	path := fmt.Sprintf("/users/%s", userID)
	query := url.Values{}
	query.Add("include[]", "email")

	data, err := c.Request("GET", path, query)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("error parsing user details: %w", err)
	}

	return &user, nil
}

// EnrollmentRequest represents the request body for enrolling a user
type EnrollmentRequest struct {
	UserID          string `json:"user_id"`
	Type            string `json:"type"`
	EnrollmentState string `json:"enrollment_state,omitempty"`
	CourseSection   string `json:"course_section_id,omitempty"`
	LimitPrivileges bool   `json:"limit_privileges_to_course_section,omitempty"`
	Notify          bool   `json:"notify,omitempty"`
}

// AddUserToCourse enrolls a user in a course
func (c *Client) AddUserToCourse(courseID, userID, enrollmentType string, notify bool) (*Enrollment, error) {
	path := fmt.Sprintf("/courses/%s/enrollments", courseID)

	// Create the enrollment request
	enrollReq := EnrollmentRequest{
		UserID: userID,
		Type:   enrollmentType, // e.g., "StudentEnrollment", "TeacherEnrollment", etc.
		Notify: notify,
	}

	// Wrap in the enrollment object expected by the API
	reqBody := map[string]EnrollmentRequest{
		"enrollment": enrollReq,
	}

	data, err := c.RequestWithBody("POST", path, nil, reqBody)
	if err != nil {
		return nil, err
	}

	var enrollment Enrollment
	if err := json.Unmarshal(data, &enrollment); err != nil {
		return nil, fmt.Errorf("error parsing enrollment response: %w", err)
	}

	return &enrollment, nil
}

// GetEnrollments retrieves enrollments for a course
func (c *Client) GetEnrollments(courseID string) ([]Enrollment, error) {
	path := fmt.Sprintf("/courses/%s/enrollments", courseID)

	data, err := c.Request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var enrollments []Enrollment
	if err := json.Unmarshal(data, &enrollments); err != nil {
		return nil, fmt.Errorf("error parsing enrollments: %w", err)
	}

	return enrollments, nil
}

// RemoveUserFromCourse deletes a user's enrollment in a course
func (c *Client) RemoveUserFromCourse(courseID, enrollmentID string) error {
	path := fmt.Sprintf("/courses/%s/enrollments/%s", courseID, enrollmentID)
	query := url.Values{}
	query.Add("task", "delete")

	_, err := c.Request("DELETE", path, query)
	return err
}

// RemoveUserByID removes a user from a course by user ID
func (c *Client) RemoveUserByID(courseID, userID string) error {
	// First, get all enrollments for the course
	enrollments, err := c.GetEnrollments(courseID)
	if err != nil {
		return fmt.Errorf("error fetching enrollments: %w", err)
	}

	// Convert userID to int for comparison
	uid, err := strconv.Atoi(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Find the enrollment for this user
	var found bool
	for _, enrollment := range enrollments {
		if enrollment.UserID == uid {
			// Found the enrollment, now remove it
			err := c.RemoveUserFromCourse(courseID, strconv.Itoa(enrollment.ID))
			if err != nil {
				return fmt.Errorf("error removing enrollment: %w", err)
			}
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("no enrollment found for user %s in course %s", userID, courseID)
	}

	return nil
}

// CreateAssignment creates a new assignment in a course
func (c *Client) CreateAssignment(courseID string, assignment *Assignment) (*Assignment, error) {
	path := fmt.Sprintf("/courses/%s/assignments", courseID)

	// Create the request body
	requestBody := map[string]interface{}{
		"assignment": map[string]interface{}{
			"name":             assignment.Name,
			"description":      assignment.Description,
			"points_possible":  assignment.PointsPossible,
			"due_at":           assignment.DueAt.Format(time.RFC3339),
			"published":        assignment.Published,
			"grading_type":     assignment.GradingType,
			"submission_types": assignment.SubmissionTypes,
		},
	}

	// For optional time fields, only include them if they are set
	if !assignment.UnlockAt.IsZero() {
		requestBody["assignment"].(map[string]interface{})["unlock_at"] = assignment.UnlockAt.Format(time.RFC3339)
	}
	if !assignment.LockAt.IsZero() {
		requestBody["assignment"].(map[string]interface{})["lock_at"] = assignment.LockAt.Format(time.RFC3339)
	}

	// Make the API request
	data, err := c.RequestWithBody("POST", path, nil, requestBody)
	if err != nil {
		return nil, fmt.Errorf("error creating assignment: %w", err)
	}

	// Parse the response
	var newAssignment Assignment
	if err := json.Unmarshal(data, &newAssignment); err != nil {
		return nil, fmt.Errorf("error parsing assignment response: %w", err)
	}

	return &newAssignment, nil
}
