package main

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"

	"codeberg.org/shinyzero0/ostrovok2024-client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := f(); err != nil {
		log.Fatal(err)
	}
}
func f() error {
	srvu, ok := os.LookupEnv("SERVER_URI")
	if !ok {
		return errors.New("fuck")
	}
	r := os.Stdin
	s := csv.NewReader(r)
	fields, err := s.Read()
	if err != nil {
		return err
	}
	fieldMap := make(map[string]int, len(fields))
	for k, v := range fields {
		fieldMap[v] = k
	}

	getBatch := MakeGetBatch(s, 1000)
	conn, err := grpc.NewClient(srvu, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	cli := proto.NewProcessorClient(conn)
	ctx := context.Background()
	doBatch := MakeDoBatch(cli, ctx, fieldMap)
	writeBatch := MakeWriteBatch(csv.NewWriter(os.Stdout), fields)
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

func MakeDoBatch(c proto.ProcessorClient, ctx context.Context, fieldMap map[string]int) DoBatch {

	rate_name := fieldMap["rate_name"]
	class := fieldMap["class"]
	quality := fieldMap["quality"]
	bathroom := fieldMap["bathroom"]
	bedding := fieldMap["bedding"]
	capacity := fieldMap["capacity"]
	club := fieldMap["club"]
	bedrooms := fieldMap["bedrooms"]
	balcony := fieldMap["balcony"]
	view := fieldMap["view"]
	floor := fieldMap["floor"]
	return func(records []Record) ([]OutputRecord, error) {
		inputs, _ := Map(records, func(i Record) (o *proto.InputRecord, err error) {
			return &proto.InputRecord{
				Description: string(i[0]),
			}, nil
		})

		in := proto.BatchedInput{
			Records: inputs,
		}
		results, err := c.ProcessBatch(ctx, &in)
		if err != nil {
			return nil, err
		}
		return MapIdx(
			results.Records,
			func(i *proto.OutputRecord, idx int) (o OutputRecord, err error) {
				o = make(OutputRecord, len(fieldMap))
				o[bedding] = BeddingType(i.BeddingType).String()
				o[floor] = Floor(i.Floor).String()
				o[club] = Club(i.IsClub).String()
				o[view] = View(i.View).String()
				o[capacity] = Capacity(i.Capacity).String()
				o[bedrooms] = Bedrooms(i.BedroomsAmount).String()
				o[balcony] = Balcony(i.HasBalcony).String()
				o[class] = RoomClass(i.Class).String()
				o[bathroom] = BathroomType(i.Bathroom).String()
				o[quality] = Quality(i.Quality).String()
				o[rate_name] = inputs[idx].Description
				return o, nil
			},
		)
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
	return out[:k], ERR
}
func MapIdx[I any, O any](in []I, f func(i I, idx int) (o O, err error)) ([]O, error) {
	out := make([]O, len(in))
	var ERR error
	k := 0
	for idx, i := range in {
		o, err := f(i, idx)
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

func MakeWriteBatch(w *csv.Writer, header []string) WriteBatch {
	w.Write(header)
	return func(records []OutputRecord) error {
		out, _ := Map(records, func(i OutputRecord) (o []string, err error) {
			return i, nil
		})
		return w.WriteAll(out)
	}
}
