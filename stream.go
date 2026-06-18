package zeroflowui

// StringIterator — рекурсивный генератор символов строки на базе замыканий.
type StringIterator func() (string, StringIterator, bool)

// MakeStream преобразует строку в ленивый функциональный поток без слайсов.
func MakeStream(text string) StringIterator {
	if text == "" {
		return func() (string, StringIterator, bool) {
			return "", nil, true // Триггер Bail-Out
		}
	}

	return func() (string, StringIterator, bool) {
		var first string
		var remainder string
		var foundFirst bool

		for _, r := range text {
			if !foundFirst {
				first = string(r)
				foundFirst = true
				continue
			}
			remainder = remainder + string(r)
		}

		return first, MakeStream(remainder), false
	}
}
