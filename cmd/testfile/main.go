package main

import (
	"encoding/csv"
	"io"
	"os"
)

func main() {
	r := csv.NewReader(os.Stdin)
	w := csv.NewWriter(os.Stdout)
	head, _ := r.Read()
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
		w.Write(fs)
	}
}
