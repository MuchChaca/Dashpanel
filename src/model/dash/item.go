package dash

// Item is an item
type Item struct {
	SessionID string `json:"-"`
	ID        int64  `json:"id,omitempty"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// // Apache is an apache server, apache2 or httpd
// type Apache struct {
// 	Version     string `json:"version"`
// 	ServiceName string `json:"service"`
// 	Status      bool   `json:"status"`
// }
