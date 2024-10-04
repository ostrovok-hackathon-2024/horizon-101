package main

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"

	"flag"

	"codeberg.org/shinyzero0/ostrovok2024-client/proto"
	. "codeberg.org/shinyzero0/ostrovok2024-client/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var batchSize = flag.Int("n", 1000, "set batch size")
var content = flag.String("content", "", "set input file path")

func main() {
	if err := f(); err != nil {
		log.Fatal(err)
	}
}
func f() error {
	flag.Parse()
	srvu, ok := os.LookupEnv("SERVER_URI")
	if !ok {
		return errors.New("$SERVER_URI undefined")
	}
	var r io.Reader
	if *content == "" {
		r = os.Stdin
	} else {
		rdr, err := os.Open(*content)
		if err != nil {
			return err
		}
		r = rdr
	}
	s := csv.NewReader(r)
	_, err := s.Read()
	if err != nil {
		return err
	}
	fields := []string{"rate_name","class","quality","bathroom","bedding","capacity","club","balcony","view"}
	fieldMap := make(map[string]int, len(fields))
	for k, v := range fields {
		fieldMap[v] = k
	}

	getBatch := MakeGetBatch(s, *batchSize)
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
	balcony := fieldMap["balcony"]
	view := fieldMap["view"]
	return func(records []Record) ([]OutputRecord, error) {
		inputs, _ := Map(records, func(i Record) (o *proto.InputRecord, err error) {
			return &proto.InputRecord{
				Description: string(i[rate_name]),
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
				o[club] = Club(i.IsClub).String()
				o[view] = View(i.View).String()
				o[capacity] = Capacity(i.Capacity).String()
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
