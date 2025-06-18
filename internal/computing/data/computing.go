package computing

type Computing struct {
	Id      string `json:"id"`
	OwnerId string `json:"ownerId"`
	Vnc     bool   `json:"vnc"`
	Address string `json:"address"`
	Port    string `json:"port"`
}
