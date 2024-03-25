package middleware

import (
	"reflect"
	"testing"
)

func Test_ContentTypeChecker(t *testing.T) {
	type args struct {
		expectedContentType string
	}
	tests := []struct {
		name string
		args args
		want Middleware
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContentTypeChecker(tt.args.expectedContentType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("сontentTypeCheckerMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
