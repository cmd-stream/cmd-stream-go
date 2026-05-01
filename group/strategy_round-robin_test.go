package group_test

import (
	"testing"

	grp "github.com/cmd-stream/cmd-stream-go/group"
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
		s = grp.NewRoundRobinStrategy([]int{5, 10, 15})
	)
	for i := range 9 {
		e, index := s.Next()
		asserterror.Equal(t, e, testCases[i].elem)
		asserterror.Equal(t, index, testCases[i].index)
	}
}
