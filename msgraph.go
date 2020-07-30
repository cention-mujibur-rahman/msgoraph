// Package msgoraph implements a Go interface for the Microsoft Graph API
package msgoraph

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// BoolValue returns the value of the bool pointer passed in or
// false if the pointer is nil.
func BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}

// BoolSlice returns a slice of bool pointers given a slice of bools.
func BoolSlice(v []bool) []*bool {
	out := make([]*bool, len(v))
	for i := range v {
		out[i] = &v[i]
	}
	return out
}

// Float64 returns a pointer to the float64 value passed in.
func Float64(v float64) *float64 {
	return &v
}

// Float64Value returns the value of the float64 pointer passed in or
// 0 if the pointer is nil.
func Float64Value(v *float64) float64 {
	if v != nil {
		return *v
	}
	return 0
}

// Float64Slice returns a slice of float64 pointers given a slice of float64s.
func Float64Slice(v []float64) []*float64 {
	out := make([]*float64, len(v))
	for i := range v {
		out[i] = &v[i]
	}
	return out
}

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// IntValue returns the value of the int pointer passed in or
// 0 if the pointer is nil.
func IntValue(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}

// IntSlice returns a slice of int pointers given a slice of int's.
func IntSlice(v []int) []*int {
	out := make([]*int, len(v))
	for i := range v {
		out[i] = &v[i]
	}
	return out
}

// Int64 is a helper routine that allocates a new int64 value
// to store v and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// Int64Value returns the value of the int64 pointer passed in or
// 0 if the pointer is nil.
func Int64Value(v *int64) int64 {
	if v != nil {
		return *v
	}
	return 0
}

// Int64Slice returns a slice of int64 pointers given a slice of int64s.
func Int64Slice(v []int64) []*int64 {
	out := make([]*int64, len(v))
	for i := range v {
		out[i] = &v[i]
	}
	return out
}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }

// StringValue returns the value of the string pointer passed in or
// "" if the pointer is nil.
func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

// StringSlice returns a slice of string pointers given a slice of strings.
func StringSlice(v []string) []*string {
	out := make([]*string, len(v))
	for i := range v {
		out[i] = &v[i]
	}
	return out
}
