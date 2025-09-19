package schema

type AuthorizerResponse struct {
	Status      string `json:"status"`
	AuthorizeID string `json:"authorize_id,omitempty"`
	Error       string `json:"errors,omitempty"`
	Warning     string `json:"warning,omitempty"`
}
