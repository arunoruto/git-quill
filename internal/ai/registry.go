package ai

import "fmt"

var providers = []Provider{
	Ollama{},
	Opencode{},
	Gemini{},
	Copilot{},
}

func GetAvailableProviders() []Provider {
	var available []Provider
	for _, p := range providers {
		if p.IsAvailable() {
			available = append(available, p)
		}
	}
	return available
}

func GetProviderByName(name string) (Provider, error) {
	for _, p := range providers {
		if p.Name() == name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("provider '%s' not found", name)
}
