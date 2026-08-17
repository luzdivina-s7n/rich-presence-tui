package discord
import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)
const (
	opHandshake uint32 = 0
	opFrame     uint32 = 1
	opClose     uint32 = 2
	opPing      uint32 = 3
	opPong      uint32 = 4
)
var errConnectionLost = errors.New("discord connection lost")
type Client struct {
	conn   net.Conn
	appID  string
	pid    int
	mu     sync.Mutex
	frames chan ipcFrame
	done   chan struct{}
	once   sync.Once
}
type ipcFrame struct {
	op   uint32
	data []byte
}
func Login(appID string) (*Client, error) {
	appID = strings.TrimSpace(appID)
	if appID == "" {
		return nil, errors.New("application id is empty")
	}
	conn, err := dialPipe()
	if err != nil {
		return nil, err
	}
	c := &Client{
		conn:   conn,
		appID:  appID,
		pid:    os.Getpid(),
		frames: make(chan ipcFrame, 8),
		done:   make(chan struct{}),
	}
	go c.readLoop()
	if err := c.handshake(); err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}
func (c *Client) AppID() string {
	return c.appID
}
func (c *Client) Connected() bool {
	select {
	case <-c.done:
		return false
	default:
		return true
	}
}
func (c *Client) readLoop() {
	for {
		op, data, err := readFrame(c.conn)
		if err != nil {
			c.once.Do(func() { close(c.done) })
			return
		}
		switch op {
		case opPing:
			_ = c.send(opPong, nil)
		case opPong:
		default:
			select {
			case c.frames <- ipcFrame{op: op, data: data}:
			case <-c.done:
				return
			}
		}
	}
}
func (c *Client) send(op uint32, payload interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var data []byte
	if payload != nil {
		var err error
		data, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	}
	buf := make([]byte, 8+len(data))
	binary.LittleEndian.PutUint32(buf[0:4], op)
	binary.LittleEndian.PutUint32(buf[4:8], uint32(len(data)))
	copy(buf[8:], data)
	if _, err := c.conn.Write(buf); err != nil {
		c.once.Do(func() { close(c.done) })
		return err
	}
	return nil
}
func (c *Client) handshake() error {
	payload := map[string]interface{}{
		"v":         1,
		"client_id": c.appID,
	}
	if err := c.send(opHandshake, payload); err != nil {
		return fmt.Errorf("discord handshake: %w", err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		select {
		case f := <-c.frames:
			var evt struct {
				Evt  string          `json:"evt"`
				Cmd  string          `json:"cmd"`
				Data json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(f.data, &evt); err != nil {
				continue
			}
			if evt.Evt == "READY" || evt.Cmd == "READY" {
				return nil
			}
			if msg := decodeError(evt.Data); msg != "" {
				return errors.New(msg)
			}
		case <-c.done:
			return errConnectionLost
		case <-time.After(100 * time.Millisecond):
		}
	}
	return errors.New("discord handshake timeout")
}
func (c *Client) SetActivity(a Activity) error {
	return c.command("SET_ACTIVITY", map[string]interface{}{
		"pid":      c.pid,
		"activity": a,
	})
}
func (c *Client) ClearActivity() error {
	return c.command("SET_ACTIVITY", map[string]interface{}{
		"pid":      c.pid,
		"activity": nil,
	})
}
func (c *Client) Close() error {
	c.once.Do(func() { close(c.done) })
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
func (c *Client) command(cmd string, args map[string]interface{}) error {
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	payload := map[string]interface{}{
		"cmd":   cmd,
		"args":  args,
		"nonce": nonce,
	}
	if err := c.send(opFrame, payload); err != nil {
		return err
	}
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		select {
		case f := <-c.frames:
			var resp struct {
				Nonce string          `json:"nonce"`
				Evt   string          `json:"evt"`
				Cmd   string          `json:"cmd"`
				Data  json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(f.data, &resp); err != nil {
				continue
			}
			if resp.Nonce != nonce {
				continue
			}
			if msg := decodeError(resp.Data); msg != "" {
				return errors.New(msg)
			}
			return nil
		case <-c.done:
			return errConnectionLost
		case <-time.After(100 * time.Millisecond):
		}
	}
	return errors.New("timeout waiting for discord response")
}
func decodeError(data json.RawMessage) string {
	var e struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(data, &e); err != nil || e.Message == "" {
		return ""
	}
	if e.Code == 0 {
		return ""
	}
	return fmt.Sprintf("discord: %s", e.Message)
}
func readFrame(r io.Reader) (uint32, []byte, error) {
	header := make([]byte, 8)
	if _, err := io.ReadFull(r, header); err != nil {
		return 0, nil, err
	}
	op := binary.LittleEndian.Uint32(header[0:4])
	n := binary.LittleEndian.Uint32(header[4:8])
	if n > 1<<20 {
		return 0, nil, errors.New("discord frame too large")
	}
	data := make([]byte, n)
	if _, err := io.ReadFull(r, data); err != nil {
		return 0, nil, err
	}
	return op, data, nil
}
