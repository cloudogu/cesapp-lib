package core

import (
	"fmt"
	"github.com/gammazero/toposort"
	"sort"
)

// SortDogusByDependency takes an unsorted slice of Dogu structs and returns a slice of Dogus ordered by the
// importance of their dependencies descending, that is: the most needed dogu will be the first element.
func SortDogusByDependency(dogus []*Dogu) []*Dogu {
	ordered := sortByDependency{dogus}
	return ordered.sortDogus()
}

// SortDogusByInvertedDependency takes an unsorted slice of Dogu structs and returns a new slice of Dogus ordered by the
// importance of their dependencies ascending, that is: the most independent dogu will be the first element.
func SortDogusByInvertedDependency(dogus []*Dogu) []*Dogu {
	orderedDesc := SortDogusByDependency(dogus)

	orderedAsc := []*Dogu{}
	for i := len(orderedDesc) - 1; i >= 0; i-- {
		orderedAsc = append(orderedAsc, orderedDesc[i])
	}

	return orderedAsc
}

// SortDogusByName takes an unsorted slice of Dogus
// and returns an copy of the slice ordered by the full name of the dogu
func SortDogusByName(dogus []*Dogu) []*Dogu {
	byName := sortByName{dogus}
	sort.Sort(&byName)
	return byName.dogus
}

// private struct which is required of ordering

type sortByName struct {
	dogus []*Dogu
}

func (byName *sortByName) Len() int {
	return len(byName.dogus)
}

func (byName *sortByName) Swap(i, j int) {
	byName.dogus[i], byName.dogus[j] = byName.dogus[j], byName.dogus[i]
}

func (byName *sortByName) Less(i, j int) bool {
	return byName.dogus[i].GetFullName() < byName.dogus[j].GetFullName()
}

type sortByDependency struct {
	dogus []*Dogu
}

func contains(slice []Dependency, item string) bool {
	for _, s := range slice {
		if s.Name == item {
			return true
		}
	}
	return false
}

func (bd *sortByDependency) sortDogus() []*Dogu {
	dependencyEdges := bd.getDependencyEdges()
	sorted, err := toposort.Toposort(dependencyEdges)
	if err != nil {
		return nil
	}

	sortedDogus, err := toDoguSlice(sorted)
	if err != nil {
		return nil
	}

	return sortedDogus
}

func (bd *sortByDependency) getDependencyEdges() []toposort.Edge {
	var dependencyEdges []toposort.Edge
	for _, dogu := range bd.dogus {
		dependencies := dogu.GetAllDependenciesOfType(DependencyTypeDogu)
		if len(dependencies) > 0 {
			dependentDogus := bd.dependenciesToDogus(dependencies)
			for _, dependency := range dependentDogus {
				dependencyEdges = append(dependencyEdges, toposort.Edge{dependency, dogu})
			}
		} else {
			dependencyEdges = append(dependencyEdges, toposort.Edge{nil, dogu})
		}
	}
	return dependencyEdges
}

func toDoguSlice(dogus []interface{}) ([]*Dogu, error) {
	result := make([]*Dogu, len(dogus))
	for i, dogu := range dogus {
		if castedDogu, ok := dogu.(*Dogu); ok {
			result[i] = castedDogu
		} else {
			return nil, fmt.Errorf("expected Dogu, got %T", dogu)
		}
	}
	return result, nil
}

func (bd *sortByDependency) dependenciesToDogus(dependencies []Dependency) []*Dogu {
	var result []*Dogu

	for _, dogu := range bd.dogus {
		if contains(dependencies, dogu.GetSimpleName()) {
			result = append(result, dogu)
		}
	}

	return result
}
