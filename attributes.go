package oteltest

// attributes in OTel are key-value pairs, where the keys always are strings,
// but the values can be of the following different types:
//  - bool and []bool
//  - int64 and []int64
//  - float64 and []float64
//  - string and []string

// HaveAttribute succeeds if a metric has an attribute with the specified name
// (key) and optional value.
//
// The value passed into the attribute parameter can be either a string or
// [types.GomegaMatcher]:
//   - a string in the form of “name” where it must match an attribute key, or
//     in the “name=value” form where it must match both the attribute key and
//     string value.
//   - a GomegaMatcher that matches the (string) key only.
//   - any other type of value is an error.
//
// To test for attribute values other than string, use [HaveAttributeWithValue]
// instead.
func HaveAttribute(attribute any) {}

// HaveAttributeWithValue succeeds if a metric has an attibute with the
// specified name (key) and value.
func HaveAttributeWithValue(name any, value any) {}
