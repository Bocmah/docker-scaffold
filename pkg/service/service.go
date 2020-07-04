package service

type PHP struct {
	Version    string
	Extensions []string
}

type Nginx struct {
	Port               int
	ServerName         string
	FastCGIPassPort    int
	FastCGIReadTimeout int
}

type NodeJS struct {
	Version string
}
