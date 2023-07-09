package responses

// ValidResponse represents a 200 success response
type ValidResponse[T interface{}] struct {
	Success  bool `json:"success"`
	Response T    `json:"response,omitempty"`
}

// InvalidResponse represents a failure from the server
type InvalidResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CurrentAgentResponse is the struct to represent the response from the HTTP server
// for the currently running agent
type CurrentAgentResponse struct {
	IsRunning bool              `json:"is_running"`
	ID        string            `json:"id"`
	Config    map[string]string `json:"config"`
}
