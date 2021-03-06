// deepcopy deep copies maps, slices, etc. A standard copy will copy the
// pointers: deep copy copies the values pointed to.
// 
// Only what is needed has been implemented. Could make more dynamic, at the 
// cost of reflection. Either adjust as needed or create a new function.
// 
// Copyright (c)2014, Joel Scoble (github.com/mohae), all rights reserved.
// License: MIT, for more details check the included LICENSE.txt.
package deepcopy

import (
	"reflect"
)

// InterfaceToSliceStrings takes an interface that is a slice of strings
// and returns a deep copy of it as a slice of strings.
func InterfaceToSliceStrings(v interface{}) []string {
	if v == nil {
		return nil
	}
	var sl []string

	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(v)
		sLen := s.Len()

		for i := 0; i < sLen; i++ {
			sl = append(sl, s.Index(i).Interface().(string))
		}

	default:
		return nil
	}

	return sl
}

// SliceStrings deep copies a slice of strings
func SliceStrings(s []string) []string{
	if s == nil {
		return nil
	}
	
	var sl []string

	sLen := len(s)

	for i := 0; i < sLen; i++ {
		sl = append(sl, s[i])
	}

	return sl
}

// InterfaceToSliceInts takes an interface that is a slice of ints and returns 
// a deep copy of it as a slice of strings. An error is returned if the 
// interface is not a slice of strings.
func InterfaceToSliceInts(v interface{}) []int {
	if v == nil {
		return nil
	}
	var sl []int

	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(v)
		sLen := s.Len()

		for i := 0; i < sLen; i++ {
			sl = append(sl, s.Index(i).Interface().(int))
		}

	default:
		return nil
	}
	return sl
}

// SliceInts deep copies a slice of ints.
func SliceInts(s []int) []int{
	if s == nil {
		return nil
	}
	
	var sl []int

	sLen := len(s)

	for i := 0; i < sLen; i++ {
		sl = append(sl, s[i])
	}

	return sl
}


// Iface recursively deep copies an interface{}
func Iface(iface interface{}) interface{} {
	if iface == nil {
		return nil
	}

	// Make the interface a reflect.Value
	original := reflect.ValueOf(iface)

	// Make a copy of the same type as the original.
	copy := reflect.New(original.Type()).Elem()

	// Recursively copy the original.
	copyRecursive(original, copy)

	// Return theb copy as an interface.
	return copy.Interface()
}

// copyRecursive does the actual copying of the interface. It currently has 
// limited support for what it can handle. Add as needed.
func copyRecursive(original, copy reflect.Value) {
	// handle according to original's Kind
	switch original.Kind() {
	case reflect.Ptr:
		// Get the actual value being pointed to.
		originalValue := original.Elem()

		// if  it isn't valid, return.
		if !originalValue.IsValid() {
			return
		}

		copy.Set(reflect.New(originalValue.Type()))
		copyRecursive(originalValue, copy.Elem())

	case reflect.Interface:
		// Get the value for the interface, not the pointer.
		originalValue := original.Elem()

		// Get the value by calling Elem().
		copyValue := reflect.New(originalValue.Type()).Elem()
		copyRecursive(originalValue, copyValue)
		copy.Set(copyValue)

	case reflect.Struct:
		// Go through each field of the struct and copy it.
		for i := 0; i < original.NumField(); i++ {
			copyRecursive(original.Field(i), copy.Field(i))
		}

	case reflect.Slice:
		// Make a new slice and copy each element.
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i:= 0; i < original.Len(); i++ {
			copyRecursive(original.Index(i), copy.Index(i))
		}

	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			copyRecursive(originalValue, copyValue)
			copy.SetMapIndex(key, copyValue)
		}

	// Set the actual values from here on.
	case reflect.String:
		copy.SetString(original.Interface().(string))
		
	case reflect.Int:
		copy.SetInt(int64(original.Interface().(int)))

	case reflect.Bool:
		copy.SetBool(original.Interface().(bool))

	case reflect.Float64:
		copy.SetFloat(original.Interface().(float64))

	default:
		copy.Set(original)
	}
}
