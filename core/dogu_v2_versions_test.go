package core

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func Test_sortDogus(t *testing.T) {
	unsortedDogus := []*Dogu{{Name: "Dogu1", Version: "11.22.33-1"}, {Name: "Dogu2", Version: "0.1.3-5"}, {Name: "Dogu3", Version: "0.5.3-3"}, {Name: "Dogu4", Version: "9.3.9"}}
	expectedDogus := []*Dogu{{Name: "Dogu1", Version: "11.22.33-1"}, {Name: "Dogu4", Version: "9.3.9"}, {Name: "Dogu3", Version: "0.5.3-3"}, {Name: "Dogu2", Version: "0.1.3-5"}}

	sort.Sort(ByDoguVersion(unsortedDogus))

	assert.Equal(t, unsortedDogus, expectedDogus)
}
