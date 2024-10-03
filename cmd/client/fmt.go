package main
type BeddingType int32

const (
	BeddingType_UndefinedBedding BeddingType = 0
	BeddingType_BunkBed          BeddingType = 1
	BeddingType_SingleBed        BeddingType = 2
	BeddingType_DoubleBed        BeddingType = 3
	BeddingType_Twin             BeddingType = 4
	BeddingType_Multiple         BeddingType = 5
)

func (b BeddingType) String() string {
	switch b {
	case BeddingType_UndefinedBedding:
		return "undefined"
	case BeddingType_BunkBed:
		return "bunk bed"
	case BeddingType_SingleBed:
		return "single bed"
	case BeddingType_DoubleBed:
		return "double/double or twin"
	case BeddingType_Twin:
		return "twin/twin or double"
	case BeddingType_Multiple:
		return "multiple"
	default:
		return "unknown"
	}
}
type BathroomType int32

const (
	BathroomType_UndefinedBathroom       BathroomType = 0
	BathroomType_SharedBathroom          BathroomType = 1
	BathroomType_PrivateBathroom         BathroomType = 2
	BathroomType_ExternalPrivateBathroom BathroomType = 3
)

func (b BathroomType) String() string {
	switch b {
	case BathroomType_UndefinedBathroom:
		return "undefined"
	case BathroomType_SharedBathroom:
		return "shared bathroom"
	case BathroomType_PrivateBathroom:
		return "private bathroom"
	case BathroomType_ExternalPrivateBathroom:
		return "external private bathroom"
	default:
		return "unknown"
	}
}
type Quality int32

const (
	Quality_UndefinedQuality Quality = 0
	Quality_Economy          Quality = 1
	Quality_Standard         Quality = 2
	Quality_Comfort          Quality = 3
	Quality_Business         Quality = 4
	Quality_Superior         Quality = 5
	Quality_Deluxe           Quality = 6
	Quality_Premier          Quality = 7
	Quality_Executive        Quality = 8
	Quality_Presidential     Quality = 9
	Quality_Premium          Quality = 10
	Quality_Classic          Quality = 11
	Quality_Ambassador       Quality = 12
	Quality_Grand            Quality = 13
	Quality_Luxury           Quality = 14
	Quality_Platinum         Quality = 15
	Quality_Prestige         Quality = 16
	Quality_Privilege        Quality = 17
	Quality_Royal            Quality = 18
)

func (q Quality) String() string {
	switch q {
	case Quality_UndefinedQuality:
		return "undefined"
	case Quality_Economy:
		return "economy"
	case Quality_Standard:
		return "standard"
	case Quality_Comfort:
		return "comfort"
	case Quality_Business:
		return "business"
	case Quality_Superior:
		return "superior"
	case Quality_Deluxe:
		return "deluxe"
	case Quality_Premier:
		return "premier"
	case Quality_Executive:
		return "executive"
	case Quality_Presidential:
		return "presidential"
	case Quality_Premium:
		return "premium"
	case Quality_Classic:
		return "classic"
	case Quality_Ambassador:
		return "ambassador"
	case Quality_Grand:
		return "grand"
	case Quality_Luxury:
		return "luxury"
	case Quality_Platinum:
		return "platinum"
	case Quality_Prestige:
		return "prestige"
	case Quality_Privilege:
		return "privilege"
	case Quality_Royal:
		return "royal"
	default:
		return "unknown"
	}
}


type RoomClass int32

const (
	RoomClass_RunOfHouse  RoomClass = 0
	RoomClass_Dorm        RoomClass = 1
	RoomClass_Capsule     RoomClass = 2
	RoomClass_Room        RoomClass = 3
	RoomClass_JuniorSuite RoomClass = 4
	RoomClass_Suite       RoomClass = 5
	RoomClass_Apartment   RoomClass = 6
	RoomClass_Studio      RoomClass = 7
	RoomClass_Villa       RoomClass = 8
	RoomClass_Cottage     RoomClass = 9
	RoomClass_Bungalow    RoomClass = 10
	RoomClass_Chalet      RoomClass = 11
	RoomClass_Camping     RoomClass = 12
)

func (r RoomClass) String() string {
	switch r {
	case RoomClass_RunOfHouse:
		return "run-of-house"
	case RoomClass_Dorm:
		return "dorm"
	case RoomClass_Capsule:
		return "capsule"
	case RoomClass_Room:
		return "room"
	case RoomClass_JuniorSuite:
		return "junior-suite"
	case RoomClass_Suite:
		return "suite"
	case RoomClass_Apartment:
		return "apartment"
	case RoomClass_Studio:
		return "studio"
	case RoomClass_Villa:
		return "villa"
	case RoomClass_Cottage:
		return "cottage"
	case RoomClass_Bungalow:
		return "bungalow"
	case RoomClass_Chalet:
		return "chalet"
	case RoomClass_Camping:
		return "camping"
	default:
		return "unknown"
	}
}
type Capacity int32

const (
	Capacity_UndefinedCapacity Capacity = 0
	Capacity_Single            Capacity = 1
	Capacity_DoubleCapacity    Capacity = 2
	Capacity_Triple            Capacity = 3
	Capacity_Quadruple         Capacity = 4
	Capacity_Quintuple         Capacity = 5
	Capacity_Sextuple          Capacity = 6
)

