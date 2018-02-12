package queue

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	// MetaID represents a unique message id
	MetaID = "id"
	// MetaCreated represents the creation timestamp
	MetaCreated = "created"
	// MetaSource represents the input source (e.g. file, bus, http, ...)
	MetaSource = "source"
	// MetaFormat represents the document/fileformat (e.g. rdi, xml, csv)
	MetaFormat = "format"
	// MetaBatch carries the id of documents belonging to the same input stream
	MetaBatch = "batch"
)

// Message represents the information being passed around by the different
// docproc services.
type Message struct {
	Metadata map[string]interface{} `json:"metadata"`
	Content  map[string]interface{} `json:"content"`
}

// Processor implementations are used to perform operations based on a Message.
type Processor interface {
	Name() string
	Process(msg *Message) error
}

// NewMessage creates a new Message with the passed content.
// It is guaranteed that the newly created Message's Metadata contain at least
// a MetaID and MetaCreated field.
func NewMessage(content map[string]interface{}) *Message {
	var id string
	if uid, err := uuid.NewRandom(); err != nil {
		id = fmt.Sprint(time.Now().Unix())
	} else {
		id = uid.String()
	}
	metadata := map[string]interface{}{
		MetaID:      id,
		MetaCreated: time.Now(),
	}
	return &Message{
		Metadata: metadata,
		Content:  content,
	}
}

// Clear resets the Content section of the Message.
func (msg *Message) Clear() {
	msg.Content = make(map[string]interface{})
}

// ToJSON returns a JSON representation of the Message
func (msg *Message) ToJSON() ([]byte, error) {
	return json.Marshal(msg)
}

// FromJSON will initialize the Message with the JSON representation of a
// Message.
func (msg *Message) FromJSON(data []byte) error {
	return json.Unmarshal(data, msg)
}

// MsgFromJSON returns a new Message from the passed in JSON representation
func MsgFromJSON(data []byte) (*Message, error) {
	msg := &Message{}
	if err := msg.FromJSON(data); err != nil {
		return nil, err
	}
	return msg, nil
}
