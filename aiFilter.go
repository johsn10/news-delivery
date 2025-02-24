package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func aiFilter(llm ollama.LLM, article Article, source Source) bool {
	ctx := context.Background()
	prefs := prefsToStr(source.Prefs)
	var passedFilter bool
	prompt := fmt.Sprintln("This article is your input\n\n",
		article.Description,
		"\n\nIs this article about either:",
		prefs,
		"\nAnswere with either {YES} or {NO}",
	)
	fmt.Println(prompt)
	completion, err := llm.Call(ctx, prompt,
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	switch completion {
	case "{YES}":
		passedFilter = true
	case "{NO}":
		passedFilter = false
	default:
		fmt.Println("Did not return valid result!")
		passedFilter = false
	}
	return passedFilter
}

func initLLM() (*ollama.LLM, error) {
	llm, err := ollama.New(ollama.WithModel("llama3.2"))
	if err != nil {
		return &ollama.LLM{}, err
	}
	return llm, nil
}

func prefsToStr(prefsList []string) string {
	var prefsStr string
	for _, pref := range prefsList {
		prefsStr += pref + ", "
	}
	return prefsStr
}
