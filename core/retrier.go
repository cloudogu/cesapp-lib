package core

import (
	"fmt"
	"github.com/eapache/go-resiliency/retrier"
	"time"
)

const (
	constantPolicyType    = "constant"
	exponentialPolicyType = "exponential"
)

var log = GetLogger()

func GetBackoff(policy RetryPolicy) ([]time.Duration, error) {
	if policy.Interval < 0 {
		return nil, fmt.Errorf("the retry interval needs to be greater or equal to 0: given '%d'", policy.Interval)
	}
	if policy.MaxRetryCount < 0 {
		return nil, fmt.Errorf("the mxa retry count needs to be greater or equal to 0: given '%d'", policy.MaxRetryCount)
	}
	switch policy.Type {
	case constantPolicyType:
		return retrier.ConstantBackoff(policy.MaxRetryCount, time.Duration(policy.Interval)*time.Millisecond), nil
	case exponentialPolicyType:
		return retrier.ExponentialBackoff(policy.MaxRetryCount, time.Duration(policy.Interval)*time.Millisecond), nil
	default:
		log.Debug("No retry policy type configured. Using type 'exponential' as default")
		return retrier.ExponentialBackoff(policy.MaxRetryCount, time.Duration(policy.Interval)*time.Millisecond), nil
	}
}
