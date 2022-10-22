package http

// SetAgentRequest is a schema for a POST request to set the current agent
type SetAgentRequest struct {
	Agent string `json:"agent"`
}
