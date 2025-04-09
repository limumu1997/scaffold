package eventbus

type EventBus struct {
	data chan any
}

// New 创建一个事件总线实例
func New(bufferSize int) *EventBus {
	return &EventBus{
		data: make(chan any, bufferSize),
	}
}

// Publish 发布事件
func (e *EventBus) Publish(value any) {
	e.data <- value
}

// Subscribe 返回只读 channel，用于订阅事件
func (e *EventBus) Subscribe() <-chan any {
	return e.data
}
