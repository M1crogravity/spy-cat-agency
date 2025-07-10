package main

// ErrorResponse represents a standard error response
// @Description Standard error response format
// @Example {"error": "error message"}
//
// swagger:model ErrorResponse
type ErrorResponseDoc struct {
	// Error message
	// Example: error message
	Error string `json:"error"`
}

// ValidationErrorResponse represents a validation error response
// @Description Validation error response format with field-specific errors
// @Example {"error": {"field1": "error message", "field2": "error message"}}
//
// swagger:model ValidationErrorResponse
type ValidationErrorResponseDoc struct {
	// Map of field names to error messages
	// Example: {"name": "must be provided", "breed": "invalid breed"}
	Error map[string]string `json:"error"`
}

// MessageResponse represents a simple message response
// @Description Simple message response format
// @Example {"message": "operation completed successfully"}
//
// swagger:model MessageResponse
type MessageResponseDoc struct {
	// Success message
	// Example: operation completed successfully
	Message string `json:"message"`
}

// SpyCatResponse represents a spy cat response
// @Description Response containing a single spy cat
// @Example {"spy-cat": {"id": 1, "name": "Agent Whiskers", "years_of_experience": 5, "breed": "Siamese", "salary": 50000.00}}
//
// swagger:model SpyCatResponse
type SpyCatResponseDoc struct {
	// Spy cat data
	SpyCat SpyCatDoc `json:"spy-cat"`
}

// SpyCatsResponse represents a list of spy cats response
// @Description Response containing a list of spy cats
// @Example {"spy-cats": [{"id": 1, "name": "Agent Whiskers", "years_of_experience": 5, "breed": "Siamese", "salary": 50000.00}]}
//
// swagger:model SpyCatsResponse
type SpyCatsResponseDoc struct {
	// List of spy cats
	SpyCats []SpyCatDoc `json:"spy-cats"`
}

// SpyCat represents a spy cat
// @Description Spy cat entity
// @Example {"id": 1, "name": "Agent Whiskers", "years_of_experience": 5, "breed": "Siamese", "salary": 50000.00}
//
// swagger:model SpyCat
type SpyCatDoc struct {
	// Unique identifier
	// Example: 1
	ID int64 `json:"id"`
	// Spy cat name
	// Example: Agent Whiskers
	Name string `json:"name"`
	// Years of experience
	// Example: 5
	YearsOfExperience int `json:"years_of_experience"`
	// Cat breed
	// Example: Siamese
	Breed string `json:"breed"`
	// Annual salary
	// Example: 50000.00
	Salary float64 `json:"salary"`
}

// CreateSpyCatRequest represents the request body for creating a spy cat
// @Description Request body for creating a new spy cat
// @Example {"name": "Agent Whiskers", "years_of_experience": 5, "breed": "Siamese", "salary": 50000.00, "password": "secretpassword123"}
//
// swagger:model CreateSpyCatRequest
type CreateSpyCatRequestDoc struct {
	// Spy cat name
	// Example: Agent Whiskers
	Name string `json:"name"`
	// Years of experience
	// Example: 5
	YearsOfExperience int `json:"years_of_experience"`
	// Cat breed
	// Example: Siamese
	Breed string `json:"breed"`
	// Annual salary
	// Example: 50000.00
	Salary float64 `json:"salary"`
	// Password for authentication
	// Example: secretpassword123
	Password string `json:"password"`
}

// UpdateSpyCatSalaryRequest represents the request body for updating spy cat salary
// @Description Request body for updating a spy cat's salary
// @Example {"salary": 55000.00}
//
// swagger:model UpdateSpyCatSalaryRequest
type UpdateSpyCatSalaryRequestDoc struct {
	// New salary amount
	// Example: 55000.00
	Salary float64 `json:"salary"`
}

// MissionResponse represents a mission response
// @Description Response containing a single mission
// @Example {"mission": {"id": 1, "state": "created", "assigned_cat_id": 1, "targets": []}}
//
// swagger:model MissionResponse
type MissionResponseDoc struct {
	// Mission data
	Mission MissionDoc `json:"mission"`
}

// MissionsResponse represents a list of missions response
// @Description Response containing a list of missions
// @Example {"missions": [{"id": 1, "state": "created", "assigned_cat_id": 1, "targets": []}]}
//
// swagger:model MissionsResponse
type MissionsResponseDoc struct {
	// List of missions
	Missions []MissionDoc `json:"missions"`
}

// Mission represents a mission
// @Description Mission entity
// @Example {"id": 1, "state": "created", "assigned_cat_id": 1, "targets": [{"id": 1, "mission_id": 1, "name": "Dr. Evil", "country": "Switzerland", "notes": "", "state": "created"}]}
//
// swagger:model Mission
type MissionDoc struct {
	// Unique identifier
	// Example: 1
	ID int64 `json:"id"`
	// Mission state (created, in_progress, completed)
	// Example: created
	State string `json:"state"`
	// ID of assigned spy cat
	// Example: 1
	AssignedCatID int64 `json:"assigned_cat_id"`
	// List of mission targets
	Targets []TargetDoc `json:"targets"`
}

