package tokenizer

type Options struct {
	KeepPunctuation bool
}

func DefaultOptions() Options {
	return Options{
		KeepPunctuation: true,
	}
}
