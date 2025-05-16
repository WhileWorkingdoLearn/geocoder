package main

/*
func main() {
	input := flag.String("in", "", "input path and filename. must be a .tsv file")
	output := flag.String("out", "", "output path and filename. must be a .tsv file")
	flag.Parse()

	printch := make(chan int64)
	conv := NewConverter(*input, *output)

	go func() {
		err := conv.Start()
		if err != nil {
			log.Fatal(err)
		}
		conv.GetProgress(printch)
	}()

	for p := range printch {
		fmt.Printf("\r%3d%% verarbeitet...", p)
	}
}

*/
