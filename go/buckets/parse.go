package buckets

import (
	"fmt"
	"strings"
)

// ParseFunc is a function that creates a BucketingStrategy from a map of arguments.
// The arguments are parsed from the specification string.
type ParseFunc func(args map[string]string) (BucketingStrategy, error)

var factories = make(map[string]ParseFunc)

// RegisterParser registers a function that can create a bucketing strategy from arguments in a specification string.
// It is meant to be called from init functions before any calls into Parse.
func RegisterParser(name string, factory ParseFunc) error {
	if factory == nil {
		return fmt.Errorf("nil bucketing strategy factory")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("empty bucketing strategy name")
	}
	name = strings.ToLower(name)
	if _, ok := factories[name]; ok {
		return fmt.Errorf("bucketing strategy %q already registered", name)
	}
	factories[name] = factory
	return nil
}

// Parse parses a specification string and returns a BucketingStrategy.
// The specification string is in the format "name:args", where name is the name of the bucketing strategy
// and args is a comma-separated list of key=value pairs.
// The arguments are parsed into a map[string]string and passed to the registered factory function.
func Parse(spec string) (BucketingStrategy, error) {
	trimmed := strings.TrimSpace(spec)
	if trimmed == "" {
		return nil, fmt.Errorf("empty bucketing strategy")
	}
	name, params, hasParams := strings.Cut(trimmed, ":")
	name = strings.ToLower(strings.TrimSpace(name))
	factory, ok := factories[name]
	if !ok {
		return nil, fmt.Errorf("unknown bucketing strategy %q", name)
	}
	args := make(map[string]string)
	if hasParams {
		for entry := range strings.SplitSeq(params, ",") {
			entry = strings.TrimSpace(entry)
			if entry == "" {
				continue
			}
			key, value, _ := strings.Cut(entry, "=")
			key = strings.ToLower(strings.TrimSpace(key))
			value = strings.TrimSpace(value)
			args[key] = value
		}
	}
	return factory(args)
}
