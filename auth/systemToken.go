package auth

// SystemToken contains data for authentication against the Cloudogu Ecosystem Backend.
type SystemToken struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}
