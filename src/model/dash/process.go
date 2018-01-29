package dash

// // Process

// Process is a process we can start or stop on the server
type Process interface {
	Start() error
	Stop() error
	GetPID() int
	GetVersion() string
	GetStatus() bool
}

// Apache is an apache server, apache2 or httpd
type Apache struct {
	Name    string `json:"service"`
	Version string `json:"version"`
	Status  bool   `json:"status"`
}

// Processus is just a simple processus
type Processus struct {
	Name    string `json:"name"`
	PID     int    `json:"pid"`
	Version string `json:"version"`
	Status  bool   `json:"status"`
}

// TreeProcess represents all processussesssss
type TreeProcess struct {
	Tree []Processus `json:"tree"`
}

// Start the apache service
func (p *Processus) Start() error {
	return nil
}

// Stop the apache service
func (p *Processus) Stop() error {
	// new command
	// cmd := exec.Command("kill", p.PID)

	return nil
}

// GetPID get the pid of the apache service
func (p *Processus) GetPID() int {
	return 0
}

// GetVersion return the version of the apache service
func (p *Processus) GetVersion() string {
	return p.Version
}

// GetStatus return false if the process is down, true otherwise
func (p *Processus) GetStatus() bool {
	return false
}

// // // Service
// type Service interface {
// 	Get() []Process
// }
//
// // ServerService is a struct
// type ServerService struct {
// 	items map[string]
// }
