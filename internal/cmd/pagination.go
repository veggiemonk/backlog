package cmd

import "github.com/spf13/viper"

// GetDefaultPageSize returns the configured default page size
func GetDefaultPageSize() int {
	return viper.GetInt(configPageSize)
}

// GetMaxLimit returns the configured maximum limit for pagination
func GetMaxLimit() int {
	return viper.GetInt(configMaxLimit)
}

// ApplyDefaultPagination applies default pagination values if not set
func ApplyDefaultPagination(limit, offset *int) (*int, *int) {
	// Don't apply defaults if user explicitly set offset (means they want pagination)
	if offset != nil && *offset > 0 {
		return limit, offset
	}

	// Don't apply defaults if user explicitly set limit
	if limit != nil && *limit > 0 {
		// Enforce max limit
		maxLimit := GetMaxLimit()
		if *limit > maxLimit {
			enforcedLimit := maxLimit
			return &enforcedLimit, offset
		}
		return limit, offset
	}

	// No pagination requested by user
	return limit, offset
}
