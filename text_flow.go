package zeroflowui

// TextFlowLine инкапсулирует состояние продвижения по текстовому сигналу.
type TextFlowLine struct {
	signal TextSignal
	stream StringIterator
	failed bool
}

// NewTextFlow — конструктор безопасной линии текстового потока.
func NewTextFlow(sig TextSignal) *TextFlowLine {
	return &TextFlowLine{
		signal: sig,
		stream: MakeStream(sig.Payload),
		failed: false,
	}
}

// SafeNextChar выполняет точечный безопасный вызов для чтения одного символа.
func (f *TextFlowLine) SafeNextChar() (*TextFlowLine, string) {
	if f.failed || f.stream == nil {
		return f, ""
	}

	char, nextStream, isEnd := f.stream()
	if isEnd {
		f.failed = true
		return f, ""
	}

	return &TextFlowLine{signal: f.signal, stream: nextStream, failed: false}, char
}
