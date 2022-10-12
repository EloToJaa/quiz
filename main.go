package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format question,answer")
	timeLimit := flag.Int("limit", 30, "time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Printf("Failed to open file: %s", *csvFilename)
		os.Exit(1)
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Failed to parse the provided csv file")
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	problemLoop(problems, timer)
}

type Problem struct {
	question string
	answer   string
}

func parseLines(lines [][]string) []Problem {
	ret := make([]Problem, len(lines))
	for i, line := range lines {
		ret[i] = Problem{
			question: strings.TrimSpace(line[0]),
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func problemLoop(problems []Problem, timer *time.Timer) {
	correct := 0
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d/%d\n", correct, len(problems))
			return

		case answer := <-answerCh:
			if answer == problem.answer {
				correct++
			}
		}
	}
	fmt.Printf("You scored %d/%d\n", correct, len(problems))
}
