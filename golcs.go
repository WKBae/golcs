// package lcs provides functions to calculate Longest Common Subsequence (LCS)
// values from two arbitrary arrays.
package lcs

import (
	"context"
	"reflect"
)

// Lcs is the interface to calculate the LCS of two arrays.
type Lcs[Slice ~[]E, E any] interface {
	// Values calculates the LCS value of the two arrays.
	Values() (values Slice)
	// ValueContext is a context aware version of Values()
	ValuesContext(ctx context.Context) (values Slice, err error)
	// IndexPairs calculates paris of indices which have the same value in LCS.
	IndexPairs() (pairs []IndexPair)
	// IndexPairsContext is a context aware version of IndexPairs()
	IndexPairsContext(ctx context.Context) (pairs []IndexPair, err error)
	// Length calculates the length of the LCS.
	Length() (length int)
	// LengthContext is a context aware version of Length()
	LengthContext(ctx context.Context) (length int, err error)
	// Left returns one of the two arrays to be compared.
	Left() (leftValues Slice)
	// Right returns the other of the two arrays to be compared.
	Right() (righttValues Slice)
}

// IndexPair represents an pair of indeices in the Left and Right arrays found in the LCS value.
type IndexPair struct {
	Left  int
	Right int
}

type lcs[Slice ~[]E, E any] struct {
	left  Slice
	right Slice
	/* for caching */
	table      [][]int
	indexPairs []IndexPair
	values     Slice
}

// New creates a new LCS calculator from two arrays.
func New[Slice ~[]E, E any](left, right Slice) Lcs[Slice, E] {
	return &lcs[Slice, E]{
		left:       left,
		right:      right,
		table:      nil,
		indexPairs: nil,
		values:     nil,
	}
}

// Table implements Lcs.Table()
func (lcs *lcs[Slice, E]) Table() (table [][]int) {
	table, _ = lcs.TableContext(context.Background())
	return table
}

// Table implements Lcs.TableContext()
func (lcs *lcs[Slice, E]) TableContext(ctx context.Context) (table [][]int, err error) {
	if lcs.table != nil {
		return lcs.table, nil
	}

	sizeX := len(lcs.left) + 1
	sizeY := len(lcs.right) + 1

	table = make([][]int, sizeX)
	for x := 0; x < sizeX; x++ {
		table[x] = make([]int, sizeY)
	}

	for y := 1; y < sizeY; y++ {
		select { // check in each y to save some time
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// nop
		}
		for x := 1; x < sizeX; x++ {
			increment := 0
			if reflect.DeepEqual(lcs.left[x-1], lcs.right[y-1]) {
				increment = 1
			}
			table[x][y] = max(table[x-1][y-1]+increment, table[x-1][y], table[x][y-1])
		}
	}

	lcs.table = table
	return table, nil
}

// Table implements Lcs.Length()
func (lcs *lcs[Slice, E]) Length() (length int) {
	length, _ = lcs.LengthContext(context.Background())
	return length
}

// Table implements Lcs.LengthContext()
func (lcs *lcs[Slice, E]) LengthContext(ctx context.Context) (length int, err error) {
	table, err := lcs.TableContext(ctx)
	if err != nil {
		return 0, err
	}
	return table[len(lcs.left)][len(lcs.right)], nil
}

// Table implements Lcs.IndexPairs()
func (lcs *lcs[Slice, E]) IndexPairs() (pairs []IndexPair) {
	pairs, _ = lcs.IndexPairsContext(context.Background())
	return pairs
}

// Table implements Lcs.IndexPairsContext()
func (lcs *lcs[Slice, E]) IndexPairsContext(ctx context.Context) (pairs []IndexPair, err error) {
	if lcs.indexPairs != nil {
		return lcs.indexPairs, nil
	}

	table, err := lcs.TableContext(ctx)
	if err != nil {
		return nil, err
	}

	pairs = make([]IndexPair, table[len(table)-1][len(table[0])-1])

	for x, y := len(lcs.left), len(lcs.right); x > 0 && y > 0; {
		if reflect.DeepEqual(lcs.left[x-1], lcs.right[y-1]) {
			pairs[table[x][y]-1] = IndexPair{Left: x - 1, Right: y - 1}
			x--
			y--
		} else {
			if table[x-1][y] >= table[x][y-1] {
				x--
			} else {
				y--
			}
		}
	}

	lcs.indexPairs = pairs

	return pairs, nil
}

// Table implements Lcs.Values()
func (lcs *lcs[Slice, E]) Values() (values Slice) {
	values, _ = lcs.ValuesContext(context.Background())
	return values
}

// Table implements Lcs.ValuesContext()
func (lcs *lcs[Slice, E]) ValuesContext(ctx context.Context) (values Slice, err error) {
	if lcs.values != nil {
		return lcs.values, nil
	}

	pairs, err := lcs.IndexPairsContext(ctx)
	if err != nil {
		return nil, err
	}

	values = make(Slice, len(pairs))
	for i, pair := range pairs {
		values[i] = lcs.left[pair.Left]
	}
	lcs.values = values

	return values, nil
}

// Table implements Lcs.Left()
func (lcs *lcs[Slice, E]) Left() (leftValues Slice) {
	leftValues = lcs.left
	return
}

// Table implements Lcs.Right()
func (lcs *lcs[Slice, E]) Right() (rightValues Slice) {
	rightValues = lcs.right
	return
}
