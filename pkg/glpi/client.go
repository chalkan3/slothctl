package glpi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GLPIClient holds the necessary configuration and tokens for API interaction.
type GLPIClient struct {
	BaseURL      string
	AppToken     string
	SessionToken string
	HTTPClient   *http.Client
}

// InitSessionResponse represents the structure of the /initSession API response.
type InitSessionResponse struct {
	SessionToken string `json:"session_token"`
}

// GLPIErrorResponse represents a generic error response from the GLPI API.
type GLPIErrorResponse struct {
	Message string `json:"message"`
}

// NewGLPIClient creates a new GLPI API client.
func NewGLPIClient(baseURL, appToken string) *GLPIClient {
	return &GLPIClient{
		BaseURL:    baseURL,
		AppToken:   appToken,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Authenticate obtains a session token from the GLPI API.
func (c *GLPIClient) Authenticate(username, password string) error {
	authURL := fmt.Sprintf("%s/apirest.php/initSession", c.BaseURL)

	payload := map[string]string{
		"login":    username,
		"password": password,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal authentication payload: %w", err)
	}

	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create authentication request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("App-Token", c.AppToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send authentication request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read authentication response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp GLPIErrorResponse
		json.Unmarshal(body, &errResp) // Attempt to unmarshal error message
		return fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, errResp.Message)
	}

	var authResponse InitSessionResponse
	if err := json.Unmarshal(body, &authResponse); err != nil {
		return fmt.Errorf("failed to unmarshal authentication response: %w", err)
	}

	c.SessionToken = authResponse.SessionToken
	return nil
}

// makeRequest performs an HTTP request to the GLPI API with proper headers.
func (c *GLPIClient) makeRequest(method, path string, body interface{}) ([]byte, error) {
	if c.SessionToken == "" {
		return nil, fmt.Errorf("session token is not set; please authenticate first")
	}

	url := fmt.Sprintf("%s/apirest.php/%s", c.BaseURL, path)

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("App-Token", c.AppToken)
	req.Header.Set("Session-Token", c.SessionToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp GLPIErrorResponse
		json.Unmarshal(respBody, &errResp)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, errResp.Message)
	}

	return respBody, nil
}

// Get performs a GET request.
func (c *GLPIClient) Get(path string) ([]byte, error) {
	return c.makeRequest("GET", path, nil)
}

// Post performs a POST request.
func (c *GLPIClient) Post(path string, body interface{}) ([]byte, error) {
	return c.makeRequest("POST", path, body)
}

// Put performs a PUT request.
func (c *GLPIClient) Put(path string, body interface{}) ([]byte, error) {
	return c.makeRequest("PUT", path, body)
}

// Delete performs a DELETE request.
func (c *GLPIClient) Delete(path string) ([]byte, error) {
	return c.makeRequest("DELETE", path, nil)
}

// ListTickets fetches a list of tickets with optional status filtering.
func (c *GLPIClient) ListTickets(statusIDs []int) ([]Ticket, error) {
	path := "Ticket"
	if len(statusIDs) > 0 {
		// GLPI API uses 'is' operator for multiple values, e.g., search[status][is]=1,2,3
		statusQuery := ""
		for i, id := range statusIDs {
			if i > 0 {
				statusQuery += ","
			}
			statusQuery += fmt.Sprintf("%d", id)
		}
		path = fmt.Sprintf("Ticket?search[status][is]=%s", statusQuery)
	}

	respBody, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets: %w", err)
	}

	var tickets []Ticket
	if err := json.Unmarshal(respBody, &tickets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tickets: %w", err)
	}
	return tickets, nil
}

// GetUsers fetches a list of GLPI users.
func (c *GLPIClient) GetUsers() ([]User, error) {
	respBody, err := c.Get("User")
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	var users []User
	if err := json.Unmarshal(respBody, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}
	return users, nil
}

// GetStatusID maps a human-readable status string to its GLPI integer ID.
func GetStatusID(status string) (int, error) {
	switch status {
	case "new":
		return 1, nil
	case "assigned":
		return 2, nil
	case "planned":
		return 3, nil
	case "pending":
		return 4, nil
	case "solved":
		return 5, nil
	case "closed":
		return 6, nil
	default:
		return 0, fmt.Errorf("unknown status: %s", status)
	}
}

// GetStatusName maps a GLPI integer ID to its human-readable status string.
func GetStatusName(statusID int) string {
	switch statusID {
	case 1:
		return "new"
	case 2:
		return "assigned"
	case 3:
		return "planned"
	case 4:
		return "pending"
	case 5:
		return "solved"
	case 6:
		return "closed"
	default:
		return "unknown"
	}
}

// CreateTicket creates a new GLPI ticket.
func (c *GLPIClient) CreateTicket(ticketInput TicketInput) (*Ticket, error) {
	input := TicketUpdateInput{Input: ticketInput}
	respBody, err := c.Post("Ticket", input)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	var createdTicket Ticket
	// GLPI API returns an array of objects with the created ID, e.g., [{id: 123}]
	var result []map[string]int
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create ticket response: %w", err)
	}

	if len(result) == 0 || result[0]["id"] == 0 {
		return nil, fmt.Errorf("invalid response when creating ticket: %s", string(respBody))
	}

	createdTicket.ID = result[0]["id"]
	// Fetch the full ticket details as the create endpoint only returns the ID
	fullTicket, err := c.GetTicket(createdTicket.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch full ticket details after creation: %w", err)
	}
	return fullTicket, nil
}

// GetTicket fetches a specific ticket by ID.
func (c *GLPIClient) GetTicket(ticketID int) (*Ticket, error) {
	respBody, err := c.Get(fmt.Sprintf("Ticket/%d", ticketID))
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket %d: %w", ticketID, err)
	}

	var ticket Ticket
	if err := json.Unmarshal(respBody, &ticket); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ticket %d: %w", ticketID, err)
	}
	return &ticket, nil
}

// AddTicketFollowup adds a follow-up to a GLPI ticket.
func (c *GLPIClient) AddTicketFollowup(ticketID int, comment string) error {
	input := AddFollowupInput{
		Input: struct {
			TicketsID int    `json:"tickets_id"`
			Content   string `json:"content"`
		}{
			TicketsID: ticketID,
			Content:   comment,
		},
	}
	_, err := c.Post("ITILFollowup", input)
	if err != nil {
		return fmt.Errorf("failed to add followup to ticket %d: %w", ticketID, err)
	}
	return nil
}

// UpdateTicket updates a GLPI ticket.
func (c *GLPIClient) UpdateTicket(ticketID int, ticketInput TicketInput) error {
	input := TicketUpdateInput{Input: ticketInput}
	_, err := c.Put(fmt.Sprintf("Ticket/%d", ticketID), input)
	if err != nil {
		return fmt.Errorf("failed to update ticket %d: %w", ticketID, err)
	}
	return nil
}
