package wildcard

import (

	"godirb/pkg/random"
	"godirb/pkg/maths"
	"strings"

	"github.com/valyala/fasthttp"
)

// random string

// Wildcard struct
type Wildcard struct {
	Status int
	Lenght int
	Active bool
	Tolerance int
}
// Wildcard funcion
func DetectWildcard(client *fasthttp.Client, baseURL string, Placeholder string,userAgents ...string) (*Wildcard, error) {
	var status int
	var lenght int
	const tries = 3
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
		
		err := client.Do(request, response)
		fasthttp.ReleaseRequest(request)
		if err != nil {
			fasthttp.ReleaseResponse(response)
			return nil, err
		}
		bodyLen := len(response.Body())
		resStatusCode := response.StatusCode()
		if index == 0 {
			status = resStatusCode
			lenght = bodyLen
		} else {
			if status == 404 {
				fasthttp.ReleaseResponse(response)
				return &Wildcard{Active: false}, nil
			}
			if status != resStatusCode || lenght != bodyLen {
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

	return &Wildcard{
		Status: status,
		Lenght: lenght,
		Active: true,
		Tolerance: tolerance,
	}, nil
}