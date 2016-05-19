package cmd

import (
	"fmt"
	"gopkg.in/readline.v1"
)

type AskOptions struct {
	Candidates []string
	Validate   func(string) error
	Default    string
}

func AskForRequiredInput(prompt string, opts ...AskOptions) string {
	pcItems := []readline.PrefixCompleterInterface{}
	validate := func(item string) error { return nil }
	defaultValue := ""
	if len(opts) > 0 {
		o := opts[0]
		for _, c := range o.Candidates {
			pcItems = append(pcItems, readline.PcItem(c))
		}
		if o.Validate != nil {
			validate = o.Validate
		}
		if o.Default != "" {
			defaultValue = o.Default
		}
	}
	var completer = readline.NewPrefixCompleter(
		pcItems...,
	)

	fullPrompt := ""
	if defaultValue != "" {
		fullPrompt = fmt.Sprintf("%s(Default: %s)> ", prompt, defaultValue)
	} else {
		fullPrompt = fmt.Sprintf("%s> ", prompt)
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       fullPrompt,
		AutoComplete: completer,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		input, err := rl.Readline()
		if err == readline.ErrInterrupt {
			panic(err)
		}
		if err == nil { // or io.EOF
			if input == "" {
				input = defaultValue
			}
			r := validate(input)
			if r == nil {
				return input
			} else {
				fmt.Printf("%s is not valid: %v\n", input, r)
			}
		}
	}
}
