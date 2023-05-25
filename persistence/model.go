package persistence

type Peer struct {
	Id        string `json:"id"`
	Address   string `json:"address"`
	PublicKey []byte `json:"publicKey"`
}

type Tenant struct {
	Id      string `json:"id"`
	Network string `json:"network"`
}

type OnboardingResponseDTO struct {
	Address string `json:"address"`
	Peers   []Peer `json:"peers"`
}
