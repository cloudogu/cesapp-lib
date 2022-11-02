package dependencies

import (
	"testing"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"

	"github.com/stretchr/testify/assert"
)

func TestDoguDependencyChecker_CheckDependencies(t *testing.T) {
	var problems error

	mockRegistry := registry.MockRegistry{}
	doguRegistry := mockRegistry.DoguRegistry()

	d1 := core.Dogu{Name: "a valid dogu", Version: "1.2.3-1", Dependencies: []core.Dependency{{Name: "dogu-b", Version: "2.2.2"}}}
	d2 := core.Dogu{Name: "a dogu for which doguRegistry.Get() will find nothing", Version: "1.2.3-1", Category: "nil"}
	d3 := core.Dogu{Name: "a dogu for which version compare will fail", Version: "1.2.3-1"}
	d4 := core.Dogu{Name: "a dogu for which allows returns false", Version: "0.2.3-1"}
	d5 := core.Dogu{Name: "a dogu for which ParseVersion will fail", Version: "0-jeff-"}
	d6 := core.Dogu{Name: "a dogu for which ParseVersionComparator will fail", Version: "0.2.3-1"}

	_ = doguRegistry.Register(&d1)
	_ = doguRegistry.Register(&d2)
	_ = doguRegistry.Register(&d3)
	_ = doguRegistry.Register(&d4)
	_ = doguRegistry.Register(&d5)
	_ = doguRegistry.Register(&d6)

	clientDependencyChecker := NewDoguDependencyChecker(doguRegistry)

	problems = clientDependencyChecker.CheckAllDependencies(core.Dogu{Name: "dogu to check", Version: "2.2.5",
		Dependencies: []core.Dependency{
			{Type: core.DependencyTypeDogu, Name: "a dogu for which doguRegistry.Get() will fail", Version: "1.0.0-1"},
			{Type: core.DependencyTypeDogu, Name: "a dogu for which doguRegistry.Get() will find nothing", Version: "1.0.0-1"},
			{Type: core.DependencyTypeDogu, Name: "a dogu for which version compare will fail", Version: "<>1.0.0-1"},
			{Type: core.DependencyTypeDogu, Name: "a dogu for which allows returns false", Version: ">1.0.0-1"},
			{Type: core.DependencyTypeDogu, Name: "a valid dogu", Version: "<=2.3.1-7"},
			{Type: core.DependencyTypeDogu, Name: "a dogu for which ParseVersion will fail", Version: "x!2.3.1-7"},
			{Type: core.DependencyTypeDogu, Name: "a dogu for which ParseVersionComparator will fail", Version: "x!2.3.1-7"},
		},
		OptionalDependencies: []core.Dependency{
			{Type: core.DependencyTypeDogu, Name: "a dogu for which doguRegistry.Get() will fail", Version: "1.0.0-1"},
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
	assert.Contains(t, problems.Error(), "An error occurred when comparing the versions")
	assert.Contains(t, problems.Error(), "0.2.3-1 parsed Version does not fulfill version requirement of >1.0.0-1")
	assert.Contains(t, problems.Error(), "failed to parse version of dependency a dogu for which ParseVersion will fail")
	assert.Contains(t, problems.Error(), "failed to parse ParseVersionComparator of version x!2.3.1-7")
	assert.NotContains(t, problems.Error(), "optionaldoguwhichisnotninstalled")
}
