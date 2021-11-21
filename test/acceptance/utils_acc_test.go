package acceptance

import (
	"log"
	"net/http"
)

func tstPerformGet(relativeUrlWithLeadingSlash string) http.Response {
	request, err := http.NewRequest(http.MethodGet, ts.URL+relativeUrlWithLeadingSlash, nil)
	if err != nil {
		log.Fatal(err)
	}
	// create a client that doesn't follow redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return *response
}

