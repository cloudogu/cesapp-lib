package remote

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

const (
	// yamlIndent consists of two whitespaces.
	yamlIndent = "  "
)

//DoguVersionPrinter provides methods for printing available dogu versions.
type DoguVersionPrinter struct {
	//Remote contains a remote dogu registry reference thath provides dogus and their versions.
	Remote Registry
}

//PrintForAllDogus prints versions for all available dogus (up to the given limit). The limit must be zero (0) or a
//positive integer whereas zero (0) will print all available versions. Otherwise a number of version up to the given
//limit can be printed.
func (dvp *DoguVersionPrinter) PrintForAllDogus(limit int) error {
	if limit < 0 {
		return errors.Errorf("invalid dogu limit '%d': limit must be zero or positive", limit)
	}

	dogus, err := dvp.Remote.GetAll()
	if err != nil {
		return errors.Wrap(err, "failed to fetch dogus from remote")
	}

	var result error

	for _, dogu := range dogus {
		err := dvp.PrintForSingleDogu(dogu, limit)
		if err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}

//PrintForSingleDogu prints versions for the given dogu (up to the given limit). The limit must be zero (0) or a
//positive integer whereas zero (0) will print all available versions. Otherwise a number of version up to the given
//limit can be printed.
func (dvp *DoguVersionPrinter) PrintForSingleDogu(dogu *core.Dogu, limit int) error {
	if limit < 0 {
		return errors.Errorf("invalid dogu limit '%d': limit must be zero or positive", limit)
	}

	versions, err := dvp.Remote.GetVersionsOf(dogu.Name)
	if err != nil {
		return errors.Wrapf(err, "could not get versions for dogu '%s' in registry", dogu.Name)
	}

	if limit == 0 {
		printDoguWithVersions(dogu, versions)
	} else {
		printDoguWithVersions(dogu, limitVersions(versions, limit))
	}

	return nil
}

//PrintDoguListInDefaultFormat prints all dogus and their latest version in tabular format.
func (dvp *DoguVersionPrinter) PrintDoguListInDefaultFormat() error {
	dogus, err := dvp.Remote.GetAll()
	if err != nil {
		return errors.Wrapf(err, "could not get all dogus in remote registry")
	}
	core.PrintDogus(dogus, true)

	return nil
}

func printDoguWithVersions(dogu *core.Dogu, versions []core.Version) {
	core.GetLogger().Printf("%s%s:\n", yamlIndent, dogu.Name)

	for _, version := range versions {
		core.GetLogger().Printf("%s%s - %s\n", yamlIndent, yamlIndent, version.Raw)
	}
}

func limitVersions(versions []core.Version, limit int) []core.Version {
	var result []core.Version

	if len(versions) <= limit {
		return versions
	}

	for i := 0; i < limit; i++ {
		result = append(result, versions[i])
	}

	return result
}
