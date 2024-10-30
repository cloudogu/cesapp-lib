package core

import (
	"fmt"
	"github.com/gammazero/toposort"
	"sort"
)

var (
	nginxIngressDependency = Dependency{
		Type: DependencyTypeDogu,
		Name: "nginx-ingress",
	}
	k8sDependencyMapping = map[string]Dependency{
		"nginx": nginxIngressDependency,
	}
)

// SortDogusByDependency takes an unsorted slice of Dogu structs and returns a slice of Dogus ordered by the
// importance of their dependencies descending, that is: the most needed dogu will be the first element.
//
// Deprecated: This function returns nil on error. For better error-handling use SortDogusByDependencyWithError instead.
//
//goland:noinspection GoDeprecation
func SortDogusByDependency(dogus []*Dogu) []*Dogu {
	orderedDogus, _ := SortDogusByDependencyWithError(dogus)
	return orderedDogus
}

// SortDogusByDependencyWithError takes an unsorted slice of Dogu structs and returns a slice of Dogus ordered by the
// importance of their dependencies descending, that is: the most needed dogu will be the first element.
func SortDogusByDependencyWithError(dogus []*Dogu) ([]*Dogu, error) {
	ordered := sortByDependency{dogus}
	orderedDogus, err := ordered.sortDogusByDependency()
	if err != nil {
		err = fmt.Errorf("error in sorting dogus by dependency: %s", err)
		log.Error(err)
	}
	return orderedDogus, err
}

// SortDogusByInvertedDependency takes an unsorted slice of Dogu structs and returns a new slice of Dogus ordered by the
// importance of their dependencies ascending, that is: the most independent dogu will be the first element.
//
// Deprecated: This function returns nil on error. For better error-handling use SortDogusByInvertedDependencyWithError instead.
//
//goland:noinspection GoDeprecation
func SortDogusByInvertedDependency(dogus []*Dogu) []*Dogu {
	orderedDogus, _ := SortDogusByInvertedDependencyWithError(dogus)
	return orderedDogus
}

// SortDogusByInvertedDependencyWithError takes an unsorted slice of Dogu structs and returns a new slice of Dogus ordered by the
// importance of their dependencies ascending, that is: the most independent dogu will be the first element.
func SortDogusByInvertedDependencyWithError(dogus []*Dogu) ([]*Dogu, error) {
	ordered := sortByDependency{dogus}
	orderedDogus, err := ordered.sortDogusByInvertedDependency()
	if err != nil {
		err = fmt.Errorf("error in sorting dogus by inverted dependency: %s", err)
		log.Error(err)
	}
	return orderedDogus, err
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

func (bd *sortByDependency) sortDogusByDependency() ([]*Dogu, error) {
	dependencyEdges := bd.getDependencyEdges()
	sorted, err := toposort.Toposort(dependencyEdges)
	return bd.handleSortResult(sorted, err)
}

func (bd *sortByDependency) getDependencyEdges() []toposort.Edge {
	var dependencyEdges []toposort.Edge
	for _, dogu := range bd.dogus {
		dependencies := dogu.GetAllDependenciesOfType(DependencyTypeDogu)
		dependentDogus := bd.dependenciesToDogus(dependencies)
		if len(dependentDogus) > 0 {
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

	// we can just append all k8s mappings here, since not installed dogus will be removed in the next step
	dependencies = appendK8sMappedDependencies(dependencies)

	for _, dogu := range bd.dogus {
		if contains(dependencies, dogu.GetSimpleName()) {
			result = append(result, dogu)
		}
	}

	return result
}

func appendK8sMappedDependencies(dependencies []Dependency) []Dependency {
	for _, dep := range dependencies {
		if _, ok := k8sDependencyMapping[dep.Name]; ok {
			dependencies = append(dependencies, k8sDependencyMapping[dep.Name])
		}
	}
	return dependencies
}

func (bd *sortByDependency) sortDogusByInvertedDependency() ([]*Dogu, error) {
	dependencyEdges := bd.getDependencyEdges()
	sorted, err := toposort.ToposortR(dependencyEdges)
	return bd.handleSortResult(sorted, err)
}

func (bd *sortByDependency) handleSortResult(sorted []interface{}, err error) ([]*Dogu, error) {
	if err != nil {
		err = fmt.Errorf("sort by dependency failed: %s", err)
		log.Error(err)
		return nil, err
	}

	sortedDogus, err := toDoguSlice(sorted)
	if err != nil {
		err = fmt.Errorf("sort by dependency failed: %s", err)
		log.Error(err)
		return nil, err
	}

	return sortedDogus, nil
}
