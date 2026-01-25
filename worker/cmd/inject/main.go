package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/OkadaSatoshi/codingworker/worker/internal/sqs"
)

func main() {
	// Flags
	repo := flag.String("repo", "", "Repository (e.g., owner/repo)")
	issue := flag.Int("issue", 0, "Issue number")
	title := flag.String("title", "", "Task title")
	body := flag.String("body", "", "Task body")
	jsonFile := flag.String("json", "", "JSON file containing message")
	output := flag.String("output", "", "Output file path (default: stdout)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: inject [options]\n\n")
		fmt.Fprintf(os.Stderr, "Generate a test message JSON for the CodingWorker.\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  inject -repo owner/repo -issue 1 -title \"Create hello.go\" -body \"Create a hello world program\"\n")
		fmt.Fprintf(os.Stderr, "  inject -json message.json\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	var msg *sqs.Message

	if *jsonFile != "" {
		// Load from JSON file
		data, err := os.ReadFile(*jsonFile)
		if err != nil {
			log.Fatalf("Failed to read JSON file: %v", err)
		}
		msg = &sqs.Message{}
		if err := json.Unmarshal(data, msg); err != nil {
			log.Fatalf("Failed to parse JSON: %v", err)
		}
	} else {
		// Create from flags
		if *repo == "" || *issue == 0 || *title == "" {
			flag.Usage()
			os.Exit(1)
		}

		msg = &sqs.Message{
			IssueNumber: *issue,
			Repository:  *repo,
			Title:       *title,
			Body:        *body,
			Labels:      []string{sqs.LabelTrigger},
			CreatedAt:   time.Now().Format(time.RFC3339),
		}
	}

	// Generate JSON output
	data, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
	}

	if *output != "" {
		if err := os.WriteFile(*output, data, 0644); err != nil {
			log.Fatalf("Failed to write output file: %v", err)
		}
		fmt.Printf("Message written to %s\n", *output)
	} else {
		fmt.Println(string(data))
	}
}
