package config

import (
	_ "embed"
	"encoding/json"
	"log"
)

//go:embed json/prompt.json
var promptsData []byte

type PromptsConfig struct {
	Prompts []Prompt `json:"prompts"`
}

type Prompt struct {
	Name            string      `json:"name"`
	System          string      `json:"system"`
	Function        []function  `json:"functions"`
	FewShotExamples interface{} `json:"few_shot_examples"`
}

type function struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parameters  parameters `json:"parameters"`
	Required    []string   `json:"required"`
}

type parameters struct {
	Type       string      `json:"type"`
	Properties interface{} `json:"properties"`
}

func GetPrompts() map[string]Prompt {

	var promptsConfig PromptsConfig
	err := json.Unmarshal(promptsData, &promptsConfig)
	if err != nil {
		log.Fatal(err)
	}

	promptsMap := make(map[string]Prompt)
	for _, val := range promptsConfig.Prompts {
		promptsMap[val.Name] = val
	}

	return promptsMap
}
