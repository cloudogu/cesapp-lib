package core

// ByDoguVersion implements sort.Interface for []Dogu to Dogus by their versions
type ByDoguVersion []*Dogu

// Len is the number of elements in the collection.
func (doguVersions ByDoguVersion) Len() int {
	return len(doguVersions)
}

// Swap swaps the elements with indexes i and j.
func (doguVersions ByDoguVersion) Swap(i, j int) {
	doguVersions[i], doguVersions[j] = doguVersions[j], doguVersions[i]
}

// Less reports whether the element with index i should sort before the element with index j.
func (doguVersions ByDoguVersion) Less(i, j int) bool {
	v1, err := ParseVersion(doguVersions[i].Version)
	if err != nil {
		GetLogger().Errorf("connot parse version %s for comparison", doguVersions[i].Version)
	}
	v2, err := ParseVersion(doguVersions[j].Version)
	if err != nil {
		GetLogger().Errorf("connot parse version %s for comparison", doguVersions[j].Version)
	}

	isNewer := v1.IsNewerThan(v2)
	return isNewer
}
