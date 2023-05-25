package persistence

// Peer todo: Implement Address as []byte?
type Peer struct {
	Id        string `json:"id"`
	Tenant    string `json:"tenant"`
	Address   string `json:"address"`
	PublicKey []byte `json:"publicKey"`
}

// Tenant todo: Implement Network as CIDR object or []byte?
type Tenant struct {
	Id      string `json:"id"`
	Network string `json:"network"`
}
