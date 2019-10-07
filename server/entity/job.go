package entity

// Payload defines the structure of the extra payload in a job request
type Payload map[string]interface{}

// JobRequest a request sent to the /jobs endpoint
type JobRequest struct {
	Name    string  `json:"name" validate:"required"`
	Payload Payload `json:"payload"`
}

// JobResponse defines a response returned from the /jobs endpoint
type JobResponse struct {
	Message string `json:"message"`
}
