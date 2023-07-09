package display

type DropDownInput struct {
}

type TextFieldInput struct {
}

type UserInput interface {
	GenerateHTML() string
}

type Config = map[string]UserInput
