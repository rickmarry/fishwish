package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type seedSpot struct {
	Name        string
	Lat         float64
	Lon         float64
	Type        string
	Description string
	AccessNotes string
	Parking     string
	Difficulty  string
	DepthInfo   string
	BottomType  string
	Species     []string
}

var spots = []seedSpot{
	{
		Name:        "Lake Tahoe - Emerald Bay",
		Lat:         38.9531,
		Lon:         -120.1028,
		Type:        "lake",
		Description: "Stunning alpine lake with crystal clear waters. Great for trout fishing with deep drop-offs near the shore.",
		AccessNotes: "Boat ramp available at Emerald Bay State Park. Shore fishing accessible from multiple points.",
		Parking:     "Large lot at state park, $10 day use fee",
		Difficulty:  "easy",
		DepthInfo:   "Shore drops to 20ft quickly, middle reaches 1600ft",
		BottomType:  "Rocky with some sandy areas",
		Species:     []string{"rainbow_trout", "lake_trout", "brown_trout", "kokanee_salmon"},
	},
	{
		Name:        "Everglades - Florida Bay",
		Lat:         25.0767,
		Lon:         -80.7489,
		Type:        "estuary",
		Description: "World-class flats fishing for tarpon, bonefish, and permit. Shallow water sight fishing paradise.",
		AccessNotes: "Best accessed by boat or guided tour. Some shore access from Card Sound Bridge.",
		Parking:     "Limited roadside parking at bridge access points",
		Difficulty:  "hard",
		DepthInfo:   "Very shallow, 1-4ft flats with deeper channels",
		BottomType:  "Mud, sand, and turtle grass flats",
		Species:     []string{"tarpon", "bonefish", "permit", "snook", "redfish"},
	},
	{
		Name:        "Mississippi River - Pool 7",
		Lat:         43.8508,
		Lon:         -91.2393,
		Type:        "river",
		Description: "Excellent walleye and sauger fishery with diverse habitat from backwaters to main channel.",
		AccessNotes: "Multiple public boat launches. Shore fishing at Lock & Dam 7.",
		Parking:     "Free parking at dam and several launch points",
		Difficulty:  "moderate",
		DepthInfo:   "Backwaters 3-8ft, main channel 10-20ft",
		BottomType:  "Sand and gravel with rock structures",
		Species:     []string{"walleye", "sauger", "smallmouth_bass", "catfish", "crappie"},
	},
	{
		Name:        "Kachemak Bay - Alaska",
		Lat:         59.6203,
		Lon:         -151.2758,
		Type:        "ocean",
		Description: "Premier saltwater fishing for halibut, salmon, and lingcod. Pristine wilderness setting.",
		AccessNotes: "Charter boats available from Homer Spit. Some shore fishing accessible.",
		Parking:     "Homer Spit public lots",
		Difficulty:  "moderate",
		DepthInfo:   "Intertidal to 300ft+ in channels",
		BottomType:  "Rocky reefs, kelp beds, and muddy bottoms",
		Species:     []string{"halibut", "silver_salmon", "king_salmon", "lingcod", "rockfish"},
	},
	{
		Name:        "White River - Arkansas",
		Lat:         36.3137,
		Lon:         -92.3899,
		Type:        "river",
		Description: "Tailwater trout fishery below Bull Shoals Dam. Year-round cold water creates perfect trout habitat.",
		AccessNotes: "Public access at multiple points. Float trips highly recommended.",
		Parking:     "Free access points every few miles",
		Difficulty:  "easy",
		DepthInfo:   "Wadeable 3-6ft, pools to 15ft",
		BottomType:  "Gravel with some boulders",
		Species:     []string{"rainbow_trout", "brown_trout", "cutthroat_trout"},
	},
	{
		Name:        "Lake Okeechobee - Florida",
		Lat:         26.9342,
		Lon:         -80.7928,
		Type:        "lake",
		Description: "America's largest lake south of the Great Lakes. Legendary largemouth bass fishery.",
		AccessNotes: "Multiple marinas and public ramps around the lake. Shore fishing on the dike.",
		Parking:     "Free at most access points",
		Difficulty:  "easy",
		DepthInfo:   "Very shallow, average 9ft",
		BottomType:  "Muck and sand with abundant vegetation",
		Species:     []string{"largemouth_bass", "crappie", "bluegill", "catfish"},
	},
	{
		Name:        "Green River - Kentucky",
		Lat:         37.3842,
		Lon:         -86.5836,
		Type:        "river",
		Description: "Trophy musky and smallmouth bass river with scenic limestone bluffs.",
		AccessNotes: "Boat launch at Brownsville. Some wade fishing available.",
		Parking:     "Small lot at boat ramp",
		Difficulty:  "moderate",
		DepthInfo:   "Runs 4-12ft, deeper holes to 20ft",
		BottomType:  "Rock and gravel with undercut banks",
		Species:     []string{"muskellunge", "smallmouth_bass", "spotted_bass", "catfish"},
	},
	{
		Name:        "Pyramid Lake - Nevada",
		Lat:         40.0386,
		Lon:         -119.5678,
		Type:        "lake",
		Description: "Iconic desert lake known for massive Lahontan cutthroat trout. Unique fishery on tribal land.",
		AccessNotes: "Tribal permit required. Multiple access points around the lake.",
		Parking:     "Designated areas at each access point",
		Difficulty:  "moderate",
		DepthInfo:   "Shore drops quickly to 100ft+",
		BottomType:  "Alkaline mud and tufa formations",
		Species:     []string{"lahontan_cutthroat_trout", "cui-ui"},
	},
	{
		Name:        "Kenai River - Alaska",
		Lat:         60.4178,
		Lon:         -150.8372,
		Type:        "river",
		Description: "World-famous for trophy king salmon and giant rainbow trout. Blue glacial waters.",
		AccessNotes: "Excellent road access with multiple parks and pullouts.",
		Parking:     "Various state recreation sites, $5 day fee",
		Difficulty:  "easy",
		DepthInfo:   "Wadeable runs 2-6ft, deeper holes 8-15ft",
		BottomType:  "Gravel with boulder structures",
		Species:     []string{"king_salmon", "silver_salmon", "sockeye_salmon", "rainbow_trout", "dolly_varden"},
	},
	{
		Name:        "Lake Powell - Utah/Arizona",
		Lat:         37.0377,
		Lon:         -111.3377,
		Type:        "reservoir",
		Description: "Massive reservoir with endless canyons and coves. Excellent striped bass and smallmouth.",
		AccessNotes: "Boat strongly recommended. Some shore access near marinas.",
		Parking:     "Marina parking, $15-20/day",
		Difficulty:  "moderate",
		DepthInfo:   "Varies wildly with water level, 10-200ft",
		BottomType:  "Rock and submerged canyon walls",
		Species:     []string{"striped_bass", "smallmouth_bass", "largemouth_bass", "walleye"},
	},
}

