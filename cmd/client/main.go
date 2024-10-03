package main

import (
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
	r := os.Stdin
	s := csv.NewReader(r)
	getBatch := MakeGetBatch(s, 1000)
	doBatch := MakeDoBatch()
	writeBatch := MakeWriteBatch(csv.NewWriter(os.Stdout))
	return MakeProcessDocument(getBatch, doBatch, writeBatch)()
}

type ProcessDocument func() error

func MakeProcessDocument(getBatch GetBatch, doBatch DoBatch, writeBatch WriteBatch) ProcessDocument {
	return func() error {
		var ERR error
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
			outRecords, err := doBatch(records)
			if err != nil {
				ERR = errors.Join(ERR, err)
				continue
			}
			ERR = errors.Join(ERR, writeBatch(outRecords))

		}
		return ERR
	}
}

type Record []string

type GetBatch func() ([]Record, error)

func MakeGetBatch(r *csv.Reader, batchSize int) GetBatch {
	return func() (records []Record, err error) {
		for i := 0; i < batchSize; i++ {
			fields, err := r.Read()
			if err != nil {
				if err == io.EOF && i != 0 {
					return records, nil
				}
				return records, err
			}
			records = append(records, fields)
		}
		return
	}
}

type OutputRecord Record

type DoBatch func([]Record) ([]OutputRecord, error)

func MakeDoBatch() DoBatch {
	return func(records []Record) ([]OutputRecord, error) {
		return Map(records, func(i Record) (o OutputRecord, err error) {
			return OutputRecord(i), nil
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
