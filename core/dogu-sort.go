package core

import (
	"sort"
)

// SortDogusByDependency takes an unsorted slice of Dogu structs and returns a slice of Dogus ordered by the
// importance of their dependencies descending, that is: the most needed dogu will be the first element.
func SortDogusByDependency(dogus []*Dogu) []*Dogu {
	ordered := sortByDependency{dogus}
	sort.Stable(&ordered)
	return ordered.dogus
}

// SortDogusByInvertedDependency takes an unsorted slice of Dogu structs and returns a new slice of Dogus ordered by the
// importance of their dependencies ascending, that is: the most independent dogu will be the first element.
func SortDogusByInvertedDependency(dogus []*Dogu) []*Dogu {
	orderedDesc := SortDogusByDependency(dogus)

	orderedAsc := []*Dogu{}
	for i := len(orderedDesc) - 1; i >= 0; i-- {
		orderedAsc = append(orderedAsc, dogus[i])
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

func (bd *sortByDependency) Len() int {
	return len(bd.dogus)
}

func (bd *sortByDependency) Swap(i, j int) {
	bd.dogus[i], bd.dogus[j] = bd.dogus[j], bd.dogus[i]
}

func (bd *sortByDependency) Less(i, j int) bool {
	leftName := bd.dogus[i].GetSimpleName()
	rightName := bd.dogus[j].GetSimpleName()

	leftDependenciesRecursive := bd.getAllDoguDependenciesRecursive(bd.dogus[i])
	rightDependenciesRecursive := bd.getAllDoguDependenciesRecursive(bd.dogus[j])

	leftDependenciesDirect := bd.dogus[i].GetAllDependenciesOfType(DependencyTypeDogu)
	rightDependenciesDirect := bd.dogus[j].GetAllDependenciesOfType(DependencyTypeDogu)

	if contains(rightDependenciesRecursive, leftName) {
		return true
	} else if contains(leftDependenciesRecursive, rightName) {
		return false
	} else if len(leftDependenciesDirect) == len(rightDependenciesDirect) {
		return leftName < rightName
	}

	return len(leftDependenciesDirect) < len(rightDependenciesDirect)
}

func (bd *sortByDependency) getAllDoguDependenciesRecursive(inputDogu *Dogu) []Dependency {
	dependencies := inputDogu.GetAllDependenciesOfType(DependencyTypeDogu)
	dependenciesAsDogus := bd.dependenciesToDogus(dependencies)

	for _, dogu := range dependenciesAsDogus {
		for _, dep := range bd.getAllDoguDependenciesRecursive(dogu) {
			if !contains(dependencies, dep.Name) {
				dependencies = append(dependencies, dep)
			}
		}
	}

	return dependencies
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
