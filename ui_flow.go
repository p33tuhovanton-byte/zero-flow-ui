package zeroflowui

// EndOfUI возвращает терминатор истории UI (Bail-Out узел).
func EndOfUI() UIEventFlow {
	return func() (UIStateDescriptor, UIEventFlow, bool) {
		return UIStateDescriptor{}, nil, true
	}
}

// LogUIEvent — декоратор для рекурсивного наращивания истории UI событий.
func LogUIEvent(prev UIEventFlow, isWidget bool, evType UIEvent, name, details string) UIEventFlow {
	return func() (UIStateDescriptor, UIEventFlow, bool) {
		desc := UIStateDescriptor{
			IsWidget:      isWidget,
			EventType:     evType,
			ComponentName: name,
			ActionDetails: details,
		}
		return desc, prev, false
	}
}
