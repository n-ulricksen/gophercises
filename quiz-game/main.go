package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type QuizQuestion struct {
	Question string
	Answer   string
}

var quizQuestions []*QuizQuestion
var isQuizOver bool
var correctCount int
var timer *time.Timer

// Default location for quiz questions/answers
const DEFAULT_QUIZFILE string = "problems.csv"
const DEFAULT_TIMELIMIT int = 20

// Command line flags
var flagFilename *string
var flagTimeLimit *int

func main() {
	parseFlags()

	getQuizQuestionsFromCSV()

	waitForKeyPress()

	startCountdown()

	activateQuiz()
}

func parseFlags() {
	flagFilename = flag.String("f", DEFAULT_QUIZFILE,
		"file containing quiz questions")
	flagTimeLimit = flag.Int("t", DEFAULT_TIMELIMIT,
		"quiz time limit in seconds")
	flag.Parse()
}

func getQuizQuestionsFromCSV() {
	// Get the list of questions/answers from CSV file
	filename := *flagFilename

	problems, err := parseCSV(filename)
	if err != nil {
		fmt.Printf("Error parsing file %v\n", filename)
		os.Exit(1)
	}

	for _, p := range problems {
		question := p[0]
		answer := p[1]

		quizQuestions = append(quizQuestions, &QuizQuestion{question, answer})
	}
}

func waitForKeyPress() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Press <ENTER> to begin quiz. You will have %v seconds.",
		*flagTimeLimit)
	for {
		scanner.Scan()
		break
	}
}

func startCountdown() {
	timeLimit := time.Duration(*flagTimeLimit) * time.Second
	timer = time.NewTimer(timeLimit)
}

func activateQuiz() {
	for _, quizQuestion := range quizQuestions {
		fmt.Println(quizQuestion.Question)

		answerCh := make(chan string)
		go func() {
			userAnswer := getUserAnswer()
			answerCh <- userAnswer
		}()

		select {
		case <-timer.C:
			displayResults()
			return

		case userAnswer := <-answerCh:
			if userAnswer == quizQuestion.Answer {
				correctCount++
			}
			fmt.Println()
		}
	}
}

func parseCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return rows, nil
}

func getUserAnswer() string {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		if len(input) == 0 {
			continue
		}

		return input
	}
}

func displayResults() {
	fmt.Printf("\nQuiz Results: %v correct out of %v\n", correctCount,
		len(quizQuestions))
}
