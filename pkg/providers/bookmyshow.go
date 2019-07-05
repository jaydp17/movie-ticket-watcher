package providers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// BookMyShow struct
type BookMyShow struct{}

var client = &http.Client{}

const (
	okHTTPUserAgent = "okhttp/3.11.0"
	token           = "67x1xa33b4x422b361ba"
)

// GetAvailableVenueCodes fetches all the venue codes where the given movie is available
func (b *BookMyShow) GetAvailableVenueCodes(childEventCode, regionCode, dateStr string) {
	url := "https://in.bookmyshow.com/api/v2/mobile/showtimes/byevent"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", okHTTPUserAgent)
	q := req.URL.Query()
	q.Add("regionCode", regionCode)
	q.Add("eventCode", childEventCode)
	q.Add("token", token)
	q.Add("bmsId", "1.82650383.1552055894719")
	q.Add("dateCode", dateStr)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	responseStr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(responseStr))
}
