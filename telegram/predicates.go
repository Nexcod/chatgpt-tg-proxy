package telegram

type Predicate func(Update) bool

func Command(command string) Predicate {
	return func(update Update) bool {
		return update.Message.Text == "/"+command
	}
}

func AnyTextMessage() Predicate {
	return func(update Update) bool {
		return update.Message.Text != ""
	}
}
