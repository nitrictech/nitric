package details

type AuthDetailsService interface {
	GetWorkOSDetails() (*WorkOSDetails, error)
}

type WorkOSDetails struct {
	ClientID    string `json:"clientId"`
	ApiHostname string `json:"apiHostname"`
}
