// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package csl

// CSL contains each record from the Consolidate Screening List, broken down by the record's original source
type CSL struct {
	// []*SSI (Sectoral Sanctions Identifications List (SSI) - Treasury Department)
	// []*EL (Entity List – Bureau of Industry and Security)
	// []*UL (Unverified List – Bureau of Industry and Security)
	// []*PSE (Foreign Sanctions Evaders (FSE) - Treasury Department)
	// []*ISN (Nonproliferation Sanctions (ISN) - State Department)
	// []*PLC (Palestinian Legislative Council List (PLC) - Treasury Department)
	// []*CAPTA (CAPTA (formerly Foreign Financial Institutions Subject to Part 561 - Treasury Department))
	// []*ADL (AECA Debarred List - State Department)
}
