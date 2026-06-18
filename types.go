package zeroflowui

// SignalType определяет тип передаваемых данных без использования чисел.
type SignalType bool
const TextType SignalType = true

// UIEvent разделяет действия пользователя и системные события.
type UIEvent bool
const (
	EventInteraction UIEvent = true  // Клик, ввод данных
	EventLifecycle   UIEvent = false // Инициализация, рендеринг, сбой
)

// TextSignal — атомарный контейнер для текстового потока.
type TextSignal struct {
	Type    SignalType
	Payload string
}

// UIStateDescriptor описывает состояние компонента интерфейса.
type UIStateDescriptor struct {
	IsWidget      bool
	EventType     UIEvent
	ComponentName string
	ActionDetails string
}

// UIEventFlow — функциональный поток истории событий UI.
type UIEventFlow func() (UIStateDescriptor, UIEventFlow, bool)
