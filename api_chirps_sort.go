package main

import (
	"cmp"
	"slices"

	"github.com/jabuta/chirpy/internal/database"
)

func ascSort(chirps []database.Chirp) {
	slices.SortStableFunc(chirps, func(i, j database.Chirp) int {
		return cmp.Compare(i.ID, j.ID)
	})
}
func dscSort(chirps []database.Chirp) {
	slices.SortStableFunc(chirps, func(i, j database.Chirp) int {
		return cmp.Compare(j.ID, i.ID)
	})
}
