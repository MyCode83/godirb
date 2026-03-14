package baseline

import (
	"fmt"
	"godirb/pkg/random"
	"godirb/pkg/maths"
	"strings"

	"github.com/valyala/fasthttp"
)
type Baseline struct {
	Status int
	Lenght int
	Tolerance int
}
func BuildBaseLine(baseUrl string, client *fasthttp.Client ,placeholder string) (*Baseline, error){
	const tries = 3
	var status int
	var lenght int

	var lenghts []int
	var tolerance int

	parts := strings.Split(baseUrl, placeholder)
	if len(parts) != 2 {
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

		err := client.Do(request, response)
		if err != nil {
			return nil, fmt.Errorf("[!] %w\n", err)
		}
		if i == 1 {
			status = response.StatusCode()
			lenght = len(response.Body())
		}
		lenghts = append(lenghts, len(response.Body()))


	}
	min, max := maths.MinMax(lenghts...)
	tolerance = max - min
	if tolerance <= 0 {
		tolerance = 10
	}
	baseline := &Baseline{
		Status: status,
		Lenght: lenght,
		Tolerance: tolerance,
	}
	return baseline, nil
}