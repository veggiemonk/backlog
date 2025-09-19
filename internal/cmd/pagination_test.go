package cmd

import (
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/viper"
)

func TestApplyDefaultPagination(t *testing.T) {
	// Save original config values
	origPageSize := viper.GetInt(configPageSize)
	origMaxLimit := viper.GetInt(configMaxLimit)
	
	// Restore at the end
	defer func() {
		viper.Set(configPageSize, origPageSize)
		viper.Set(configMaxLimit, origMaxLimit)
	}()
	
	// Set test config values
	viper.Set(configPageSize, 10)
	viper.Set(configMaxLimit, 50)
	
	tests := []struct {
		name           string
		inputLimit     *int
		inputOffset    *int
		expectedLimit  *int
		expectedOffset *int
	}{
		{
			name:           "no_pagination_requested",
			inputLimit:     nil,
			inputOffset:    nil,
			expectedLimit:  nil,
			expectedOffset: nil,
		},
		{
			name:           "explicit_limit_under_max",
			inputLimit:     intPtr(20),
			inputOffset:    nil,
			expectedLimit:  intPtr(20),
			expectedOffset: nil,
		},
		{
			name:           "explicit_limit_over_max",
			inputLimit:     intPtr(100),
			inputOffset:    nil,
			expectedLimit:  intPtr(50), // Should be capped at max
			expectedOffset: nil,
		},
		{
			name:           "explicit_offset_no_limit",
			inputLimit:     nil,
			inputOffset:    intPtr(5),
			expectedLimit:  nil,
			expectedOffset: intPtr(5),
		},
		{
			name:           "explicit_limit_and_offset",
			inputLimit:     intPtr(15),
			inputOffset:    intPtr(5),
			expectedLimit:  intPtr(15),
			expectedOffset: intPtr(5),
		},
		{
			name:           "zero_limit",
			inputLimit:     intPtr(0),
			inputOffset:    nil,
			expectedLimit:  intPtr(0),
			expectedOffset: nil,
		},
		{
			name:           "zero_offset",
			inputLimit:     intPtr(10),
			inputOffset:    intPtr(0),
			expectedLimit:  intPtr(10),
			expectedOffset: intPtr(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			
			resultLimit, resultOffset := ApplyDefaultPagination(tt.inputLimit, tt.inputOffset)
			
			if tt.expectedLimit == nil {
				is.True(resultLimit == nil)
			} else {
				is.True(resultLimit != nil)
				is.Equal(*resultLimit, *tt.expectedLimit)
			}
			
			if tt.expectedOffset == nil {
				is.True(resultOffset == nil)
			} else {
				is.True(resultOffset != nil)
				is.Equal(*resultOffset, *tt.expectedOffset)
			}
		})
	}
}

func TestGetConfigurationValues(t *testing.T) {
	is := is.New(t)
	
	// Save original values
	origPageSize := viper.GetInt(configPageSize)
	origMaxLimit := viper.GetInt(configMaxLimit)
	
	// Restore at the end
	defer func() {
		viper.Set(configPageSize, origPageSize)
		viper.Set(configMaxLimit, origMaxLimit)
	}()
	
	// Test setting custom values
	viper.Set(configPageSize, 30)
	viper.Set(configMaxLimit, 200)
	
	is.Equal(GetDefaultPageSize(), 30)
	is.Equal(GetMaxLimit(), 200)
	
	// Test with different values
	viper.Set(configPageSize, 5)
	viper.Set(configMaxLimit, 100)
	
	is.Equal(GetDefaultPageSize(), 5)
	is.Equal(GetMaxLimit(), 100)
}

func intPtr(i int) *int {
	return &i
}
