package persistence

type Peer struct {
	Id        string `json:"id"`
	Tenant    string `json:"tenant"`
	Address   string `json:"address"`
	PublicKey []byte `json:"publicKey"`
}

type Tenant struct {
	Id      string `json:"id"`
	Network string `json:"network"`
}
