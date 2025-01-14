package handlers

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User          string `json:"user"`
	Authenticated bool   `json:"authenticated"`
}

type DownloadRequest struct {
	User     string `json:"user"`
	Category string `json:"category"`
	Year     int    `json:"year"`
}

//type DownloadResponse struct{}
