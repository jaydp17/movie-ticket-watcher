package core

import (
	"github.com/jaydp17/movie-ticket-watcher/pkg/dao"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
	"github.com/jaydp17/movie-ticket-watcher/pkg/utils"
	"regexp"
	"strings"
)

// mergeCinemasByName converts the name of the cinema to a slug by removing punctuation, lower casing, etc.
// and then compares to see how many cinemas have the same NameSlug
func mergeCinemasByName(bmsCinemas, ptmCinemas []providers.Cinema) (dao.Cinemas, []providers.Cinema, []providers.Cinema) {
	maxCinemas := utils.MaxInt(len(bmsCinemas), len(ptmCinemas))
	cinemasByName := make(map[string]dao.Cinema, maxCinemas)

	// match using key = NameSlug
	for _, cinema := range bmsCinemas {
		nameSlug := cinema.NameSlug()
		cinemaWithBmsID := dao.Cinema{
			Cinema: providers.Cinema{
				ID:        nameSlug,
				Name:      cinema.Name,
				Provider:  cinema.Provider,
				CityID:    cinema.CityID,
				Latitude:  cinema.Latitude,
				Longitude: cinema.Longitude,
				Address:   cinema.Address,
			},
			BookmyshowID: cinema.ID,
		}
		cinemasByName[nameSlug] = cinemaWithBmsID
	}
	for _, cinema := range ptmCinemas {
		nameSlug := cinema.NameSlug()
		if existingCinema, ok := cinemasByName[nameSlug]; ok {
			existingCinema.PaytmID = cinema.ID
			cinemasByName[nameSlug] = existingCinema
		} else {
			cinemaWithPytmID := dao.Cinema{
				Cinema: providers.Cinema{
					ID:        nameSlug,
					Name:      cinema.Name,
					Provider:  cinema.Provider,
					CityID:    cinema.CityID,
					Latitude:  cinema.Latitude,
					Longitude: cinema.Longitude,
					Address:   cinema.Address,
				},
				PaytmID: cinema.ID,
			}
			cinemasByName[nameSlug] = cinemaWithPytmID
		}
	}

	// generate list of merged cinemas
	mergedCinemas := make([]dao.Cinema, 0, maxCinemas)
	for _, c := range cinemasByName {
		if c.HasAllProviderIDs() {
			mergedCinemas = append(mergedCinemas, c)
		}
	}

	// find the remaining cinemas
	remainingBmsCinemas := make([]providers.Cinema, 0)
	for _, cinema := range bmsCinemas {
		if foundCinema, ok := cinemasByName[cinema.NameSlug()]; ok && foundCinema.HasAllProviderIDs() {
			// it's matched
		} else {
			remainingBmsCinemas = append(remainingBmsCinemas, cinema)
		}
	}
	remainingPtmCinemas := make([]providers.Cinema, 0)
	for _, cinema := range ptmCinemas {
		if foundCinema, ok := cinemasByName[cinema.NameSlug()]; ok && foundCinema.HasAllProviderIDs() {
			// it's matched
		} else {
			remainingPtmCinemas = append(remainingPtmCinemas, cinema)
		}
	}

	return mergedCinemas, remainingBmsCinemas, remainingPtmCinemas
}

