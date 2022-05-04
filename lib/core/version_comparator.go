package core

import (
	"regexp"

	"github.com/pkg/errors"
)

// go does not support lookaheads or lookbehinds (tried "^[=,>,<]+[\"d\"-,.]*" "^[=,>,<]+(?=[d])+"
const OperatorRegex = `^[=><]+`

// =, <, >, <=, >=
type Operator string

type VersionComparator struct {
	version  Version
	operator Operator
}

func ParseVersionComparator(raw string) (VersionComparator, error) {
	var version Version

	operator, err := parseOperator(raw)
	if err != nil {
		return VersionComparator{}, errors.Wrap(err, "failed to parse operator")
	}

	if raw == "" {
		return VersionComparator{}, nil
	}

	version, err = ParseVersion(raw[len(operator):])
	if err != nil {
		return VersionComparator{}, errors.Wrap(err, "failed to parse version")
	}
	return VersionComparator{
		version:  version,
		operator: operator,
	}, nil
}

func (v VersionComparator) Allows(version Version) (bool, error) {
	switch v.operator {
	case "=", "==": //Approximately equivalent to version
		return v.version.IsEqualTo(version), nil
	case ">":
		return v.version.IsOlderThan(version), nil
	case "<":
		return v.version.IsNewerThan(version), nil
	case ">=":
		return v.version.IsOlderOrEqualThan(version), nil
	case "<=":
		return v.version.IsNewerOrEqualThan(version), nil
	case "":
		// this is the edge-case that the dogu.json has no version field or an empty version. At this point we allow
		// every version. if a list of dogus is used it should be best practice to sort the list by version
		if v.version.Raw == "" {
			return true, nil
		}
		return v.version.IsEqualTo(version), nil
	default:
		/* We could use the default operator here but this is a conscious choice to not do this,
		as there are places in the calling code where the error is caught and assessed accordingly */
		return false, errors.Errorf("could not find suitable comperator for '%s' operator", v.operator)
	}
}

func parseOperator(raw string) (Operator, error) {
	var operator Operator

	r, _ := regexp.Compile(OperatorRegex)
	if r.MatchString(raw) {
		idx := r.FindStringIndex(raw)
		operator = Operator(raw[idx[0]:idx[1]]) //cut operator
		if len(operator) > 2 {
			err := errors.Errorf("dependency operator %s of Version %s cannot contain more than two characters. Allowed operators are =,==,>,<,>= and <=", operator, raw)
			GetLogger().Error(err)
			return "", err
		}
	}
	if len(operator) == 0 {
		GetLogger().Debug(errors.Errorf("no dependency operator in %s could be found", raw))
		return "", nil
	}
	return operator, nil
}
