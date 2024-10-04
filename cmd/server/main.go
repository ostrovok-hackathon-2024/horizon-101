package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"

	"codeberg.org/shinyzero0/ostrovok2024-client/proto"
	"codeberg.org/shinyzero0/ostrovok2024-client/utils"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"google.golang.org/grpc"
)

func main() {
	if err := f(); err != nil {
		log.Fatal(err)
	}
}

type Matchers []Matcher

func (ms Matchers) Apply(s string, rec *proto.OutputRecord) error {
	mp := bleve.NewIndexMapping()
	mp.DefaultAnalyzer = "en"
	index, err := bleve.NewMemOnly(mp)
	if err != nil {
		return err
	}
	defer index.Close()
	index.Index("this", s)
	for _, m := range ms {
		result, _ := index.Search(bleve.NewSearchRequest(m.Query))
		if result.MaxScore > 0 {
			m.Action(rec)
		}
	}
	return nil
}

type BestMatcher struct {
	Matchers []Matcher
}
type Action func(*proto.OutputRecord)

func (bm BestMatcher) Apply(s string, rec *proto.OutputRecord) error {
	mp := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(mp)
	if err != nil {
		return err
	}
	defer index.Close()
	index.Index("this", s)
	var best struct {
		score float64
		Action
	}
	for _, m := range bm.Matchers {
		result, _ := index.Search(bleve.NewSearchRequest(m.Query))
		if result.MaxScore > best.score {
			best.Action = m.Action
		}
	}
	if best.Action != nil {
		best.Action(rec)
	}
	return nil
}

type Matcher struct {
	query.Query
	Action func(*proto.OutputRecord)
}

type server struct {
	proto.UnimplementedProcessorServer
	Matchers []BestMatcher
}

func (s *server) ProcessBatch(_ context.Context, in *proto.BatchedInput) (*proto.BatchedOutput, error) {
	outs, err := utils.Map(
		in.Records,
		func(i *proto.InputRecord) (*proto.OutputRecord, error) {
			o := &proto.OutputRecord{
				Bathroom:       proto.BathroomType_PrivateBathroom,
				Class:          proto.RoomClass_Room,
				BedroomsAmount: 1,
				BeddingType:    proto.BeddingType_UndefinedBedding,
				HasBalcony:     false,
				IsClub:         false,
				View:           proto.View_UndefinedView,
				Floor:          proto.Floor_UndefinedFloor,
				Capacity:       proto.Capacity_DoubleCapacity,
				Quality:        proto.Quality_Standard,
			}
			_, err := utils.Map(s.Matchers, func(m BestMatcher) (*struct{}, error) {
				return nil, m.Apply(i.Description, o)
			})
			if o.Capacity == 0 {
				o.Capacity = proto.Capacity(o.BedroomsAmount)
			}
			return o, err
		},
	)
	return &proto.BatchedOutput{
		Records: outs,
	}, err
}

