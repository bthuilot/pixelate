package requests

// SetScreenRequest is a schema for a POST request to set the current agent
type SetScreenRequest struct {
	Screen string `json:"screen"`
}
