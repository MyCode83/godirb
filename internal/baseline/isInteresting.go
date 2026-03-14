package baseline
import (
	"godirb/internal/wildcard"
)
func (b *Baseline) IsInteresting(status int, lenght int, tolerance int) bool {
	if status != b.Status {
		return true
	}

	if !wildcard.IsSimilarSize(lenght, b.Lenght, tolerance){
		return true
	}
	return false
}