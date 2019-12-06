// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"strings"

	"github.com/moov-io/watchman/pkg/ofac"

	"github.com/abadojack/whatlanggo"
	"github.com/bbalet/stopwords"
	"github.com/pariz/gountries"
)

const (
	minConfidence = 0.50
)

func removeStopwords(in string, lang whatlanggo.Lang) string {
	return strings.TrimSpace(stopwords.CleanString(strings.ToLower(in), lang.Iso6391(), false))
}

// detectLanguage will return a guess as to the appropiate language a given SDN's name
// is written in. The addresses must be linked to the SDN whose name is detected.
func detectLanguage(in string, addrs []*ofac.Address) whatlanggo.Lang {
	info := whatlanggo.Detect(in)
	if info.IsReliable() {
		// Return the detected language if whatlanggo is confident enough
		return info.Lang
	}

	if len(addrs) == 0 {
		// If no addresses are associated to this text blob then fallback to English
		return whatlanggo.Eng
	}

	// Return the countries primary langauge associated to the primary address for this SDN.
	//
	// TODO(adam): Should we do this only if there's one address? If there are multiple should we
	// fallback to English or a mixed set?
	country, err := gountries.New().FindCountryByName(addrs[0].Country)
	if len(country.Languages) == 0 || err != nil {
		return whatlanggo.Eng
	}

	// If the language is spoken in the country and we're somewhat confident in the original detection
	// then return that language.
	if info.Confidence > minConfidence {
		for key, _ := range country.Languages {
			if strings.EqualFold(key, info.Lang.Iso6393()) {
				return info.Lang
			}
		}
	}

	// How should we pick the language for countries with multiple languages? A hardcoded map?
	// What if we found the language whose name is closest to the country's name and returned that?
	//
	// Should this fallback be the mixed set that contains stop words from several popular languages
	// in the various data sets?

	return whatlanggo.Eng
}
