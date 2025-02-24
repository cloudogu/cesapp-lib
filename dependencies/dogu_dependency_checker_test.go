package dependencies

import (
	"errors"
	"github.com/cloudogu/cesapp-lib/registry/mocks"
	"testing"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
)

func TestDoguDependencyChecker_CheckDependencies(t *testing.T) {
	var problems error

	doguRegistry := mocks.NewDoguRegistry(t)

	d1 := &core.Dogu{Name: "a valid dogu", Version: "1.2.3-1", Dependencies: []core.Dependency{{Name: "dogu-b", Version: "2.2.2"}}}
	d3 := &core.Dogu{Name: "a dogu for which version compare will fail", Version: "1.2.3-1"}
	d4 := &core.Dogu{Name: "a dogu for which allows returns false", Version: "0.2.3-1"}
	d5 := &core.Dogu{Name: "a dogu for which ParseVersion will fail", Version: "0-jeff-"}
	d6 := &core.Dogu{Name: "a dogu for which ParseVersionComparator will fail", Version: "0.2.3-1"}

	doguRegistry.On("Get", "a dogu for which doguRegistry.Get() will fail").Return(nil, assert.AnError)
	doguRegistry.On("Get", "a dogu for which doguRegistry.Get() will return nil nil optional").Return(nil, nil)
	doguRegistry.On("Get", "a dogu for which doguRegistry.Get() will fail with keyNotFound error").Return(nil, errors.New("error: Key not found"))
	doguRegistry.On("Get", "a valid dogu").Return(d1, nil)
	doguRegistry.On("Get", "a dogu for which doguRegistry.Get() will find nothing").Return(nil, nil)
	doguRegistry.On("Get", "a dogu for which version compare will fail").Return(d3, nil)
	doguRegistry.On("Get", "a dogu for which allows returns false").Return(d4, nil)
	doguRegistry.On("Get", "a dogu for which ParseVersion will fail").Return(d5, nil)
	doguRegistry.On("Get", "a dogu for which ParseVersionComparator will fail").Return(d6, nil)

	clientDependencyChecker := NewDoguDependencyChecker(doguRegistry)

	problems = clientDependencyChecker.CheckAllDependencies(core.Dogu{Name: "dogu to check", Version: "2.2.5",
		Dependencies: []core.Dependency{
			{Type: core.DependencyTypeDogu, Name: "a dogu for which doguRegistry.Get() will find nothing", Version: "1.2.3-1"},
			{Type: core.DependencyTypeDogu, Name: "a valid dogu", Version: "<=2.3.1-7"},
		},
		OptionalDependencies: []core.Dependency{
			{Type: core.DependencyTypeDogu, Name: "a dogu for which doguRegistry.Get() will fail", Version: "1.0.0-1"},
			{Type: core.DependencyTypeDogu, Name: "a dogu for which doguRegistry.Get() will return nil nil optional", Version: "1.0.0-1"},
			{Type: core.DependencyTypeDogu, Name: "a dogu for which doguRegistry.Get() will fail with keyNotFound error", Version: "1.0.0-1"},
			{Type: core.DependencyTypeDogu, Name: "a dogu for which doguRegistry.Get() will find nothing", Version: "1.0.0-1"},
			{Type: core.DependencyTypeDogu, Name: "a dogu for which version compare will fail", Version: "<>1.0.0-1"},
			{Type: core.DependencyTypeDogu, Name: "a dogu for which allows returns false", Version: ">1.0.0-1"},
			{Type: core.DependencyTypeDogu, Name: "a valid dogu", Version: "<=2.3.1-7"},
			{Type: core.DependencyTypeDogu, Name: "a dogu for which ParseVersion will fail", Version: "x!2.3.1-7"},
			{Type: core.DependencyTypeDogu, Name: "a dogu for which ParseVersionComparator will fail", Version: "x!2.3.1-7"},
		},
	})

	assert.NotNil(t, problems)
	assert.NotContains(t, problems.Error(), "a valid dogu")
	assert.Contains(t, problems.Error(), "failed to resolve dependencies a dogu for which doguRegistry.Get() will fail")
	assert.Contains(t, problems.Error(), "an error occurred when comparing the versions")
	assert.Contains(t, problems.Error(), "The parsed version of the dogu a dogu for which allows returns false (0.2.3-1) does not fulfill the version requirement of the dogu dependency a dogu for which allows returns false (>1.0.0-1)")
	assert.Contains(t, problems.Error(), "failed to parse version of dependency a dogu for which ParseVersion will fail")
	assert.Contains(t, problems.Error(), "failed to parse ParseVersionComparator of version x!2.3.1-7")
	assert.Contains(t, problems.Error(), "dependency a dogu for which doguRegistry.Get() will find nothing seems not to be installed")
	assert.NotContains(t, problems.Error(), "optionaldoguwhichisnotninstalled")
}