var speciesData = []struct {
	Name        string
	CommonName  string
	Description string
	AvgWeight   float64
	BestSeasons []string
	Bait        []string
}{
	{"largemouth_bass", "Largemouth Bass", "Most popular gamefish in North America", 4.5, []string{"spring", "fall"}, []string{"plastic worms", "crankbaits", "topwater"}},
	{"smallmouth_bass", "Smallmouth Bass", "Hard-fighting fighter in clear water", 3.0, []string{"spring", "summer", "fall"}, []string{"tubes", "jigs", "finesse worms"}},
	{"rainbow_trout", "Rainbow Trout", "Widely stocked and aggressively bites", 2.5, []string{"spring", "fall"}, []string{"powerbait", "spinners", "flies"}},
	{"brown_trout", "Brown Trout", "Cunning and challenging to catch", 4.0, []string{"fall", "winter"}, []string{"streamers", "minnows", "nightcrawlers"}},
	{"lake_trout", "Lake Trout", "Deep-water predator of northern lakes", 12.0, []string{"winter", "spring"}, []string{"jigs", "spoons", "trolling lures"}},
	{"walleye", "Walleye", "Excellent table fare, active at dawn/dusk", 3.5, []string{"spring", "fall"}, []string{"jigs", "crankbaits", "nightcrawlers"}},
	{"catfish", "Channel Catfish", "Bottom feeder with keen sense of smell", 5.0, []string{"summer"}, []string{"cut bait", "chicken liver", "stink bait"}},
	{"crappie", "Crappie", "Popular panfish, great for frying", 0.75, []string{"spring"}, []string{"small jigs", "minnows", "spiders"}},
	{"tarpon", "Tarpon", "Silver king, spectacular acrobatic fighter", 80.0, []string{"spring", "summer"}, []string{"live crabs", "mullet", "artificial flies"}},
	{"redfish", "Red Drum", "Powerful inshore saltwater species", 12.0, []string{"fall", "spring"}, []string{"gold spoons", "shrimp", "crabs"}},
	{"snook", "Snook", "Aggressive estuarine predator", 15.0, []string{"summer", "fall"}, []string{"live shrimp", "pilchards", "topwater plugs"}},
	{"king_salmon", "King Salmon", "Largest Pacific salmon, trophy fish", 25.0, []string{"summer"}, []string{"salmon roe", "spinners", "flies"}},
	{"silver_salmon", "Coho Salmon", "Acrobatic and aggressive fighter", 8.0, []string{"summer", "fall"}, []string{"spinners", "salmon eggs", "flies"}},
	{"halibut", "Pacific Halibut", "Massive flatfish, ultimate bottom fish", 40.0, []string{"summer"}, []string{"herring", "squid", "circle hooks"}},
	{"muskellunge", "Musky", "The fish of 10,000 casts", 20.0, []string{"fall"}, []string{"large jerkbaits", "bucktails", "topwater"}},
}

