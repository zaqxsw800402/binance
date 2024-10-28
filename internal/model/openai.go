package model

import "sync"

type Prompt struct {
	Question string
	Data     string
}

func NewPrompt(question string, data []byte) *Prompt {
	return &Prompt{Question: question, Data: string(data)}
}

type AllPrompt struct {
	mu      sync.Mutex
	Prompts []*Prompt
}

func NewAllPrompt() *AllPrompt {
	return &AllPrompt{}
}

func (a *AllPrompt) AddPrompt(prompt *Prompt) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Prompts = append(a.Prompts, prompt)
}
