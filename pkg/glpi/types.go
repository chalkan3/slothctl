package glpi

// GLPIInstance represents a configured GLPI instance.
type GLPIInstance struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	AppToken string `json:"app_token"`
	User     string `json:"user"`
	Password string `json:"password,omitempty"` // Omit if empty for security
}

// Ticket represents a GLPI ticket.
type Ticket struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Content     string `json:"content"`
	Status      int    `json:"status"` // 1=new, 2=assigned, 3=planned, 4=pending, 5=solved, 6=closed
	Urgency     int    `json:"urgency"`
	Impact      int    `json:"impact"`
	RequesterID int    `json:"users_id_requester"`
	AssigneeID  int    `json:"users_id_assign"`
	// Add more fields as needed based on GLPI API documentation
}

// TicketInput represents the data for creating or updating a ticket.
type TicketInput struct {
	Name        string `json:"name,omitempty"`
	Content     string `json:"content,omitempty"`
	Status      int    `json:"status,omitempty"`
	Urgency     int    `json:"urgency,omitempty"`
	Impact      int    `json:"impact,omitempty"`
	RequesterID int    `json:"users_id_requester,omitempty"`
	AssigneeID  int    `json:"users_id_assign,omitempty"`
}

// TicketUpdateInput represents the data for updating a ticket.
type TicketUpdateInput struct {
	Input TicketInput `json:"input"`
}

// AddFollowupInput represents the data for adding a followup to a ticket.
type AddFollowupInput struct {
	Input struct {
		TicketsID int    `json:"tickets_id"`
		Content   string `json:"content"`
	} `json:"input"`
}

// User represents a GLPI user.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
