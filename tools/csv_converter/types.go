package csvconverter

import "sync"

type Normalizer struct {
	filter         string
	pattern        string
	streetSuffixes []string
	suffixMap      map[string]string
}

type Converter struct {
	InputPath      string
	OutputPath     string
	ProcessedBytes int64
	totalSize      int64
	input          chan *Task
	output         chan *Task
	printOut       chan int
	wg             sync.WaitGroup
	normalizer     Normalizer
}

type Task struct {
	Data []string
	lock *sync.Mutex
}
