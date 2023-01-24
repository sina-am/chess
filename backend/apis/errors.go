package apis

type HTTPError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Details    string `json:"details"`
}

func (e *HTTPError) Error() string {
	return e.Details
}
