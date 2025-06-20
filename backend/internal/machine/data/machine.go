package machine

import "github.com/google/uuid"

type Machine struct {
	Id       string `json:"id"`
	OwnerId  string `json:"ownerid"`
	Vnc      bool   `json:"vnc"`
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Options struct {
	OwnerId  string
	Vnc      bool
	Username string
	Password string
}

func New(options Options) (*Machine, error) {
	return &Machine{
		Id:       uuid.NewString(),
		OwnerId:  options.OwnerId,
		Vnc:      options.Vnc,
		Username: options.Username,
		Password: options.Password,
	}, nil
}
