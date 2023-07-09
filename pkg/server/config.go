package server

type Config struct {
	// Address is the IP:port combo to host on
	Address string
}

func NewDefaultConfig() Config {
	return Config{
		Address: "0.0.0.0:8341",
	}
}
