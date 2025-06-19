package main

import (
	"cronparser/cronapp"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s \"<cron expression>\"", os.Args[0])
	}

	if err := run(os.Args[1]); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run(cronString string) error {
	schedule, err := cronapp.Parse(cronString)
	if err != nil {
		return fmt.Errorf("failed to parse cron expression: %w", err)
	}

	fmt.Println(schedule.String())
	return nil
}
