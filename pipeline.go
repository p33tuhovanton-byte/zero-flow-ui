package zeroflowui

import "fmt"

// SignalProcessor задает интерфейс для цепочки декораторов (Middleware).
type SignalProcessor interface {
	Process(textFlow *TextFlowLine, uiFlow UIEventFlow)
}

// TerminalProcessor завершает выполнение конвейера.
type TerminalProcessor struct{}
func (t TerminalProcessor) Process(textFlow *TextFlowLine, uiFlow UIEventFlow) {}

// SystemPipelineDecorator — главный обрабатывающий декоратор.
type SystemPipelineDecorator struct {
	Next SignalProcessor
}

func (spd SystemPipelineDecorator) Process(textFlow *TextFlowLine, uiFlow UIEventFlow) {
	if textFlow.signal.Type == TextType {
		fmt.Println("--- [ОБРАБОТКА ТЕКСТОВОГО СИГНАЛА FLOWLINE] ---")
		spd.transmitTextLoop(textFlow)
	}

	fmt.Println("\n--- [АНАЛИЗ ПОЛНОГО ЦИКЛА UI И СОСТОЯНИЯ ОБОЛОЧКИ] ---")
	spd.interpretUILoop(uiFlow, false)

	spd.Next.Process(textFlow, uiFlow)
}

func (spd SystemPipelineDecorator) transmitTextLoop(f *TextFlowLine) {
	nextFlow, char := f.SafeNextChar()
	if nextFlow.failed {
		return
	}
	fmt.Print(char)
	spd.transmitTextLoop(nextFlow)
}

func (spd SystemPipelineDecorator) interpretUILoop(flow UIEventFlow, disasterState bool) {
	if flow == nil {
		return
	}

	descriptor, nextFlow, isEnd := flow()
	if isEnd {
		return
	}

	if descriptor.ActionDetails == "Crash" || descriptor.ActionDetails == "Panic" {
		disasterState = true
	}

	spd.interpretUILoop(nextFlow, disasterState)

	fmt.Print("[UI ")
	if descriptor.IsWidget {
		fmt.Print("WIDGET")
	} else {
		fmt.Print("SCREEN")
	}
	fmt.Print("] Компонент: ", descriptor.ComponentName)

	if descriptor.EventType == EventInteraction {
		fmt.Print(" | Клиентский ввод: ", descriptor.ActionDetails)
	} else {
		fmt.Print(" | Системный статус: ", descriptor.ActionDetails)
	}

	if disasterState {
		fmt.Print(" -> !!! [ОБНАРУЖЕНО БЕДСТВИЕ ОБОЛОЧКИ (Crash/Panic State)]")
	}
	fmt.Println()
}
