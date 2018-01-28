package dash

import "sync"

// Service is our service .... -_-
type Service interface {
	Get(owner string) []Item
	Save(owner string, newItems []Item) error
}

// MemoryService is a struct
type MemoryService struct {
	// key = session id, value the list of dash items that this session id has.
	items map[string][]Item
	// protected by locker for concurrent access
	mu sync.RWMutex
}

// NewMeMoryService create a new MeMoryService
func NewMeMoryService() *MemoryService {
	return &MemoryService{
		items: make(map[string][]Item, 0),
	}
}

// Save saves the MemoryService -_- ?
func (s *MemoryService) Save(sessionOwner string, newItems []Item) error {
	var prevID int64
	for i := range newItems {
		if newItems[i].ID == 0 {
			newItems[i].ID = prevID
			prevID++
		}
	}

	s.mu.Lock()
	s.items[sessionOwner] = newItems
	s.mu.Unlock()
	return nil
}
