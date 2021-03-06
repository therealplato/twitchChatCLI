package main

import "sort"

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}

		if len(vsf) > 10 {
			break
		}
	}
	return vsf
}

func Contains(slice []string, thing string) bool {
	sort.Strings(slice)
	i := sort.SearchStrings(slice, thing)
	return i < len(slice) && slice[i] == thing
}