func main() {
	_ = godotenv.Load("services/spot-service/.env.local")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("DB_USER", "fishwish"),
		getEnv("DB_PASSWORD", "fishwish"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "fishwish"),
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("database not reachable: %v", err)
	}

	ctx := context.Background()

	log.Println("seeding species...")
	for _, sp := range speciesData {
		_, err := pool.Exec(ctx, `
			INSERT INTO species (name, common_name, description, avg_weight_lbs, best_seasons, preferred_bait)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (name) DO NOTHING`,
			sp.Name, sp.CommonName, sp.Description, sp.AvgWeight, sp.BestSeasons, sp.Bait,
		)
		if err != nil {
			log.Printf("error inserting species %s: %v", sp.Name, err)
		}
	}

	log.Println("seeding spots...")
	for _, s := range spots {
		var spotID string
		err := pool.QueryRow(ctx, `
			INSERT INTO spots (name, location, type, description, access_notes, parking, difficulty, depth_info, bottom_type)
			VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326), $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT DO NOTHING RETURNING id`,
			s.Name, s.Lon, s.Lat, s.Type, s.Description, s.AccessNotes, s.Parking, s.Difficulty, s.DepthInfo, s.BottomType,
		).Scan(&spotID)
		if err != nil {
			log.Printf("error inserting spot %s: %v", s.Name, err)
			continue
		}

		for _, spName := range s.Species {
			_, err := pool.Exec(ctx, `
				INSERT INTO spot_species (spot_id, species_id)
				SELECT $1, id FROM species WHERE name = $2
				ON CONFLICT DO NOTHING`,
				spotID, spName,
			)
			if err != nil {
				log.Printf("error linking species %s to %s: %v", spName, s.Name, err)
			}
		}
	}

	log.Println("seeding complete!")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
