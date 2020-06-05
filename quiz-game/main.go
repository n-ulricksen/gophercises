package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
)

type QuizQuestion struct {
	Question string
	Answer   string
}

var quizQuestions []*QuizQuestion
var correctCount int

// Default location for quiz questions/answers
const DEFAULT_QUIZFILE string = "problems.csv"

// Command line flags
var flagFilename *string

func main() {
	parseFlags()

	getQuizQuestionsFromCSV()

	activateQuiz()

	displayResults()
}

func parseFlags() {
	flagFilename = flag.String("f", "", "file containing quiz questions")
	flag.Parse()
}

func getQuizQuestionsFromCSV() {
	// Get the list of questions/answers from CSV file
	filename := DEFAULT_QUIZFILE
	if *flagFilename != "" {
		filename = *flagFilename
	}

	problems, err := parseCSV(filename)
	if err != nil {
		fmt.Println("Empty file...")
	}

	for _, p := range problems {
		question := p[0]
		answer := p[1]

		quizQuestions = append(quizQuestions, &QuizQuestion{question, answer})
	}
}

func activateQuiz() {
	for _, quizQuestion := range quizQuestions {
		fmt.Println(quizQuestion.Question)

		userAnswer := getUserAnswer()
		if userAnswer == quizQuestion.Answer {
			correctCount++
		}
		fmt.Println()
	}
}

func displayResults() {
	fmt.Printf("Quiz Results: %v correct out of %v\n", correctCount,
		len(quizQuestions))
}

func parseCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
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