func (c Capacity) String() string {
	switch c {
	case Capacity_UndefinedCapacity:
		return "undefined"
	case Capacity_Single:
		return "single"
	case Capacity_DoubleCapacity:
		return "double"
	case Capacity_Triple:
		return "triple"
	case Capacity_Quadruple:
		return "quadruple"
	case Capacity_Quintuple:
		return "quintuple"
	case Capacity_Sextuple:
		return "sextuple"
	default:
		return "unknown"
	}
}


type View int32

const (
	View_UndefinedView    View = 0
	View_BayView          View = 1
	View_BosphorusView    View = 2
	View_BurjKhalifaView  View = 3
	View_CanalView        View = 4
	View_CityView         View = 5
	View_CourtyardView    View = 6
	View_DubaiMarinaView  View = 7
	View_GardenView       View = 8
	View_GolfView         View = 9
	View_HarbourView      View = 10
	View_InlandView       View = 11
	View_KremlinView      View = 12
	View_LakeView         View = 13
	View_LandView         View = 14
	View_MountainView     View = 15
	View_OceanView        View = 16
	View_PanoramicView    View = 17
	View_ParkView         View = 18
	View_PartialOceanView View = 19
	View_PartialSeaView   View = 20
	View_PartialView      View = 21
	View_PoolView         View = 22
	View_RiverView        View = 23
	View_SeaView          View = 24
	View_SheikhZayedView  View = 25
	View_StreetView       View = 26
	View_SunriseView      View = 27
	View_SunsetView       View = 28
	View_WaterView        View = 29
	View_WithView         View = 30
	View_Beachfront       View = 31
	View_OceanFront       View = 32
	View_SeaFront         View = 33
)

func (v View) String() string {
	switch v {
	case View_UndefinedView:
		return "undefined"
	case View_BayView:
		return "bay view"
	case View_BosphorusView:
		return "bosphorus view"
	case View_BurjKhalifaView:
		return "burj-khalifa view"
	case View_CanalView:
		return "canal view"
	case View_CityView:
		return "city view"
	case View_CourtyardView:
		return "courtyard view"
	case View_DubaiMarinaView:
		return "dubai-marina view"
	case View_GardenView:
		return "garden view"
	case View_GolfView:
		return "golf view"
	case View_HarbourView:
		return "harbour view"
	case View_InlandView:
		return "inland view"
	case View_KremlinView:
		return "kremlin view"
	case View_LakeView:
		return "lake view"
	case View_LandView:
		return "land view"
	case View_MountainView:
		return "mountain view"
	case View_OceanView:
		return "ocean view"
	case View_PanoramicView:
		return "panoramic view"
	case View_ParkView:
		return "park view"
	case View_PartialOceanView:
		return "partial-ocean view"
	case View_PartialSeaView:
		return "partial-sea view"
	case View_PartialView:
		return "partial view"
	case View_PoolView:
		return "pool view"
	case View_RiverView:
		return "river view"
	case View_SeaView:
		return "sea view"
	case View_SheikhZayedView:
		return "sheikh-zayed view"
	case View_StreetView:
		return "street view"
	case View_SunriseView:
		return "sunrise view"
	case View_SunsetView:
		return "sunset view"
	case View_WaterView:
		return "water view"
	case View_WithView:
		return "with view"
	case View_Beachfront:
		return "beachfront"
	case View_OceanFront:
		return "ocean front"
	case View_SeaFront:
		return "sea front"
	default:
		return "unknown"
	}
}

type Floor int32

const (
	Floor_UndefinedFloor Floor = 0
	Floor_PenthouseFloor Floor = 1
	Floor_DuplexFloor    Floor = 2
	Floor_BasementFloor  Floor = 3
	Floor_AtticFloor     Floor = 4
)

// String метод для типа Floor
func (f Floor) String() string {
	switch f {
	case Floor_UndefinedFloor:
		return "undefined"
	case Floor_PenthouseFloor:
		return "penthouse floor"
	case Floor_DuplexFloor:
		return "duplex floor"
	case Floor_BasementFloor:
		return "basement floor"
	case Floor_AtticFloor:
		return "attic floor"
	default:
		return "unknown floor"
	}
}

type Club bool

func (c Club) String() string {
	if c {return "club"} else {return "not club"}
}
type Balcony bool

func (b Balcony) String() string {
	if b {return "balcony"} else {return "no balcony"}
}


type Bedrooms int32

const (
	Bedrooms_Undefined Bedrooms = 0
	Bedrooms_One       Bedrooms = 1
	Bedrooms_Two       Bedrooms = 2
	Bedrooms_Three     Bedrooms = 3
	Bedrooms_Four      Bedrooms = 4
	Bedrooms_Five      Bedrooms = 5
	Bedrooms_Six       Bedrooms = 6
)

// String метод для типа Bedrooms
func (b Bedrooms) String() string {
	switch b {
	case Bedrooms_Undefined:
		return "undefined"
	case Bedrooms_One:
		return "1 bedroom"
	case Bedrooms_Two:
		return "2 bedrooms"
	case Bedrooms_Three:
		return "3 bedrooms"
	case Bedrooms_Four:
		return "4 bedrooms"
	case Bedrooms_Five:
		return "5 bedrooms"
	case Bedrooms_Six:
		return "6 bedrooms"
	default:
		return "unknown number of bedrooms"
	}
}
