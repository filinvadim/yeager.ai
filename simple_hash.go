package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// questions:
// - is strict sequence increment required? not clear
// - what to do if given nil sequence number with non-nil value?
// - order for same sequence fragments?

type integrityError string

func (e integrityError) Error() string {
	return string(e)
}

const (
	ErrIntegrityVerification integrityError = "Error: Data integrity verification failed.\""
	ErrMissingFragment       integrityError = "Error: Missing fragment of the sequence."
	ErrBrokenOrder           integrityError = "Error: Sequence order is broken."

	dataKey                = "data"
	hashKey                = "hash"
	missingDataPlaceholder = "..."
)

func NewFragment(data, hash string) Fragment {
	return map[string]string{dataKey: data, hashKey: hash}
}

type Fragment map[string]string

func (f Fragment) Data() string {
	return f[dataKey]
}
func (f Fragment) Hash() string {
	return f[hashKey]
}
func (f Fragment) IsNil() bool {
	return f == nil
}

// readability aliases
type (
	sequence = *int

	verifiedFragment struct {
		int
		f Fragment
	}
)

func reconstructData(unorderedFragments map[sequence]Fragment) (string, error) {
	verifiedFragments, err := verify(unorderedFragments)
	if errors.Is(err, ErrIntegrityVerification) { // critical error
		return "", ErrIntegrityVerification
	}

	sort.SliceStable(verifiedFragments, func(i, j int) bool {
		return verifiedFragments[i].int <= verifiedFragments[j].int
	})

	return assemble(verifiedFragments), err
}

func verify(fragmentsMap map[sequence]Fragment) (_ []verifiedFragment, baseErr error) {
	var verifiedFragments = make([]verifiedFragment, 0, len(fragmentsMap))

	// in Go map access is always random
	for i, fr := range fragmentsMap {
		if err := verifySequence(i); err != nil {
			baseErr = errors.Join(baseErr, err)
			// append policy for missing order seq
			verifiedFragments = append(verifiedFragments, verifiedFragment{len(verifiedFragments) - 1, fr})
			continue
		}

		fr, err := verifyMissing(fr, *i) // add placeholder if missing - reconstructed data might still be readable
		if err != nil {
			baseErr = errors.Join(baseErr, err)
			verifiedFragments = append(verifiedFragments, verifiedFragment{*i, fr})
			continue
		}

		if err = verifyHash(fr); err != nil {
			return nil, err
		}

		// append used in case of position collision (same index). But what order?
		verifiedFragments = append(verifiedFragments, verifiedFragment{*i, fr})
	}
	return verifiedFragments, baseErr
}

func verifySequence(i *int) error {
	if i == nil {
		return ErrBrokenOrder
	}
	return nil
}

func verifyMissing(fr Fragment, i int) (_ Fragment, err error) {
	if fr.IsNil() {
		return NewFragment(missingDataPlaceholder, ""), fmt.Errorf("%w: %d", ErrMissingFragment, i)
	}
	return fr, nil
}

func verifyHash(fr Fragment) error {
	calculatedHash := simpleHash(fr.Data())
	if fr.Hash() != "" && calculatedHash != fr.Hash() {
		return ErrIntegrityVerification // critical
	}
	return nil
}

func assemble(orderedFragments []verifiedFragment) string {
	reconstructedData := new(strings.Builder)
	for i := 0; i < len(orderedFragments); i++ {
		fragment := orderedFragments[i].f
		data := fragment.Data()

		reconstructedData.WriteString(data)
	}

	reconstructed := reconstructedData.String()
	reconstructedData.Reset()
	return reconstructed
}

const hashLength = 30

// simple hash function
func simpleHash(input string) string {
	if input == "" {
		return ""
	}
	var (
		firstComplicator  uint64 = 31 // random value
		secondComplicator uint64 = 53 // random value
		baseHashValue     uint64 = 7  // random value
	)

	for i, char := range input {
		baseHashValue += uint64(char) * firstComplicator                            // int64(char) - symbol's ASCII-code
		baseHashValue = baseHashValue ^ (baseHashValue << 5) ^ (baseHashValue >> 3) // just random movements
		baseHashValue = baseHashValue * secondComplicator                           // random multiplication
		baseHashValue = baseHashValue ^ uint64(i)                                   // also random XOR
	}

	hashString := strconv.FormatUint(baseHashValue, 16) // hex

	// adjust to len limit
	if len(hashString) > hashLength {
		hashString = hashString[:hashLength] // cut
	} else if len(hashString) < hashLength {
		padding := strings.Repeat("#", hashLength-len(hashString))
		hashString = hashString + padding // grow
	}

	return hashString
}
