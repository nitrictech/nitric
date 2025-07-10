package details

type AuthDetailsService interface {
	GetWorkOSDetails() (*WorkOSDetails, error)
}

type WorkOSDetails struct {
	ClientID    string `json:"client_id"`
	ApiHostname string `json:"api_hostname"`
}
