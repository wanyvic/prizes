package refresh

import (
	"container/heap"
	"fmt"
	"testing"
	"time"
)

func Test_service_check(t *testing.T) {
	h := &ServiceLoop{}
	heap.Init(h)

	item1 := CheckItem{ServiceID: "1", RemoveAt: time.Date(2018, 1, 2, 15, 30, 10, 0, time.Local)}
	item2 := CheckItem{ServiceID: "2", RemoveAt: time.Date(2016, 1, 2, 13, 30, 10, 0, time.Local)}
	item3 := CheckItem{ServiceID: "3", RemoveAt: time.Date(2019, 1, 2, 25, 30, 10, 0, time.Local)}
	item4 := CheckItem{ServiceID: "4", RemoveAt: time.Date(2011, 1, 2, 15, 30, 10, 0, time.Local)}
	item5 := CheckItem{ServiceID: "5", RemoveAt: time.Date(2020, 1, 2, 13, 30, 10, 0, time.Local)}
	item6 := CheckItem{ServiceID: "6", RemoveAt: time.Date(2019, 1, 2, 25, 30, 10, 0, time.Local)}

	heap.Push(h, &item1)
	heap.Push(h, &item2)
	heap.Push(h, &item3)
	heap.Push(h, &item4)
	heap.Push(h, &item5)
	heap.Push(h, &item6)
	for _, a := range *h {

		fmt.Println(a)
	}
	fmt.Println()
	for i := 0; i < len(*h); i++ {
		a := *h
		if a[i].ServiceID == "5" {
			a[i].RemoveAt = time.Date(2000, 1, 2, 25, 30, 10, 0, time.Local)
			heap.Fix(h, i)
		}
	}

	// // heap.Fix(h, 3)
	// fmt.Println(*h)
	for _, a := range *h {

		fmt.Println(a)
	}

	fmt.Println()
	for h.Len() > 0 {
		fmt.Println(heap.Pop(h))
	}
}
