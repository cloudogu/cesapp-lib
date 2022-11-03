package dependencies

import (
	"fmt"
	"strings"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"

	"github.com/hashicorp/go-multierror"
)

type doguDependencyChecker struct {
	doguRegistry registry.DoguRegistry
}

var (
	log = core.GetLogger()
)

// NewDoguDependencyChecker creates a new checker for dogu dependencies with a given dogu registry.
func NewDoguDependencyChecker(doguRegistry registry.DoguRegistry) *doguDependencyChecker {
	return &doguDependencyChecker{
		doguRegistry: doguRegistry,
	}
}

// CheckAllDependencies checks mandatory and optional dependencies from a dogu.
func (dc *doguDependencyChecker) CheckAllDependencies(dogu core.Dogu) error {
	var allProblems error

	err := dc.CheckMandatoryDependencies(dogu)
	if err != nil {
		allProblems = multierror.Append(allProblems, err)
	}

	err = dc.CheckOptionalDependencies(dogu)
	if err != nil {
		allProblems = multierror.Append(allProblems, err)
	}

	return allProblems
}

// CheckMandatoryDependencies checks only mandatory dependencies from a dogu.
func (dc *doguDependencyChecker) CheckMandatoryDependencies(dogu core.Dogu) error {
	dependencies := dogu.GetDependenciesOfType(core.DependencyTypeDogu)

	return dc.checkDoguDependencies(dependencies, false)
}

// CheckOptionalDependencies checks only optional dependencies from a dogu.
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

// CheckDoguDependency checks a single dependency from a dogu.
func (dc *doguDependencyChecker) CheckDoguDependency(doguDependency core.Dependency, optional bool) error {
	log.Debugf("checking dogu dependency %s:%s", doguDependency.Name, doguDependency.Version)
	localDependency, err := dc.doguRegistry.Get(doguDependency.Name)
	if err != nil {
		if optional && strings.Contains(err.Error(), "Key not found") { // if a dogu is not found an error is returned
			return nil // not installed => no error as this is ok for optional dependencies
		}
		return fmt.Errorf("failed to resolve dependencies %s: %w", doguDependency.Name, err)
	}
	if localDependency == nil {
		if optional {
			return nil // not installed => no error as this is ok for optional dependencies
		}
		return fmt.Errorf("dependency %s seems not to be installed", doguDependency.Name)
	}
	// it does not count as an error if no version is specified as the field is optional
	if doguDependency.Version != "" {
		localDependencyVersion, err := core.ParseVersion(localDependency.Version)
		if err != nil {
			return fmt.Errorf("failed to parse version of dependency %s: %w", localDependency.Name, err)
		}
		comparator, err := core.ParseVersionComparator(doguDependency.Version)
		if err != nil {
			return fmt.Errorf("failed to parse ParseVersionComparator of version %s for doguDependency %s: %w", doguDependency.Version, doguDependency.Name, err)
		}
		allows, err := comparator.Allows(localDependencyVersion)
		if err != nil {
			return fmt.Errorf("an error occurred when comparing the versions: %w", err)
		}
		if !allows {
			return fmt.Errorf("%s parsed Version does not fulfill version requirement of %s dogu %s", localDependency.Version, doguDependency.Version, doguDependency.Name)
		}
	}
	return nil // no error, dependency is ok
}
