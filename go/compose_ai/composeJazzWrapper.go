package compose_ai

import "fmt"

func ComposeJazz() error {

	// only jazz for now
	generatorPath := "/compose_ai/python/jazz/genrator.py"
	generateCmd := "python2 generator.py 1"
	generateEpochs := 1

	fmt.Println(generatorPath)
	fmt.Println(generateCmd)
	fmt.Println(generateEpochs)

	return nil
}
