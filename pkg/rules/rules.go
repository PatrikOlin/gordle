package rules

type Rules struct {
	MaxGuesses int
	RuleExp    string
}

var RuleExp = `Rules:
		- You have 6 guesses
		- You can only guess using 5 letters words
		- You can only guess using lowercase letters
		- Green means correct at the right position
		- Yellow means correct at the wrong position
	Give a filename in argument to have more words

	Good luck!
`

func Get() Rules {
	return Rules{
		MaxGuesses: 6,
		RuleExp:    RuleExp,
	}
}
