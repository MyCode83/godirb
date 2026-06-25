package baseline

import (
	"fmt"
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/MyCode83/godirb/pkg/maths"
	"github.com/MyCode83/godirb/pkg/random"
	"strings"
)

type Baseline struct {
	Status    int
	Lenght    int
	Tolerance int
}

func BuildBaseLine(baseUrl string, client *transport.Client, placeholder string) (*Baseline, error) {
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
	for i := range tries {
		randomText := random.RandomString(12)
		url := parts[0] + randomText + parts[1]

		response, err := client.Do(&transport.RequestOptions{
			URL:        url,
			Method:     transport.MethodGET,
			MethodMode: transport.MethodModeFixed,
		})
		if err != nil {
			debug.Error("baseline", err)
			return nil, fmt.Errorf("[!] %w\n", err)
		}
		debug.Printf("baseline response status=%d body=%d", response.StatusCode, response.Lenght)
		if i == 1 {
			status = response.StatusCode
			lenght = response.Lenght
			debug.Printf("baseline reference sample status=%d length=%d", status, lenght)
		}
		lenghts = append(lenghts, response.Lenght)

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
