package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var count int = 0
	var c = make(chan []int)
	var resChannel = make(chan string)

	reader := bufio.NewReader(file)

	for {
		if count == 0 {
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			if line == "" {
				break
			}
			line = strings.TrimSuffix(line, "\n")

			scannedCount, err := strconv.Atoi(line)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}

			count = scannedCount
		} else {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			if line == "" {
				break
			}

			sliceOfDigitsAsStrings := strings.Split(line, " ")
			sliceOfDigitsAsInteger := make([]int, len(sliceOfDigitsAsStrings))

			for i, value := range sliceOfDigitsAsStrings {
				sliceOfDigitsAsInteger[i], err = strconv.Atoi(value)

				if err != nil {
					log.Fatal(err)
				}
			}

			go findAbsentNumbers(sliceOfDigitsAsInteger, c, count)
			go writeToFile(<-c, resChannel)
			resChannel <- "Final!"
		}
	}
}

func writeToFile(result []int, c chan string) {
	f, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for i, value := range result {
		if i == len(result) - 1 {
			_, err = f.WriteString(fmt.Sprintf("%d", value))
		} else {
			_, err = f.WriteString(fmt.Sprintf("%d ", value))
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	<-c
}

func findAbsentNumbers(digits []int, c chan []int, count int) {
	sliceOfSortedDigitsAsInteger:= quickSort(digits)
	var sliceOfFindedDigits []int

	if sliceOfSortedDigitsAsInteger[0] != 1 {
		sliceOfFindedDigits = append(sliceOfFindedDigits, 1)
		if sliceOfSortedDigitsAsInteger[0] != 2 {
			sliceOfFindedDigits = append(sliceOfFindedDigits, 2)
		}
	}

	for i, _ := range sliceOfSortedDigitsAsInteger {
		if i < len(sliceOfSortedDigitsAsInteger) - 1 {
			if sliceOfSortedDigitsAsInteger[i + 1] - sliceOfSortedDigitsAsInteger[i] == 3 {
				sliceOfFindedDigits = append(sliceOfFindedDigits, sliceOfSortedDigitsAsInteger[i] + 1)
				sliceOfFindedDigits = append(sliceOfFindedDigits, sliceOfSortedDigitsAsInteger[i] + 2)
			} else if sliceOfSortedDigitsAsInteger[i + 1] - sliceOfSortedDigitsAsInteger[i] == 2 {
				sliceOfFindedDigits = append(sliceOfFindedDigits, sliceOfSortedDigitsAsInteger[i] + 1)
			}
		}
	}

	if len(sliceOfSortedDigitsAsInteger) < count && len(sliceOfFindedDigits) != 2 {
		sliceOfFindedDigits = append(sliceOfFindedDigits, sliceOfSortedDigitsAsInteger[len(sliceOfSortedDigitsAsInteger) - 1] + 1)
	}

	c <- sliceOfFindedDigits
}

func quickSort(digits[]int) []int {
	if len(digits) < 2 {
		return digits
	}

	var left []int
	var right []int
	pivot := digits[rand.Intn(len(digits))]

	for _, value := range digits {
		if pivot > value {
			left = append(left, value)
		} else {
			right = append(right, value)
		}
	}

	return append(quickSort(left), quickSort(right)...)
}
