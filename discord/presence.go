package discord
import "sync"
type Presence struct {
	mu     sync.Mutex
	client *Client
	appID  string
}
func NewPresence() *Presence {
	return &Presence{}
}
func (p *Presence) Set(appID string, a Activity) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.client != nil && p.client.appID != appID {
		_ = p.client.Close()
		p.client = nil
	}
	if p.client == nil {
		c, err := Login(appID)
		if err != nil {
			return err
		}
		p.client = c
		p.appID = appID
	}
	return p.client.SetActivity(a)
}
func (p *Presence) Clear() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.client == nil {
		return nil
	}
	return p.client.ClearActivity()
}
func (p *Presence) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.client == nil {
		return nil
	}
	return p.client.Close()
}
