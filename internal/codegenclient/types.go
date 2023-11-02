package codegenclient

type InputStringArray []string

func (i *InputStringArray) String() string {
	return ""
}

func (i *InputStringArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type AppFlags struct {
	InputFile InputStringArray
}
