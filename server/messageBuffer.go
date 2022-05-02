package server

type TailMessage struct {
	Line    string `json:"line"`
	Context string `json:"context,omitempty"`
}

type MessageBuffer struct {
	Size    int
	content []TailMessage
}

// NewMessageBuffer initializes a MessageBuffer with the given size.
func NewMessageBuffer(size int) (buffer MessageBuffer) {
	buffer = MessageBuffer{
		Size:    size,
		content: []TailMessage{},
	}
	return
}

// mabeInitEmpty initializes the MessageBuffer if it hasn't been done yet.
// It gives a default size of 10 if the size is 0 or less.
func (buff *MessageBuffer) maybeInitEmpty() {
	if buff.content == nil {
		buff.content = []TailMessage{}
	}
	if buff.Size <= 0 {
		buff.Size = 10
	}
}

// Add appends a TailMessage to the buffer, knocking off the first value if the size exceeds the size.
// Note that this will not *shorten* the buffer, it can only ever grow it.
func (buff *MessageBuffer) Add(entry TailMessage) {

	buff.maybeInitEmpty()

	if len(buff.content) >= buff.Size {
		for i := 0; i < len(buff.content)-1; i++ {
			buff.content[i] = buff.content[i+1]
		}
		buff.content[len(buff.content)-1] = entry
	} else {
		buff.content = append(buff.content, entry)
	}
}

// Get returns a copy of the TailMessage elements contained in the buffer.
func (buff *MessageBuffer) Get() []TailMessage {

	buff.maybeInitEmpty()

	content := make([]TailMessage, len(buff.content))
	copy(content, buff.content)
	return content
}

// Clear reinitializes the buffer, making it so it contains no elements.
func (buff *MessageBuffer) Clear() {
	buff.content = []TailMessage{}
}
