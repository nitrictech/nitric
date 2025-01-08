package deploytf

type EventGridSubscriber struct {
	Url                       *string `json:"url"`
	ActiveDirectoryAppIdOrUri *string `json:"active_directory_app_id_or_uri"`
	ActiveDirectoryTenantId   *string `json:"active_directory_tenant_id"`
	EventToken                *string `json:"event_token"`
}
