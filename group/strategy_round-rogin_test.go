package cgrp

import (
	"testing"

	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestRoundRobinStrategy(t *testing.T) {
	var (
		testCases = []struct {
			elem  int
			index int64
		}{
			{elem: 5, index: 0},
			{elem: 10, index: 1},
			{elem: 15, index: 2},
			{elem: 5, index: 0},
			{elem: 10, index: 1},
			{elem: 15, index: 2},
			{elem: 5, index: 0},
			{elem: 10, index: 1},
			{elem: 15, index: 2},
		}
		s = NewRoundRobinStrategy([]int{5, 10, 15})
	)
	for i := range 9 {
		e, index := s.Next()
		asserterror.Equal(e, testCases[i].elem, t)
		asserterror.Equal(index, testCases[i].index, t)
	}
}
