### Zero-Flow UI Library

Библиотека для потоковой передачи текстовых сигналов и декларативного логирования UI-компонентов. Спроектирована в рамках методологии **Zero-Collection Architecture**.

### 🛡️ Архитектурные гарантии (Zero-Collection)

*   **Полный отказ от коллекций**: В коде отсутствуют слайсы (`[]`), массивы (`[N]`) и хэш-таблицы (`map`).
*   **Исключение базовой математики**: Типы `int`, `uint`, `uintptr` и пакет `unsafe` полностью запрещены.
*   **Функциональные потоки**: Итерация строк и ленты событий UI построена на ленивых замыканиях и хвостовой рекурсии.
*   **Паттерн Bail-Out**: Точечные безопасные вызовы гарантируют завершение конвейера без генерации паник и исключений рантайма.

### 📊 Графические схемы архитектуры

### 1\. Архитектура текстового конвейера (StringIterator Closure Chain)

Каждый вызов `SafeNextChar()` не сдвигает числовой индекс, а переходит к следующему ленивому замыканию, изолирующему свой остаток строки в куче:

text

      [ String Stream Initializer ]
                   │
                   ▼
      ┌──────────────────────────────────────────────┐
      │ MakeStream("Go")                             │
      │  ├─ Current Char: "G"                        │
      │  └─ Next Stream ───► ┌───────────────────────┤
      └────────────────────► │ MakeStream("o")       │
                             │  ├─ Current Char: "o" │
                             │  └─ Next Stream ──────┼──► ┌─────────────────────┐
                             └───────────────────────┘    │ MakeStream("")      │
                                                          │  ├─ Char: ""        │
                                                          │  ├─ End: true       │
                                                          │  └─ [ Trigger Bail ]│
                                                          └─────────────────────┘
    

Используйте код с осторожностью.

### 2\. Стак событий UI и мониторинг бедствий (Deferred Action Chain)

Лента событий UI наслаивается в памяти в обратном порядке через декоратор `LogUIEvent`. При интерпретации рекурсия разворачивает стек, а флаг `disasterState` сквозным образом маркирует состояние системы:

text

      [ Емкость истории: Стек замыканий ]             [ Рантайм интерпретатора ]
      
      ┌───────────────────────────────────┐            ┌────────────────────────┐
      │ LogUIEvent (RetryButton, Clicked) │ ◄────────  │ Оболочка в режиме      │
      │  └─ Previous Node ────────────────┼───┐        │ КРИТИЧЕСКОГО БЕДСТВИЯ  │
      └───────────────────────────────────┘   │        │ (disasterState = true) │
                                              ▼        └───────────▲────────────┘
      ┌───────────────────────────────────┐                        │
      │ LogUIEvent (MainWindow, Crash)    │ ───────────────────────┘
      │  └─ Previous Node ────────────────┼───┐        [ Триггер аварии ]
      └───────────────────────────────────┘   │
                                              ▼
      ┌───────────────────────────────────┐            ┌────────────────────────┐
      │ LogUIEvent (MainWindow, Rendered) │ ◄────────  │ Оболочка стабильна     │
      │  └─ Previous Node: EndOfUI()      │            │ (disasterState = false)│
      └───────────────────────────────────┘            └────────────────────────┘
    

Используйте код с осторожностью.

### 📦 Установка

bash

    go get ://github.com
    

Используйте код с осторожностью.

### 🚀 Быстрый старт

Пример инициализации сквозного конвейера обработки сигналов и анализа интерфейса на предмет аварийных состояний (Crash/Panic):

go

    package main
    
    import (
    	"://github.com"
    )
    
    func main() {
    	// 1. Инициализация текстового сигнала
    	textSignal := zeroflowui.TextSignal{
    		Type:    zeroflowui.TextType,
    		Payload: "Потоковый текст\n",
    	}
    
    	// 2. Декларативная сборка истории UI (без append и списков)
    	uiTimeline := zeroflowui.EndOfUI()
    	uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "MainWindow", "Rendered")
    	uiTimeline = zeroflowui.LogUIEvent(uiTimeline, true,  zeroflowui.EventInteraction, "RefreshButton", "Clicked")
    	
    	// Симуляция бедствия оболочки
    	uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "MainWindow", "Crash")
    
    	// 3. Запуск конвейера декораторов
    	pipeline := zeroflowui.SystemPipelineDecorator{
    		Next: zeroflowui.TerminalProcessor{},
    	}
    
    	pipeline.Process(zeroflowui.NewTextFlow(textSignal), uiTimeline)
    }
    

Используйте код с осторожностью.

### 🗺️ Спецификация API

### Подсистема TextFlow

*   `MakeStream(text string) StringIterator` — Преобразует скалярную строку в ленивый посимвольный функциональный поток.
*   `NewTextFlow(sig TextSignal) *TextFlowLine` — Создает новое изолированное состояние линии конвейера.
*   `SafeNextChar() (*TextFlowLine, string)` — Точечный безопасный вызов. Возвращает следующий символ и сдвинутое состояние потока. При достижении конца строки переводит флаг `failed` в `true` (Bail-Out).

### Подсистема UI & Disaster Recovery

*   `EndOfUI() UIEventFlow` — Возвращает терминальный узел истории.
*   `LogUIEvent(prev UIEventFlow, isWidget bool, evType UIEvent, name, details string) UIEventFlow` — Декоратор, наслаивающий новое событие UI поверх предыдущего контекста выполнения.
*   Если в `details` передаются маркеры `"Crash"` или `"Panic"`, встроенный интерпретатор автоматически переключает контекст в аварийный режим, маркируя все последующие действия UI как произошедшие во время бедствия оболочки.

### 📄 Лицензия

MIT

# zero-flow-ui