package wildcard

import (
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/pkg/maths"
	"github.com/MyCode83/godirb/pkg/random"
	"strings"

	"github.com/valyala/fasthttp"
)

// random string

// Wildcard struct
type Wildcard struct {
	Status    int
	Lenght    int
	Active    bool
	Tolerance int
}

// Wildcard funcion
func DetectWildcard(client *fasthttp.Client, baseURL string, Placeholder string, userAgents ...string) (*Wildcard, error) {
	var status int
	var lenght int
	const tries = 3
	debug.Printf("wildcard detection start base_url=%q placeholder=%q tries=%d", baseURL, Placeholder, tries)
	var fullURL string
	var isPortFuzz bool = false
	var lenghts []int
	var tolerance int
	for index := range tries {
		if strings.Contains(baseURL, Placeholder) && !isPortFuzz {
			fullURL = strings.Replace(baseURL, Placeholder, random.RandomString(12), 1)
		} else {
			fullURL = baseURL + "/" + random.RandomString(12)
		}

		request := fasthttp.AcquireRequest()
		response := fasthttp.AcquireResponse()
		request.SetRequestURI(fullURL)
		request.Header.SetUserAgent(random.RandChoice(userAgents))
		request.Header.SetMethod("GET")

		debug.Request("wildcard", request)
		err := client.Do(request, response)
		fasthttp.ReleaseRequest(request)
		if err != nil {
			debug.Error("wildcard", err)
			fasthttp.ReleaseResponse(response)
			return nil, err
		}
		debug.Response("wildcard", response)
		bodyLen := len(response.Body())
		resStatusCode := response.StatusCode()
		if index == 0 {
			status = resStatusCode
			lenght = bodyLen
			debug.Printf("wildcard first sample status=%d length=%d", status, lenght)
		} else {
			if status == 404 {
				debug.Printf("wildcard inactive because first status is 404")
				fasthttp.ReleaseResponse(response)
				return &Wildcard{Active: false}, nil
			}
			if status != resStatusCode || lenght != bodyLen {
				debug.Printf("wildcard inactive mismatch index=%d status=%d length=%d expected_status=%d expected_length=%d",
					index, resStatusCode, bodyLen, status, lenght)
				fasthttp.ReleaseResponse(response)
				return &Wildcard{Active: false}, nil
			}
		}
		lenghts = append(lenghts, bodyLen)
		fasthttp.ReleaseResponse(response)

	}
	min, max := maths.MinMax(lenghts...)
	tolerance = max - min
	if tolerance <= 0 {
		tolerance = 10
	}
	debug.Printf("wildcard active status=%d length=%d tolerance=%d samples=%v", status, lenght, tolerance, lenghts)

	return &Wildcard{
		Status:    status,
		Lenght:    lenght,
		Active:    true,
		Tolerance: tolerance,
	}, nil
}
