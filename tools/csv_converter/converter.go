package csvconverter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	workerCount = 4
	batchSize   = 100
)

func NewConverter(InputPath string,
	OutputPath string) *Converter {
	return &Converter{
		InputPath:  InputPath,
		OutputPath: OutputPath,
		input:      make(chan *Task, workerCount),
		output:     make(chan *Task, workerCount),
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

	for range workerCount {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()

			for {
				select {
				case task, ok := <-c.input: // you must check for readable state of the channel.
					if !ok {
						return
					}
					c.output <- c.parseStruct(task)
				}
			}
		}()
	}

	go c.scanFile(inputFile)

	go c.waitForWorkersDone()

	writer := bufio.NewWriter(outputFile)

	for r := range c.output {
		for _, line := range r.Data {
			writer.WriteString(line + "\n")
		}

	}
	writer.Flush()

	fmt.Println("\nFinished with converting file, to: ", c.OutputPath)
	return nil
}

func (c *Converter) scanFile(inputFile *os.File) {

	scanner := bufio.NewScanner(inputFile)
	task := &Task{Data: make([]string, 0), lock: &sync.Mutex{}}
	for scanner.Scan() {
		line := scanner.Text()
		atomic.AddInt64(&c.ProcessedBytes, int64(len(line)+1))
		task.Data = append(task.Data, line)
		if len(task.Data) >= batchSize {
			c.input <- task
			task = &Task{Data: make([]string, 0), lock: &sync.Mutex{}}
		}
	}
	if len(task.Data) >= 0 {
		c.input <- task
	}
	close(c.input)
}

func (c *Converter) parseStruct(task *Task) *Task {
	defer task.lock.Unlock()
	task.lock.Lock()
	result := make([]string, 0)
	for _, line := range task.Data {
		columns := strings.Split(line, "\t")
		if len(columns) < 6 {
			continue
		}
		newcolumns := []string{
			columns[0],
			columns[1],
			c.normalizer.ReplaceStreetSuffixWithMarker(columns[2]),
			columns[4],
			columns[6],
			c.normalizer.NormalizeTokenFromStreet(columns[2]),
		}
		newLine := strings.Join(newcolumns, "\t")
		result = append(result, newLine)
	}
	task.Data = result
	return task
}

func (c *Converter) GetProgress(listener chan<- int64) {
	go func() {
		for {
			done := atomic.LoadInt64(&c.ProcessedBytes)
			percent := done * 100 / c.totalSize
			listener <- percent
			if done >= c.totalSize {
				break
			}
		}
		close(listener)
	}()
}

func (c *Converter) waitForWorkersDone() {
	c.wg.Wait()
	close(c.output) // when you close(res) it breaks the below loop.
}
