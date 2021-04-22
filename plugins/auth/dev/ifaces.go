package auth_service

type ScribbleIface interface {
	Read(string, string, interface{}) error
	ReadAll(string) ([]string, error)
	Write(string, string, interface{}) error
	Delete(string, string) error
}
