package core

import (
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
	"testing"
)

var bmsCinemas = []providers.Cinema{
	{
		ID:        "PVDR",
		Name:      "PVR: 4DX, Orion Mall, Dr Rajkumar Road",
		Provider:  "PVR",
		Latitude:  13.011,
		Longitude: 77.5551,
		Address:   "Orion Mall, 3rd Floor, Brigade Gateway, Dr. Rajkumar Road, Malleshwaram West, Bengaluru, Karnataka 560055, India",},
	{
		ID:        "PMDX",
		Name:      "PVR: 4DX, Phoenix Marketcity Mall, Whitefield Road",
		Provider:  "PVR",
		Latitude:  12.9973,
		Longitude: 77.6957,
		Address:   "Phoenix Market City Mall, Whitefield Road, Mahadevapura, Krishnarajapura, Devasandra Industrial Estate, Bengaluru, Karnataka 560048, India",},
	{
		ID:        "PVDX",
		Name:      "PVR: 4DX, Vega City, Bannerghatta Road, Bengaluru",
		Provider:  "PVR",
		Latitude:  12.9077,
		Longitude: 77.6012,
		Address:   "PVR Cinemas, 4th and 5th Floor, Vega City Mall,172/1, Bannerghatta Main Road, JP Nagar 3rd Phase, Dollar Layout, Bengaluru, Karnataka 560076, India",},
	{
		ID:        "CXSS",
		Name:      "PVR: Arena Mall, Doddanakundi, Marthalli Ring Road",
		Provider:  "PVR",
		Latitude:  12.9799,
		Longitude: 77.6932,
		Address:   "5th Floor, Soul Space Arena Mall, Akme Ballet Inner Road, Doddanakundi, Marathahalli Outer Ring Road, Mahadevapura, Bengaluru, Karnataka 560037, India",},
	{
		ID:        "CXBL",
		Name:      "PVR: Central Spirit Mall, Bellandur",
		Provider:  "PVR",
		Latitude:  12.9258,
		Longitude: 77.6752,
		Address:   "Survey No.89/6, 78/7, Soul Space Spirit, Outer Ring Road, Bellandur Junction, Opposite Acme Hormony Plaza, Bengaluru, Karnataka 560103, India",},
	{
		ID:        "PVBG",
		Name:      "PVR: Gold, Forum Mall, Koramangala",
		Provider:  "PVR",
		Latitude:  12.9346,
		Longitude: 77.6111,
		Address:   "The Forum Mall, 21-22, Adugodi Main Road, Koramangala, Chikku Lakshmaiah Layout, Bengaluru, Karnataka 560095, India",},
	{
		ID:        "PVOG",
		Name:      "PVR: Gold, Orion Mall, Dr Rajkumar Road",
		Provider:  "PVR",
		Latitude:  13.011,
		Longitude: 77.5551,
		Address:   "Orion Mall, 3rd Floor, Brigade Gateway, Dr. Rajkumar Road, Malleshwaram West, Bengaluru, Karnataka 560055, India",},
	{
		ID:        "PVGD",
		Name:      "PVR: Gold, Vega City, Bannerghatta Road",
		Provider:  "PVR",
		Latitude:  12.9077,
		Longitude: 77.6012,
		Address:   "PVR Cinemas, 4th and 5th Floor, Vega City Mall,172/1, Bannerghatta Main Road, JP Nagar 3rd Phase, Dollar Layout, Bengaluru, Karnataka 560076, India",},
	{
		ID:        "PVRB",
		Name:      "PVR: MSR Elements Mall, Tanisandhra Main Road",
		Provider:  "PVR",
		Latitude:  13.0452,
		Longitude: 77.6266,
		Address:   "5th Floor, MSR Elements Regalia Mall, 100 Main Tanisandra Road, Kasba Hobli, Nagavara, Bengaluru, Karnataka 560077, India",},
	{
		ID:        "PVOR",
		Name:      "PVR: Orion Mall, Dr Rajkumar Road",
		Provider:  "PVR",
		Latitude:  13.011,
		Longitude: 77.5551,
		Address:   "Orion Mall, 3rd Floor, Brigade Gateway, Dr. Rajkumar Road, Malleshwaram West, Bengaluru, Karnataka 560055, India",},
	{
		ID:        "PVBM",
		Name:      "PVR: Phoenix Marketcity Mall, Whitefield Road",
		Provider:  "PVR",
		Latitude:  12.9973,
		Longitude: 77.6957,
		Address:   "Phoenix Market City Mall, Whitefield Road, Mahadevapura, Krishnarajapura, Devasandra Industrial Estate, Bengaluru, Karnataka 560048, India",},
	{
		ID:        "PVLY",
		Name:      "PVR: Play House, Vega City, Bannerghatta Road",
		Provider:  "PVR",
		Latitude:  12.9077,
		Longitude: 77.6012,
		Address:   "PVR Cinemas, 4th and 5th Floor, Vega City Mall,172/1, Bannerghatta Main Road, JP Nagar 3rd Phase, Dollar Layout, Bengaluru, Karnataka 560076, India",},
	{
		ID:        "PPPM",
		Name:      "PVR: Play House,Phoenix Marketcity Whitefield Road",
		Provider:  "PVR",
		Latitude:  12.9973,
		Longitude: 77.6957,
		Address:   "Phoenix Market City Mall, Whitefield Road, Mahadevapura, Krishnarajapura, Devasandra Industrial Estate, Bengaluru, Karnataka 560048, India",},
	{
		ID:        "PVVS",
		Name:      "PVR: Vaishnavi Sapphire Mall, Yeshwanthpur",
		Provider:  "PVR",
		Latitude:  13.0247,
		Longitude: 77.5486,
		Address:   "3rd floor, Vaishnavi Sapphire Mall, Tumkur Main Road, Yeshwanthpur Industrial Area, Phase 1, Yeshwanthpur, Bengaluru, Karnataka 560022, India",},
	{
		ID:        "PVEG",
		Name:      "PVR: Vega City,  Bannerghatta Road",
		Provider:  "PVR",
		Latitude:  12.9077,
		Longitude: 77.6012,
		Address:   "PVR Cinemas, 4th and 5th Floor, Vega City Mall,172/1, Bannerghatta Main Road, JP Nagar 3rd Phase, Dollar Layout, Bengaluru, Karnataka 560076, India",},
	{
		ID:        "PVPB",
		Name:      "PVR: VR Bengaluru, Whitefield Road",
		Provider:  "PVR",
		Latitude:  12.9965,
		Longitude: 77.6953,
		Address:   "3rd Floor, VR Mall, 408 3rd Cross, Singayana Palya, Main Whitefield Road, Mahadevpura, Bengaluru, Karnataka 560048, India",},
	{
		ID:        "PVVR",
		Name:      "PVR: VR GOLD, VR Bengaluru, Whitefield Road",
		Provider:  "PVR",
		Latitude:  12.9965,
		Longitude: 77.6953,
		Address:   "3rd Floor, VR Mall, 408 3rd Cross, Singayana Palya, Main Whitefield Road, Mahadevpura, Bengaluru, Karnataka 560048, India",},
}
var ptmCinemas = []providers.Cinema{
	{
		ID:        "7261",
		Name:      "PVR Play House Phoenix Marketcity Mall, Whitefield Road",
		Provider:  "PVR",
		Latitude:  12.997251,
		Longitude: 77.69573,
		Address:   "Phoenix Market City Mall, Whitefield Road, Mahadevapura, Krishnarajapura, Bengaluru, Karnataka 560048, India",
	},
	{
		ID:        "207",
		Name:      "PVR Gold Forum Mall, Koramangala",
		Provider:  "PVR",
		Latitude:  12.934661,
		Longitude: 77.611314,
		Address:   "The Forum Mall, 21-22, Adugodi Main Road, Chikku Lakshmaiah Layout, Koramangala, Bengaluru, Karnataka 560095, India",
	},
	{
		ID:        "208",
		Name:      "PVR Gold Orion Mall, Dr Rajkumar Road",
		Provider:  "PVR",
		Latitude:  13.011014,
		Longitude: 77.555058,
		Address:   "3rd Floor, Orion Mall, Brigade Gateway, Dr. Rajkumar Road, Malleshwaram West, Bengaluru, Karnataka 560055, India",
	},
	{
		ID:        "2845",
		Name:      "PVR Phoenix Marketcity Mall 4DX , Whitefield Road",
		Provider:  "PVR",
		Latitude:  12.997251,
		Longitude: 77.69573,
		Address:   "Phoenix Market City Mall, Whitefield Road, Mahadevapura, Krishnarajapura, Bengaluru, Karnataka 560048, India",
	},
	{
		ID:        "3672",
		Name:      "PVR VEGA 4DX, Bannerghatta Road",
		Provider:  "PVR",
		Latitude:  12.907664,
		Longitude: 77.601161,
		Address:   "4th and 5th Floor, Vega City Mall, 172/1, Bannerghatta Main Road, Dollar Layout, JP Nagar 3rd Phase, Bengaluru, Karnataka 560076, India",
	},
	{
		ID:        "277",
		Name:      "PVR Forum IMAX, Koramangala",
		Provider:  "PVR",
		Latitude:  12.934661,
		Longitude: 77.611314,
		Address:   "The Forum Mall, 21-22, Adugodi Main Road, Chikku Lakshmaiah Layout, Koramangala, Bengaluru, Karnataka 560095, India",
	},
	{
		ID:        "209",
		Name:      "PVR Orion Mall, Dr Rajkumar Road",
		Provider:  "PVR",
		Latitude:  13.011014,
		Longitude: 77.555058,
		Address:   "3rd Floor, Orion Mall, Brigade Gateway, Dr. Rajkumar Road, Malleshwaram West, Bengaluru, Karnataka 560055, India",
	},
	{
		ID:        "205",
		Name:      "PVR Central Spirit Mall, Bellandur",
		Provider:  "PVR",
		Latitude:  12.925817,
		Longitude: 77.67521,
		Address:   "Survey No.89/6, 78/7, Soul Space Spirit, Outer Ring Road, Opposite Acme Hormony Plaza, Bellandur Junction, Bengaluru, Karnataka 560103, India",
	},
	{
		ID:        "211",
		Name:      "PVR Phoenix Marketcity Mall, Whitefield Road",
		Provider:  "PVR",
		Latitude:  12.997251,
		Longitude: 77.69573,
		Address:   "Phoenix Market City Mall, Whitefield Road, Mahadevapura, Krishnarajapura, Bengaluru, Karnataka 560048, India",
	},
	{
		ID:        "193",
		Name:      "PVR Vaishnavi Sapphire Mall, Yashwantpur",
		Provider:  "PVR",
		Latitude:  13.024686,
		Longitude: 77.548585,
		Address:   "3rd floor, Vaishnavi Sapphire Mall, Tumkur Main Road, Yeshwanthpur Industrial Area Phase 1, Yeshwanthpur, Bengaluru, Karnataka 560022, India",
	},
	{
		ID:        "304",
		Name:      "PVR VR GOLD, Whitefield Road",
		Provider:  "PVR",
		Latitude:  12.996519,
		Longitude: 77.695256,
		Address:   "3rd Floor, VR Mall, Singayana Palya, Main Whitefield Road, Mahadevpura, Bengaluru, Karnataka 560048, India",
	},
	{
		ID:        "306",
		Name:      "PVR VR, Whitefield Road",
		Provider:  "PVR",
		Latitude:  12.996519,
		Longitude: 77.695256,
		Address:   "3rd Floor, VR Mall, Singayana Palya, Main Whitefield Road, Mahadevpura, Bengaluru, Karnataka 560048, India",
	},
	{
		ID:        "3674",
		Name:      "PVR Paytm IMAX Vega City, Bannerghatta Road",
		Provider:  "PVR",
		Latitude:  12.907664,
		Longitude: 77.601161,
		Address:   "4th and 5th Floor, Vega City Mall, 172/1, Bannerghatta Main Road, Dollar Layout, JP Nagar 3rd Phase, Bengaluru, Karnataka 560076, India",
	},
	{
		ID:        "204",
		Name:      "PVR Arena Mall, Doddanakundi",
		Provider:  "PVR",
		Latitude:  12.979911,
		Longitude: 77.693189,
		Address:   "5th Floor, Soul Space Arena Mall, Akme Ballet Inner Road, Doddanakundi, Marathahalli Outer Ring Road, Mahadevapura, Bengaluru, Karnataka 560037, India",
	},
	{
		ID:        "4945",
		Name:      "PVR 4DX Orion Mall , Dr Rajkumar Road",
		Provider:  "PVR",
		Latitude:  13.011014,
		Longitude: 77.555058,
		Address:   "PVR Orion Mall, 3rd Floor, Brigade Gateway, Malleshwaram (East), Bengaluru, Karnataka 560001, India",
	},
	{
		ID:        "3671",
		Name:      "PVR VEGA, Bannerghatta Road",
		Provider:  "PVR",
		Latitude:  12.907664,
		Longitude: 77.601161,
		Address:   "4th and 5th Floor, Vega City Mall, 172/1, Bannerghatta Main Road, Dollar Layout, JP Nagar 3rd Phase, Bengaluru, Karnataka 560076, India",
	},
	{
		ID:        "7145",
		Name:      "PVR Play House Forum Mall, Koramangala",
		Provider:  "PVR",
		Latitude:  12.934661,
		Longitude: 77.611314,
		Address:   "The Forum Mall, 21-22, Adugodi Main Road, Chikku Lakshmaiah Layout, Koramangala, Bengaluru, Karnataka 560095, India",
	},
	{
		ID:        "206",
		Name:      "PVR Forum Mall, Koramangala",
		Provider:  "PVR",
		Latitude:  12.934661,
		Longitude: 77.611314,
		Address:   "The Forum Mall, 21-22, Adugodi Main Road, Chikku Lakshmaiah Layout, Koramangala, Bengaluru, Karnataka 560095, India",
	},
	{
		ID:        "287",
		Name:      "PVR MSR Elements Mall, Tanisandhra Main Road",
		Provider:  "PVR",
		Latitude:  13.045163,
		Longitude: 77.62658,
		Address:   "5th Floor, MSR Elements Regalia Mall, 100 Main Tanisandra Road, Kasba Hobli, Nagavara, Bengaluru, Karnataka 560077, India",
	},
	{
		ID:        "305",
		Name:      "PVR VR IMAX, Whitefield Road",
		Provider:  "PVR",
		Latitude:  12.996519,
		Longitude: 77.695256,
		Address:   "3rd Floor, VR Mall, Singayana Palya, Main Whitefield Road, Mahadevpura, Bengaluru, Karnataka 560048, India",
	},
	{
		ID:        "7647",
		Name:      "PVR 4DX Forum Mall, Koramangala",
		Provider:  "PVR",
		Latitude:  12.934661,
		Longitude: 77.611314,
		Address:   "The Forum Mall, 21-22, Adugodi Main Road, Chikku Lakshmaiah Layout, Koramangala, Bengaluru, Karnataka 560095, India",
	},
	{
		ID:        "3675",
		Name:      "PVR VEGA Play House, Bannerghatta Road",
		Provider:  "PVR",
		Latitude:  12.907664,
		Longitude: 77.601161,
		Address:   "4th and 5th Floor, Vega City Mall, 172/1, Bannerghatta Main Road, Dollar Layout, JP Nagar 3rd Phase, Bengaluru, Karnataka 560076, India",
	},
	{
		ID:        "3474",
		Name:      "PVR PXL Orion Mall, Dr Rajkumar Road",
		Provider:  "PVR",
		Latitude:  13.011014,
		Longitude: 77.555058,
		Address:   "3rd Floor, Orion Mall, Brigade Gateway, Dr. Rajkumar Road, Malleshwaram West, Bengaluru, Karnataka 560055, India",
	},
	{
		ID:        "3673",
		Name:      "PVR VEGA Gold, Bannerghatta Road",
		Provider:  "PVR",
		Latitude:  12.907664,
		Longitude: 77.601161,
		Address:   "4th and 5th Floor, Vega City Mall, 172/1, Bannerghatta Main Road, Dollar Layout, JP Nagar 3rd Phase, Bengaluru, Karnataka 560076, India",
	},
}
var bmsPtmMatchCinemas = map[string]string{
	"PVOG": "208",
	"PVBM": "211",
	"PVDR": "4945",
	"CXBL": "205",
	"PVRB": "287",
	"PVOR": "209",
	"PVBG": "207",
	"PMDX": "2845",
	"PVDX": "3672",
	"CXSS": "204",
	"PVGD": "3673",
	"PVLY": "3675",
	"PPPM": "7261",
	"PVVS": "193",
	"PVEG": "3671",
	"PVPB": "306",
	"PVVR": "304",
}

func TestMergeCinemas(t *testing.T) {
	merged := MergeCinemas(bmsCinemas, ptmCinemas)

	expectedMatches := len(bmsPtmMatchCinemas)
	var actualMatches int
	for _, val := range merged {
		if val.HasAllProviderIDs() {
			actualMatches++
		}
	}
	if actualMatches != expectedMatches {
		t.Errorf("expected matches: %d, actual matches: %d", expectedMatches, actualMatches)
	}

	for _, val := range merged {
		if !val.HasAllProviderIDs() {
			continue
		}
		expectedPtmID := bmsPtmMatchCinemas[val.BookmyshowID]
		if expectedPtmID != val.PaytmID {
			t.Errorf("expected bmsID (%s) to match with ptmID (%s), but got ptmID (%s)", val.BookmyshowID, expectedPtmID, val.PaytmID)
		}
	}
}