func f() error {
	srvu, ok := os.LookupEnv("SERVER_URI")
	if !ok {
		return errors.New("fuck")
	}
	lis, err := net.Listen("tcp", srvu)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	var srv server
	viewMatcher := func(q string, v proto.View) Matcher {
		return Matcher{
			Action: func(o *proto.OutputRecord) { o.View = v },
			Query:  bleve.NewQueryStringQuery(q)}
	}
	srv.Matchers = []BestMatcher{
		// capacity
		BestMatcher{Matchers: []Matcher{
			Matcher{
				Action: func(o *proto.OutputRecord) { o.Capacity = 2 },
				Query:  bleve.NewQueryStringQuery(`+/single|one/ -"/2|double/ bed" room`)},
			Matcher{
				Action: func(o *proto.OutputRecord) { o.Capacity = 3 },
				Query:  bleve.NewQueryStringQuery(`+/3|triple|three/ -"/3|triple/ bed" room`)},
			Matcher{
				Action: func(o *proto.OutputRecord) { o.Capacity = 4 },
				Query:  bleve.NewQueryStringQuery(`+/4|quadruple/ -"/4|quadruple/ bed" room`)},
			Matcher{
				Action: func(o *proto.OutputRecord) { o.Capacity = 5 },
				Query:  bleve.NewQueryStringQuery(`+/5|quintuple/ -"/5|quintuple/ bed" room`)},
			Matcher{
				Action: func(o *proto.OutputRecord) { o.Capacity = 6 },
				Query:  bleve.NewQueryStringQuery(`+/6|sextuple/ -"/6|sextuple/ bed" room`)},
		}},
		// class
		BestMatcher{Matchers: []Matcher{
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+villa`),
				Action: func(o *proto.OutputRecord) { o.Class = proto.RoomClass_Villa }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+"run of house"`),
				Action: func(o *proto.OutputRecord) { o.Class = proto.RoomClass_RunOfHouse }},
			Matcher{
				Query: bleve.NewQueryStringQuery(`dorm dormitory dormbed`),
				Action: func(o *proto.OutputRecord) {
					o.Class = proto.RoomClass_Dorm
					o.Bathroom = proto.BathroomType_SharedBathroom
				}},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+capsule`),
				Action: func(o *proto.OutputRecord) { o.Class = proto.RoomClass_Capsule }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+"junior suite"`),
				Action: func(o *proto.OutputRecord) { o.Class = proto.RoomClass_JuniorSuite }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+suite`),
				Action: func(o *proto.OutputRecord) { o.Class = proto.RoomClass_Suite }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+apartment`),
				Action: func(o *proto.OutputRecord) { o.Class = proto.RoomClass_Apartment }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+studio`),
				Action: func(o *proto.OutputRecord) { o.Class = proto.RoomClass_Studio }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+cottage`),
				Action: func(o *proto.OutputRecord) { o.Class = proto.RoomClass_Cottage }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+bungalow`),
				Action: func(o *proto.OutputRecord) { o.Class = proto.RoomClass_Bungalow }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+chalet`),
				Action: func(o *proto.OutputRecord) { o.Class = proto.RoomClass_Chalet }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+camping`),
				Action: func(o *proto.OutputRecord) { o.Class = proto.RoomClass_Camping }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+tent`),
				Action: func(o *proto.OutputRecord) { o.Class = proto.RoomClass_Tent }},
		}},
		// quality
		BestMatcher{Matchers: []Matcher{
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+economy`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Economy }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+standard`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Standard }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+comfort`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Comfort }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+business`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Business }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+superior`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Superior }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+deluxe`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Deluxe }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+executive`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Executive }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+premier`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Premier }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+presidential`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Presidential }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+premium`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Premium }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+classic`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Classic }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+ambassador`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Ambassador }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+grand`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Grand }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+luxury`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Luxury }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+platinum`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Platinum }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+prestige`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Prestige }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+privilege`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Privilege }},
			Matcher{
				Query:  bleve.NewQueryStringQuery(`+royal`),
				Action: func(o *proto.OutputRecord) { o.Quality = proto.Quality_Royal }},
		}},
		BestMatcher{Matchers: []Matcher{
			Matcher{
				Query:  bleve.NewQueryStringQuery(`"-not club"`),
				Action: func(o *proto.OutputRecord) { o.IsClub = true }}}},
		// capacity
		BestMatcher{Matchers: []Matcher{
			Matcher{
				Action: func(o *proto.OutputRecord) { o.BeddingType = proto.BeddingType_Twin },
				Query:  bleve.NewQueryStringQuery(`+"/2|double|twin/ bed"`)},
			Matcher{
				Action: func(o *proto.OutputRecord) { o.BeddingType = proto.BeddingType_SingleBed },
				Query:  bleve.NewQueryStringQuery(`+"/1|single/ bed"`)},
			Matcher{
				Action: func(o *proto.OutputRecord) { o.BeddingType = proto.BeddingType_BunkBed },
				Query:  bleve.NewQueryStringQuery(`+bunk`)},
			Matcher{
				Action: func(o *proto.OutputRecord) { o.BeddingType = proto.BeddingType_Multiple },
				Query:  bleve.NewQueryStringQuery(`+"/3|4|triple|multiple/ bed" -room`)},
		},
		},
		BestMatcher{Matchers: []Matcher{
			Matcher{
				Action: func(o *proto.OutputRecord) { o.Bathroom = proto.BathroomType_SharedBathroom },
				Query:  bleve.NewQueryStringQuery(`+shared bathroom`)},
		}},
		BestMatcher{Matchers: []Matcher{
			viewMatcher(`bay`, proto.View_BayView),
			viewMatcher(`bosphorus`, proto.View_BosphorusView),
			viewMatcher(`burj-khalifa`, proto.View_BurjKhalifaView),
			viewMatcher(`canal`, proto.View_CanalView),
			viewMatcher(`city view`, proto.View_CityView),
			viewMatcher(`courtyard`, proto.View_CourtyardView),
			viewMatcher(`dubai-marina`, proto.View_DubaiMarinaView),
			viewMatcher(`garden`, proto.View_GardenView),
			viewMatcher(`golf`, proto.View_GolfView),
			viewMatcher(`harbour`, proto.View_HarbourView),
			viewMatcher(`inland`, proto.View_InlandView),
			viewMatcher(`kremlin`, proto.View_KremlinView),
			viewMatcher(`lake`, proto.View_LakeView),
			viewMatcher(`land`, proto.View_LandView),
			viewMatcher(`mountain`, proto.View_MountainView),
			viewMatcher(`ocean`, proto.View_OceanView),
			viewMatcher(`panoramic`, proto.View_PanoramicView),
			viewMatcher(`park`, proto.View_ParkView),
			viewMatcher(`partial ocean`, proto.View_PartialOceanView),
			viewMatcher(`partial sea`, proto.View_PartialSeaView),
			viewMatcher(`pool`, proto.View_PoolView),
			viewMatcher(`river`, proto.View_RiverView),
			viewMatcher(`sea`, proto.View_SeaView),
			viewMatcher(`sheikh-zayed`, proto.View_SheikhZayedView),
			viewMatcher(`street view`, proto.View_StreetView),
			viewMatcher(`sunrise view`, proto.View_SunriseView),
			viewMatcher(`sunset view`, proto.View_SunsetView),
			viewMatcher(`/water|fountain|creek|waterfront|/ view`, proto.View_WaterView),
			viewMatcher(`view`, proto.View_WithView),
			viewMatcher(`beach front`, proto.View_Beachfront),
			viewMatcher(`sea front`, proto.View_SeaFront),
		}},
		// floors
		// BestMatcher{Matchers: []Matcher{
		// 	Matcher{
		// 		Action: func(o *proto.OutputRecord) { o.Floor = proto.Floor_AtticFloor },
		// 		Query:  bleve.NewQueryStringQuery(`+attic`)},
		// 	Matcher{
		// 		Action: func(o *proto.OutputRecord) { o.Floor = proto.Floor_PenthouseFloor },
		// 		Query:  bleve.NewQueryStringQuery(`+penthouse`)},
		// 	Matcher{
		// 		Action: func(o *proto.OutputRecord) { o.Floor = proto.Floor_DuplexFloor },
		// 		Query:  bleve.NewQueryStringQuery(`+duplex`)},
		// 	Matcher{
		// 		Action: func(o *proto.OutputRecord) { o.Floor = proto.Floor_BasementFloor },
		// 		Query:  bleve.NewQueryStringQuery(`+basement`)},
		// }},
	}
	proto.RegisterProcessorServer(s, &srv)
	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis)
}
