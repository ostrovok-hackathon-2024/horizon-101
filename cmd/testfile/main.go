package main

import (
	"encoding/csv"
	"fmt"
	// "fmt"
	"io"
	"os"
)

func main() {
	r := csv.NewReader(os.Stdin)
	w := csv.NewWriter(os.Stdout)
	head, _ := r.Read()
	l := len(head)
	// fmt.Println(l) 
	w.Write(head)
	for {
		fs, err := r.Read()
		if err == io.EOF {
			break
		}
		for k, _ := range fs {
			if k != 0 {
				fs[k] = ""
			}
		}
		fs[l-1] = ""
		w.Write(fs)
		w.Flush()
		if err := w.Error(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		// fmt.Println(len(fs)) 
	}
}
