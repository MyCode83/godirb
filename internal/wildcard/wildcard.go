package wildcard

import (
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/MyCode83/godirb/pkg/maths"
	"github.com/MyCode83/godirb/pkg/random"
	"strings"
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
func DetectWildcard(client *transport.Client, baseURL string, Placeholder string, userAgents ...string) (*Wildcard, error) {
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

		response, err := client.Do(&transport.RequestOptions{
			URL:        fullURL,
			Method:     transport.MethodGET,
			MethodMode: transport.MethodModeFixed,
			UserAgent:  random.RandChoice(userAgents),
		})
		if err != nil {
			debug.Error("wildcard", err)
			return nil, err
		}
		debug.Printf("wildcard response status=%d body=%d", response.StatusCode, response.Lenght)
		bodyLen := response.Lenght
		resStatusCode := response.StatusCode
		if index == 0 {
			status = resStatusCode
			lenght = bodyLen
			debug.Printf("wildcard first sample status=%d length=%d", status, lenght)
		} else {
			if status == 404 {
				debug.Printf("wildcard inactive because first status is 404")
				return &Wildcard{Active: false}, nil
			}
			if status != resStatusCode || lenght != bodyLen {
				debug.Printf("wildcard inactive mismatch index=%d status=%d length=%d expected_status=%d expected_length=%d",
					index, resStatusCode, bodyLen, status, lenght)
				return &Wildcard{Active: false}, nil
			}
		}
		lenghts = append(lenghts, bodyLen)

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
