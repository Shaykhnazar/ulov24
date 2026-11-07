package rides

import (
	"slices"
	"testing"
)

func TestPaginate(t *testing.T) {
	var items []int
	for i := 1; i <= 23; i++ {
		items = append(items, i)
	}
	p := Paginate(items, 20, 10)
	if !slices.Equal(p.Items, []int{21, 22, 23}) || p.Total != 23 {
		t.Fatalf("offset 20: %+v", p)
	}
	if p := Paginate(items, 30, 10); len(p.Items) != 0 || p.Total != 23 {
		t.Fatalf("offset 30: %+v", p)
	}
	if p := Paginate(items, 0, 0); len(p.Items) != 0 {
		t.Fatalf("limit 0: %+v", p)
	}
	if p := Paginate(items, -5, 2); !slices.Equal(p.Items, []int{1, 2}) {
		t.Fatalf("negative offset: %+v", p)
	}
}
