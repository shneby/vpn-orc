package persistence

type Peer struct {
	Id        string `json:"id"`
	VAddr     string `json:"virtualAddress"`
	RAddr     string `json:"-"`
	PublicKey []byte `json:"publicKey"`
}

type Tenant struct {
	Id      int    `json:"id"`
	Network string `json:"network"`
}

type OnboardingResponse struct {
	Address string `json:"address"`
	Peers   []Peer `json:"peers"`
}

type OnboardingRequest struct {
	Id        string `json:"id"`
	TenantId  int    `json:"tenantId"`
	PublicKey string `json:"publicKey"`
	Addr      string `json:"address"`
}
