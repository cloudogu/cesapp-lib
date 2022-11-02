package dependencies

import (
	"strings"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"

	"github.com/hashicorp/go-multierror"

	"github.com/pkg/errors"
)

type doguDependencyChecker struct {
	doguRegistry registry.DoguRegistry
}

var (
	log = core.GetLogger()
)

func NewDoguDependencyChecker(doguRegistry registry.DoguRegistry) *doguDependencyChecker {
	return &doguDependencyChecker{
		doguRegistry: doguRegistry,
	}
}

func (dc *doguDependencyChecker) CheckAllDependencies(dogu core.Dogu) error {
	var allProblems error

	err := dc.CheckMandatoryDependencies(dogu)
	if err != nil {
		allProblems = multierror.Append(err)
	}

	err = dc.CheckOptionalDependencies(dogu)
	if err != nil {
		allProblems = multierror.Append(err)
	}

	return allProblems
}

func (dc *doguDependencyChecker) CheckMandatoryDependencies(dogu core.Dogu) error {
	dependencies := dogu.GetDependenciesOfType(core.DependencyTypeDogu)

	return dc.checkDoguDependencies(dependencies, false)
}

func (dc *doguDependencyChecker) CheckOptionalDependencies(dogu core.Dogu) error {
	dependencies := dogu.GetOptionalDependenciesOfType(core.DependencyTypeDogu)

	return dc.checkDoguDependencies(dependencies, true)
}

func (dc *doguDependencyChecker) checkDoguDependencies(dependencies []core.Dependency, optional bool) error {
	var problems error

	for _, doguDependency := range dependencies {
		err := dc.CheckDoguDependency(doguDependency, optional)
		if err != nil {
			problems = multierror.Append(problems, err)
		}
	}
	return problems
}

func (dc *doguDependencyChecker) CheckDoguDependency(doguDependency core.Dependency, optional bool) error {
	log.Debugf("checking dogu dependency %s:%s", doguDependency.Name, doguDependency.Version)
	localDependency, err := dc.doguRegistry.Get(doguDependency.Name)
	if err != nil {
		if optional && strings.Contains(err.Error(), "Key not found") { // if a dogu is not found an error is returned
			return nil // not installed => no error as this is ok for optional dependencies
		}
		return errors.Wrapf(err, "failed to resolve dependencies %s", doguDependency.Name)
	}
	if localDependency == nil {
		if optional {
			return nil // not installed => no error as this is ok for optional dependencies
		}
		return errors.Errorf("dependency %s seems not to be installed", doguDependency.Name)
	}
	// it does not count as a error if no version is specified as the field is optional
	if doguDependency.Version != "" {
		localDependencyVersion, err := core.ParseVersion(localDependency.Version)
		if err != nil {
			return errors.Wrapf(err, "failed to parse version of dependency %s", localDependency.Name)
		}
		comparator, err := core.ParseVersionComparator(doguDependency.Version)
		if err != nil {
			return errors.Wrapf(err, "failed to parse ParseVersionComparator of version %s for doguDependency %s", doguDependency.Version, doguDependency.Name)
		}
		allows, err := comparator.Allows(localDependencyVersion)
		if err != nil {
			return errors.Wrapf(err, "An error occurred when comparing the versions")
		}
		if !allows {
			return errors.Errorf("%s parsed Version does not fulfill version requirement of %s dogu %s", localDependency.Version, doguDependency.Version, doguDependency.Name)
		}
	}
	return nil // no error, dependency is ok
}
