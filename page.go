package main

import (
	"github.com/halberdholder/bbs/data"
)

type PageInfo struct {
	TotalThreads 	int64
	CurrentPage 	int64
	Min				int64
	Max 			int64
	Left            int64
	Right			int64
	TotalPages		int64
	PageList		[]int64
	Threads 		[]data.Thread
}

func (pi *PageInfo) Pagination() {
	pi.TotalPages = (pi.TotalThreads+(config.PageSize-1)) / config.PageSize
	if pi.TotalPages <= 1 {
		pi.Left = 1
		pi.Right = 1
		pi.Min = 1
		pi.Max = 1
		pi.TotalPages = 1
		return
	}
	PageListCount := pi.TotalPages
	if PageListCount > config.MaxPageList {
		PageListCount = config.MaxPageList
	}

	if pi.Left = pi.CurrentPage - 1; pi.Left < 1 {
		pi.Left = 1
	}
	if pi.Right = pi.CurrentPage + 1; pi.Right > pi.TotalPages {
		pi.Right = pi.TotalPages
	}
	pi.Min = 1
	pi.Max = pi.TotalPages

	pi.PageList = make([]int64, pi.TotalPages + 1)
	for i := int64(1); i <= pi.TotalPages; i++ {
		pi.PageList[i] = i
	}

	for min := int64(1); min <= pi.CurrentPage; min++ {
		max := min + PageListCount - 1
		if max < pi.CurrentPage {
			continue
		}
		if (max - pi.CurrentPage >= pi.CurrentPage - min) ||  (max == pi.TotalPages ) {
			pi.PageList = pi.PageList[min:max+1]
			break
		}
	}
}
