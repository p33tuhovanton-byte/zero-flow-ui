package zeroflowui

import "fmt"

// UIListener — функциональный тип для прослушивания событий.
// Вызывается при срабатывании действия.
type UIListener func(desc UIStateDescriptor)

// ActionBuilder инкапсулирует в себе данные о собираемом действии 
// и цепочку прослушивателей (вместо слайса функций).
type ActionBuilder struct {
	descriptor UIStateDescriptor
	listener   UIListener // Цепочка замыканий-слушателей
}

// CreateAction — точка входа для сборщика действия интерфейса.
func CreateAction() ActionBuilder {
	return ActionBuilder{}
}

// SetComponent задает имя и тип компонента.
func (ab ActionBuilder) SetComponent(name string, isWidget bool) ActionBuilder {
	ab.descriptor.ComponentName = name
	ab.descriptor.IsWidget = isWidget
	return ab
}

// SetEvent задает тип события и его текстовое описание.
func (ab ActionBuilder) SetEvent(evType UIEvent, details string) ActionBuilder {
	ab.descriptor.EventType = evType
	ab.descriptor.ActionDetails = details
	return ab
}

// Listen добавляет прослушиватель для интерфейса.
// Новые слушатели наслаиваются друг на друга без использования списков.
func (ab ActionBuilder) Listen(newListener UIListener) ActionBuilder {
	if ab.listener == nil {
		ab.listener = newListener
		return ab
	}

	// Сохраняем старого слушателя во фрейме замыкания
	oldListener := ab.listener
	ab.listener = func(desc UIStateDescriptor) {
		oldListener(desc) // Сначала уведомляем предыдущего
		newListener(desc) // Затем нового
	}
	return ab
}

// Emit производит "выброс" действия: запускает прослушивание 
// и рекурсивно возвращает новое состояние ленты UI (UIEventFlow).
func (ab ActionBuilder) Emit(currentTimeline UIEventFlow) UIEventFlow {
	// 1. Активируем прослушивание для интерфейса
	if ab.listener != nil {
		ab.listener(ab.descriptor)
	}

	// 2. Добавляем собранное действие в общую функциональную историю UI
	return LogUIEvent(
		currentTimeline,
		ab.descriptor.IsWidget,
		ab.descriptor.EventType,
		ab.descriptor.ComponentName,
		ab.descriptor.ActionDetails,
	)
}