// CreateMissionRequest represents the request body for creating a mission
// @Description Request body for creating a new mission
// @Example {"targets": [{"name": "Dr. Evil", "country": "Switzerland"}]}
//
// swagger:model CreateMissionRequest
type CreateMissionRequestDoc struct {
	// List of targets for the mission
	Targets []CreateTargetRequestDoc `json:"targets"`
}

// Target represents a mission target
// @Description Mission target entity
// @Example {"id": 1, "mission_id": 1, "name": "Dr. Evil", "country": "Switzerland", "notes": "Target spotted at secret lair", "state": "created"}
//
// swagger:model Target
type TargetDoc struct {
	// Unique identifier
	// Example: 1
	ID int64 `json:"id"`
	// Mission ID this target belongs to
	// Example: 1
	MissionID int64 `json:"mission_id"`
	// Target name
	// Example: Dr. Evil
	Name string `json:"name"`
	// Target country
	// Example: Switzerland
	Country string `json:"country"`
	// Notes about the target
	// Example: Target spotted at secret lair
	Notes string `json:"notes"`
	// Target state (created, in_progress, completed)
	// Example: created
	State string `json:"state"`
}

// TargetResponse represents a target response
// @Description Response containing a single target
// @Example {"target": {"id": 1, "mission_id": 1, "name": "Dr. Evil", "country": "Switzerland", "notes": "", "state": "created"}}
//
// swagger:model TargetResponse
type TargetResponseDoc struct {
	// Target data
	Target TargetDoc `json:"target"`
}

// CreateTargetRequest represents the request body for creating a target
// @Description Request body for creating a new target
// @Example {"name": "Dr. Evil", "country": "Switzerland"}
//
// swagger:model CreateTargetRequest
type CreateTargetRequestDoc struct {
	// Target name
	// Example: Dr. Evil
	Name string `json:"name"`
	// Target country
	// Example: Switzerland
	Country string `json:"country"`
}

// UpdateTargetNotesRequest represents the request body for updating target notes
// @Description Request body for updating target notes
// @Example {"notes": "Target spotted at secret lair"}
//
// swagger:model UpdateTargetNotesRequest
type UpdateTargetNotesRequestDoc struct {
	// Notes about the target
	// Example: Target spotted at secret lair
	Notes string `json:"notes"`
}

// AgentResponse represents an agent response
// @Description Response containing a single agent
// @Example {"agent": {"id": 1, "name": "Agent Smith"}}
//
// swagger:model AgentResponse
type AgentResponseDoc struct {
	// Agent data
	Agent AgentDoc `json:"agent"`
}

// Agent represents an agent
// @Description Agent entity
// @Example {"id": 1, "name": "Agent Smith"}
//
// swagger:model Agent
type AgentDoc struct {
	// Unique identifier
	// Example: 1
	ID int64 `json:"id"`
	// Agent name
	// Example: Agent Smith
	Name string `json:"name"`
}

// CreateAgentRequest represents the request body for creating an agent
// @Description Request body for creating a new agent
// @Example {"name": "Agent Smith", "password": "agentpassword123"}
//
// swagger:model CreateAgentRequest
type CreateAgentRequestDoc struct {
	// Agent name
	// Example: Agent Smith
	Name string `json:"name"`
	// Password for authentication
	// Example: agentpassword123
	Password string `json:"password"`
}

// TokenResponse represents an authentication token response
// @Description Response containing an authentication token
// @Example {"authentication_token": {"plaintext": "ABCDEF123456", "user_id": 1, "expiry": "2024-01-01T00:00:00Z", "scope": "authentication"}}
//
// swagger:model TokenResponse
type TokenResponseDoc struct {
	// Authentication token data
	AuthenticationToken TokenDoc `json:"authentication_token"`
}

// Token represents an authentication token
// @Description Authentication token entity
// @Example {"plaintext": "ABCDEF123456", "user_id": 1, "expiry": "2024-01-01T00:00:00Z", "scope": "authentication"}
//
// swagger:model Token
type TokenDoc struct {
	// Token plaintext value
	// Example: ABCDEF123456
	Plaintext string `json:"plaintext"`
	// User ID associated with the token
	// Example: 1
	UserID int64 `json:"user_id"`
	// Token expiry date
	// Example: 2024-01-01T00:00:00Z
	Expiry string `json:"expiry"`
	// Token scope
	// Example: authentication
	Scope string `json:"scope"`
}

// AuthenticationRequest represents the request body for authentication
// @Description Request body for user authentication
// @Example {"name": "Agent Smith", "password": "password123"}
//
// swagger:model AuthenticationRequest
type AuthenticationRequestDoc struct {
	// User name
	// Example: Agent Smith
	Name string `json:"name"`
	// User password
	// Example: password123
	Password string `json:"password"`
}
