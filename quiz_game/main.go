package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

// problem statement :
// make a quiz program that read csv file as a mathematic question with given format : question,answer
// print question sequently and read user input to get the user answer
// generate report from user answer
// print quiz result at the end of the program
// also given a timer, stop program and print user result when timer is over

func main() {
	userInput := readUserInput()
	problems := getProblems(*userInput.filename)
	testResult := runQuizes(problems, *userInput.timelimit)
	generateReport(testResult)
}

type userInput struct {
	filename  *string
	timelimit *int
}

type problem struct {
	question string
	answer   string
}

type testResult struct {
	totalQuestion int
	correctAnswer int
	details       []reportDetail
}

func newTestResult(question int) testResult {
	return testResult{
		totalQuestion: question,
		correctAnswer: 0,
		details:       make([]reportDetail, 0, question),
	}
}

func (tr *testResult) addDetail(detail reportDetail) {
	tr.details = append(tr.details, detail)
}

func (tr *testResult) calculateScore(detail reportDetail) {
	if detail.isCorrect {
		tr.correctAnswer++
	}
}

type reportDetail struct {
	question      string
	userAnswer    string
	correctAnswer string
	isCorrect     bool
}

func newReportDetail(prb problem, answer string) reportDetail {
	return reportDetail{
		question:      prb.question,
		userAnswer:    answer,
		correctAnswer: prb.answer,
		isCorrect:     prb.answer == answer,
	}
}

func readUserInput() *userInput {
	filename := flag.String("problems", "problems.csv", "problem list")
	timelimit := flag.Int("time", 5, "limit (in second) to answer one question")

	flag.Parse()

	return &userInput{
		filename:  filename,
		timelimit: timelimit,
	}
}

func getProblems(filename string) []problem {
	lines, err := readCSV(filename)
	if err != nil {
		exit(err)
	}

	return parseLinesProblem(lines)
}

func readCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func parseLinesProblem(lines [][]string) []problem {
	problems := make([]problem, 0, len(lines))
	for _, v := range lines {
		problem := problem{question: v[0], answer: v[1]}
		problems = append(problems, problem)
	}

	return problems
}

func runQuizes(problems []problem, timelimit int) testResult {
	timer := time.NewTimer(time.Duration(timelimit) * time.Second)

	result := newTestResult(len(problems))
	for i, problem := range problems {
		fmt.Printf("#%d %s : ", i+1, problem.question)

		answerChannel := make(chan string)
		go getUserAnswer(answerChannel)

		select {
		case <-timer.C:
			return result
		case userAnswer := <-answerChannel:
			report := newReportDetail(problem, userAnswer)
			result.calculateScore(report)
			result.addDetail(report)
		}
	}

	return result
}

func getUserAnswer(channel chan string) {
	var userAnswer string
	fmt.Scanf("%s\n", &userAnswer)
	channel <- userAnswer
}

func generateReport(testResult testResult) {
	fmt.Printf("you get %d corrent answer from %d question.\n", testResult.correctAnswer, testResult.totalQuestion)
}

func exit(err error) {
	log.Fatalln(err)
}