// mergeCinemasByGeoDistance clubs cinemas which are geographically within 50m from the other provider's cinema
// and if we find multiple such cinemas, we use a fuzzy search to see which out of them matches the name of the cinema more
// the later part of the function is really useful when we have cinemas like PVR Forum mall, where multiple cinemas are
// in the same Geo location, i.e. on 4th floor there's the normal PVR & on the 5th floor there's PVR Gold
func mergeCinemasByGeoDistance(bmsCinemas, ptmCinemas []providers.Cinema) (dao.Cinemas, []providers.Cinema, []providers.Cinema) {
	bmsMatchedIDs := make(map[string]bool, len(bmsCinemas))
	ptmMatchedIDs := make(map[string]bool, len(ptmCinemas))

	type BmsPtmMatch struct {
		bms      providers.Cinema
		ptm      providers.Cinema
		distance float64
	}
	matches := make([]BmsPtmMatch, 0)

	type matchPlusDistance struct {
		cinema   providers.Cinema
		distance float64
	}
	for _, bCinema := range bmsCinemas {
		minMeters := 50.0 // this is set to 50 because we don't want to match cinemas that are more than 50 meters apart
		ptmMatches := make([]matchPlusDistance, 0)
		for _, pCinema := range ptmCinemas {
			if _, ok := ptmMatchedIDs[pCinema.ID]; ok {
				// this paytmID is already matched with some other bookmyShowID so let's move on
				continue
			}
			meters := utils.GeoDistance(float64(bCinema.Latitude), float64(bCinema.Longitude), float64(pCinema.Latitude), float64(pCinema.Longitude))
			if meters < minMeters {
				// if two cinemas are < 50m apart lets push it into an array
				// which we then compare by their names using a fuzzy string compare
				ptmMatches = append(ptmMatches, matchPlusDistance{cinema: pCinema, distance: meters})
			}
		}

		// if there's just 1 cinema matched then let's just use that
		if len(ptmMatches) == 1 {
			bmsMatchedIDs[bCinema.ID] = true
			ptmCinema := ptmMatches[0].cinema
			ptmMatchedIDs[ptmCinema.ID] = true
			matches = append(matches, BmsPtmMatch{
				bms:      bCinema,
				ptm:      ptmCinema,
				distance: ptmMatches[0].distance,
			})
		}

		// if there's more than 1 cinema matched
		// let's use fuzzy string match to match their names & select one out of the nearby cinemas
		if len(ptmMatches) > 1 {
			bmsMatchedIDs[bCinema.ID] = true
			maxMatchScore := 0.0
			maxMatchIndex := -1
			for index, match := range ptmMatches {
				score := matchScore(bCinema.Name, match.cinema.Name)
				if score > maxMatchScore {
					maxMatchScore = score
					maxMatchIndex = index
				}
			}
			if maxMatchIndex > -1 {
				matchedCinema := ptmMatches[maxMatchIndex].cinema
				ptmMatchedIDs[matchedCinema.ID] = true
				matches = append(matches, BmsPtmMatch{
					bms:      bCinema,
					ptm:      matchedCinema,
					distance: ptmMatches[maxMatchIndex].distance,
				})
			}
		}
	}

	mergedCinemas := make([]dao.Cinema, 0, len(matches))
	for _, match := range matches {
		cinema := dao.Cinema{
			Cinema: providers.Cinema{
				ID:        match.bms.NameSlug(),
				Name:      match.bms.Name,
				Provider:  match.bms.Provider,
				CityID:    match.bms.CityID,
				Latitude:  match.bms.Latitude,
				Longitude: match.bms.Longitude,
				Address:   match.bms.Address,
			},
			BookmyshowID: match.bms.ID,
			PaytmID:      match.ptm.ID,
		}
		mergedCinemas = append(mergedCinemas, cinema)
	}

	// find the remaining cinemas
	remainingBmsCinemas := make([]providers.Cinema, 0)
	for _, cinema := range bmsCinemas {
		if _, ok := bmsMatchedIDs[cinema.ID]; !ok {
			remainingBmsCinemas = append(remainingBmsCinemas, cinema)
		}
	}
	remainingPtmCinemas := make([]providers.Cinema, 0)
	for _, cinema := range ptmCinemas {
		if _, ok := ptmMatchedIDs[cinema.ID]; !ok {
			remainingPtmCinemas = append(remainingPtmCinemas, cinema)
		}
	}

	return mergedCinemas, remainingBmsCinemas, remainingPtmCinemas
}

// matchScore tokenizes the strings and then calculates how many tokens are common across two strings
// the 2 strings having the highest common ratio win
func matchScore(a, b string) float64 {
	aTokens := tokenize(a)
	bTokens := tokenize(b)

	if len(aTokens) < len(bTokens) {
		aTokens, bTokens = bTokens, aTokens
	}

	matchCount := 0
	for _, aTok := range aTokens {
		for _, bTok := range bTokens {
			if aTok == bTok {
				matchCount++
			}
		}
	}

	score := float64(2*matchCount) / float64(len(aTokens)+len(bTokens))
	return score
}

// tokenize is used split a string into normalized tokens
func tokenize(str string) []string {
	reg := regexp.MustCompile(`[^0-9a-zA-z ]+`)
	cleanStr := reg.ReplaceAllString(str, "")
	splits := strings.Split(cleanStr, " ")
	lowerSplits := make([]string, 0)
	for _, s := range splits {
		if len(s) > 0 {
			lowerSplits = append(lowerSplits, strings.ToLower(s))
		}
	}
	return lowerSplits
}
