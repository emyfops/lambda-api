package setting

type Server struct {
	Host         string // Host is the host of the server
	Port         int    // Port is the port of the server (default 8080)
	ReadTimeout  int    // ReadTimeout is the maximum time for reading the entire request
	WriteTimeout int    // WriteTimeout is the  maximum time for writing the response
}

type Party struct {
	MaxPlayers int // MaxPlayers is the maximum of slots in a party
}

var DefaultServer = &Server{
	Host:         "localhost",
	Port:         8080,
	ReadTimeout:  5000,
	WriteTimeout: 5000,
}

var DefaultParty = &Party{
	MaxPlayers: 16,
}

func init() {} // TODO: Read config
