package star

import (
	"fmt"
	"math"
	"math/rand"
)

func createRandomID(namespace string) string {
	k := rand.Intn(math.MaxInt64-1000000000000000000) + 1000000000000000000
	return fmt.Sprintf("%s-%d", namespace, k)
}
