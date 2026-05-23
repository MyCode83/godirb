package baseline

import (
	"fmt"
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/pkg/maths"
	"github.com/MyCode83/godirb/pkg/random"
	"strings"

	"github.com/valyala/fasthttp"
)

type Baseline struct {
	Status    int
	Lenght    int
	Tolerance int
}

func BuildBaseLine(baseUrl string, client *fasthttp.Client, placeholder string) (*Baseline, error) {
	const tries = 3
	debug.Printf("baseline start base_url=%q placeholder=%q tries=%d", baseUrl, placeholder, tries)
	var status int
	var lenght int

	var lenghts []int
	var tolerance int

	parts := strings.Split(baseUrl, placeholder)
	if len(parts) != 2 {
		debug.Printf("baseline invalid placeholder occurrences=%d", len(parts)-1)
		return nil, fmt.Errorf("[!] Placeholder must appear only once\n")
	}
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)
	for i := range tries {
		request.Reset()
		response.Reset()
		randomText := random.RandomString(12)
		url := parts[0] + randomText + parts[1]

		request.SetRequestURI(url)

		debug.Request("baseline", request)
		err := client.Do(request, response)
		if err != nil {
			debug.Error("baseline", err)
			return nil, fmt.Errorf("[!] %w\n", err)
		}
		debug.Response("baseline", response)
		if i == 1 {
			status = response.StatusCode()
			lenght = len(response.Body())
			debug.Printf("baseline reference sample status=%d length=%d", status, lenght)
		}
		lenghts = append(lenghts, len(response.Body()))

	}
	min, max := maths.MinMax(lenghts...)
	tolerance = max - min
	if tolerance <= 0 {
		tolerance = 10
	}
	baseline := &Baseline{
		Status:    status,
		Lenght:    lenght,
		Tolerance: tolerance,
	}
	debug.Printf("baseline ready status=%d length=%d tolerance=%d samples=%v", status, lenght, tolerance, lenghts)
	return baseline, nil
}
