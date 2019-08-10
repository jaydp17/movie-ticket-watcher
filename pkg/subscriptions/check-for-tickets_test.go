package subscriptions

import (
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"testing"
	"time"
)

func TestGroupSimilarSubscriptions(t *testing.T) {
	date := db.YYYYMMDDTime{Time: time.Now()}
	input := []Subscription{
		{
			CityID:        "a",
			MovieID:       "b",
			ScreeningDate: date,
		},
		{
			CityID:        "a",
			MovieID:       "b",
			ScreeningDate: date,
		},
		{
			CityID:        "b",
			MovieID:       "c",
			ScreeningDate: date,
		},
	}
	expectedOutput := []groupSubscriptions{
		{
			subscriptions: []Subscription{
				{
					CityID:        "a",
					MovieID:       "b",
					ScreeningDate: date,
				},
				{
					CityID:        "a",
					MovieID:       "b",
					ScreeningDate: date,
				},
			},
		},
		{
			subscriptions: []Subscription{
				{
					CityID:        "b",
					MovieID:       "c",
					ScreeningDate: date,
				},
			},
		},
	}

	result := groupSimilarSubscriptions(input)
	if len(result) != len(expectedOutput) {
		t.Errorf("expectedOutput(%d groups) & result(%d groups) don't have the same number of groups", len(expectedOutput), len(result))
	}

	for i := 0; i < len(result); i++ {
		resultGroup := result[i].subscriptions
		expectedGroup := expectedOutput[i].subscriptions
		if len(resultGroup) != len(expectedGroup) {
			t.Errorf("subscriptions in the groups don't match with the expected output")
		}
		for j := 0; j < len(resultGroup); j++ {
			if resultGroup[j].CityID != expectedGroup[j].CityID || resultGroup[j].MovieID != expectedGroup[j].MovieID || resultGroup[j].ScreeningDate.ToYYYYMMDD() != expectedGroup[j].ScreeningDate.ToYYYYMMDD() {
				t.Errorf("subscriptions in the groups don't match")
			}
		}
	}
}
