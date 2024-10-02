package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
)

func main() {
	if err := f(); err != nil {
		log.Fatal(err)
	}
}
func f() error {
	getBatch := MakeGetBatch(os.Stdin, 1000)
	doBatch := MakeDoBatch()
	writeBatch := MakeWriteBatch(csv.NewWriter(os.Stdout))
	return MakeProcessDocument(getBatch, doBatch, writeBatch)()
}

type ProcessDocument func() error

func MakeProcessDocument(getBatch GetBatch, doBatch DoBatch, writeBatch WriteBatch) ProcessDocument {
	type Worker struct {
		Errors  chan error
		Results chan []OutputRecord
	}
	return func() error {
		var ERR error
		workers := make([]Worker, 0)
		for {
			records, err := getBatch()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					ERR = errors.Join(ERR, err)
					return ERR
				}
			}
			w := Worker{make(chan error, 1), make(chan []OutputRecord, 1)}
			go func() {
				outRecords, err := doBatch(records)
				w.Errors <- err
				w.Results <- outRecords
			}()
			workers = append(workers, w)
		}
		for _, w := range workers {
			err := <-w.Errors
			outRecords := <-w.Results
			if err != nil {
				ERR = errors.Join(ERR, err)
				continue
			}
			ERR = errors.Join(ERR, writeBatch(outRecords))
		}
		return ERR
	}
}

type Record string

type GetBatch func() ([]Record, error)

func MakeGetBatch(r io.Reader, batchSize int) GetBatch {
	return func() (records []Record, err error) {
		s := bufio.NewScanner(r)
		for i := 0; i < batchSize; i++ {
			if s.Scan() {
				field := s.Text()
				records = append(records, Record(field))
			} else {
				if i == 0 {
					return records, io.EOF
				}
				return records, nil
			}
		}
		return
	}
}

type OutputRecord []string

type DoBatch func([]Record) ([]OutputRecord, error)

func MakeDoBatch() DoBatch {
	return func(records []Record) ([]OutputRecord, error) {
		return Map(records, func(i Record) (o OutputRecord, err error) {
			return OutputRecord([]string{string(i)}), nil
		})
	}
}
func Map[I any, O any](in []I, f func(i I) (o O, err error)) ([]O, error) {
	out := make([]O, len(in))
	var ERR error
	k := 0
	for _, i := range in {
		o, err := f(i)
		if err != nil {
			ERR = errors.Join(ERR, err)
			continue
		}
		out[k] = o
		k++
	}
	return out, ERR
}

type WriteBatch func(records []OutputRecord) error

func MakeWriteBatch(w *csv.Writer) WriteBatch {
	return func(records []OutputRecord) error {
		out, _ := Map(records, func(i OutputRecord) (o []string, err error) {
			return i, nil
		})
		return w.WriteAll(out)
	}
}
