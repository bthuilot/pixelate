package http

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
