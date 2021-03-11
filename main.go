package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var count int

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(2)
}
func main() {

	csvFileName := flag.String("csv", "problems.csv", "A csv file with format question,answer")
	timeLimit := flag.Int("limit", 5, "Time limit for quiz in seconds")
	flag.Parse()
	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Could not open the file  %s\n", *csvFileName))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Could not read the file ")
	}
	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	for i, p := range problems {
		fmt.Printf("Problem #%d : %s\n", i+1, p.question)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("Your score is %d out of %d\n", count, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.answer {
				count++
			}
		}
	}
	fmt.Printf("Your score is %d out of %d\n", count, len(problems))

}
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			line[0],
			strings.TrimSpace(line[1]),
		}
	}
	return ret
}
