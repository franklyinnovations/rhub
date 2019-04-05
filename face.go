package rhub

import (
	"encoding/json"
)

type Handler func(msg *RedisHubMessage)
type Filter func(msg *RedisHubMessage, next func())
type HandlerWs func(message *ClientHubMessage)

//IHub like chat room
type IHub interface {
	//Get hub's id
	Id() interface{}
	//before new client join
	BeforeJoin(callback func(client IClient) error)
	AfterJoin(callback func(client IClient))
	//On client leave
	BeforeLeave(callback func(client IClient))
	AfterLeave(callback func(client IClient))
	BeforeWsMsg(callback func(msg *ClientHubMessage) bool)
	//Add filter
	// Use(filter Filter)
	// Attach an event handler function
	On(subject string, handler Handler)
	OnWs(subject string, handler HandlerWs)
	// OnLocal(subject string, handler HandlerLocal)
	//simulate client send msg
	// Emit(msg *ClientHubMessage)
	//Dettach an event handler function
	Off(subject string, handler Handler)
	OffWs(subject string, handler HandlerWs)
	SendRedisRaw(msg *HubMessageIn, props map[string]interface{})
	SendRedis(subject string, data interface{}, props map[string]interface{})
	//Send message to all clients
	SendWsAll(subject string, message interface{})
	SendWs(subject string, message interface{}, receivers []IClient)
	EchoWs(msg *ClientHubMessage)
	Close()
	// CloseMessageLoop()
	// SetSelf(self IHub)
	Run()
	UnregisterChan() chan IClient
	RegisterChan() chan IClient
	MessageChan() chan *ClientHubMessage
	// MessageLocalChan() chan<- MessageLocal
	CloseChan() chan struct{}
	Clients() map[IClient]bool

	OnTick(func(int))
	GetSeconds() int
	SendWsClient(client IClient, subject string, message interface{})
	SendWsBytes(client IClient, bs []byte)
}
type IClient interface {
	Close()
	// Send() chan []byte
	SendChan() chan []byte
	// Send(subject string, msg interface{})
	Hub() IHub
	WritePump()
	ReadPump()
	GetProps() map[string]interface{}
	// Get(key interface{}) interface{}
	// Set(key, value interface{})
	NewClientHubMessage(data []byte) (*ClientHubMessage, error)
	GetClient() *Client
}

type IFilters interface {
	Do(fn Handler)
}

//find callbacks in hub
type IRoute interface {
	Route(subject string) []Handler
	// Attach an event handler function
	On(subject string, handler Handler)
	//Dettach an event handler function
	Off(subject string, handler Handler)
}
type MessageLocal struct {
	Subject string
	Data    interface{}
	Error   error
}

//send Message should have this format
type HubMessageOut struct {
	Subject string
	Data    interface{}
}
type HubMessageIn struct {
	Subject string
	Data    *json.RawMessage
}

type RedisHubMessage struct {
	*HubMessageIn
	Props map[string]interface{}
}
type ClientHubMessage struct {
	*HubMessageIn
	Client IClient
}

func NewClientHubMessage(subject string, data interface{}, client IClient) *ClientHubMessage {
	r := &ClientHubMessage{Client: client}
	r.HubMessageIn = &HubMessageIn{Subject: subject}
	r.HubMessageIn.SetData(data)
	return r
}
func (m *HubMessageIn) SetData(data interface{}) error {
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r := json.RawMessage(bs)
	m.Data = &r
	return nil
}
func (m *HubMessageIn) Decode(obj interface{}) error {
	if nil != m.Data {
		bs, err := m.Data.MarshalJSON()
		err = json.Unmarshal(bs, &obj)
		return err
	}
	return nil
}

func (m HubMessageOut) Encode() ([]byte, error) {
	return json.Marshal(m)
}
