package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

func parseStruct(job []string) string {
	return strings.Join(job, "*")
}

const (
	workerCount = 4
	batchSize   = 100
)

type Converter struct {
	InputPath      string
	OutputPath     string
	ProcessedBytes int64
	totalSize      int64
	input          chan []string
	output         chan string
	printOut       chan int
	wg             sync.WaitGroup
	normalizer     Normalizer
}

type Task struct {
	Data []string
}

func NewConverter(InputPath string,
	OutputPath string) *Converter {
	return &Converter{
		InputPath:  InputPath,
		OutputPath: OutputPath,
		input:      make(chan []string, workerCount),
		output:     make(chan string, workerCount),
		printOut:   make(chan int),
		normalizer: NewNormalizer(),
	}
}

func (c *Converter) Start() error {
	inputFile, err := os.Open(c.InputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(c.OutputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	stat, errSize := inputFile.Stat()
	if errSize != nil {
		return errSize
	}

	c.totalSize = stat.Size()

	worker := func(jobs <-chan []string, results chan<- string) {
		for {
			select {
			case job, ok := <-jobs: // you must check for readable state of the channel.
				if !ok {
					return
				}
				results <- parseStruct(job)
			}
		}
	}

	for range workerCount {
		c.wg.Add(1)
		go func() {
			// this line will exec when chan `res` processed output at line 107 (func worker: line 71)
			defer c.wg.Done()
			worker(c.input, c.output)
		}()
	}

	go func() {
		scanner := bufio.NewScanner(inputFile)
		for scanner.Scan() {
			line := scanner.Text()
			atomic.AddInt64(&c.ProcessedBytes, int64(len(line)+1))

			c.input <- strings.Split(line, "\t")
		}
		close(c.input) // close jobs to signal workers that no more job are incoming.
	}()

	go func() {
		for {
			done := atomic.LoadInt64(&c.ProcessedBytes)
			percent := done * 100 / c.totalSize
			fmt.Printf("\r%3d%% verarbeitet...", percent)
			if done >= c.totalSize {
				break
			}
		}
	}()

	go func() {
		c.wg.Wait()
		close(c.output) // when you close(res) it breaks the below loop.
	}()

	writer := bufio.NewWriter(outputFile)

	for r := range c.output {
		writer.WriteString(r + "\n")
	}
	writer.Flush()

	fmt.Println("\nFinished with converting file, to: ", c.OutputPath)
	return nil
}

func main() {
	input := flag.String("in", "", "input path and filename. must be a .tsv file")
	output := flag.String("out", "", "output path and filename. must be a .tsv file")
	flag.Parse()

	conv := NewConverter(*input, *output)
	err := conv.Start()
	if err != nil {
		log.Fatal(err)
	}

}
