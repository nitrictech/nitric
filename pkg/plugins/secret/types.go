package secret

// Secret - Represents a container for secret versions
type Secret struct {
	Name string
}

// SecretVersion - A version of a secret
type SecretVersion struct {
	Secret  *Secret
	Version string
}

// SecretAccessResponse - Return value for a secret access request
type SecretAccessResponse struct {
	SecretVersion *SecretVersion
	Value         []byte
}

// SecretPutResponse - Return value for a secret put request
type SecretPutResponse struct {
	SecretVersion *SecretVersion
}
