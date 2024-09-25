package cps2rom

import (
	_ "embed"
	"encoding/json"
)

type Roms = map[string]RomDefinition

type RomRegion struct {
	Size       int                  `json:"size,omitempty"`
	Operations []RomRegionOperation `json:"operations,omitempty"`
}

type RomDefinition struct {
	Maincpu  RomRegion `json:"maincpu,omitempty"`
	Gfx      RomRegion `json:"gfx,omitempty"`
	Audiocpu RomRegion `json:"audiocpu,omitempty"`
	Qsound   RomRegion `json:"qsound,omitempty"`
	Key      RomRegion `json:"key,omitempty"`
}

type RomRegionOperation struct {
	Offset    int    `json:"offset,omitempty"`
	Length    int    `json:"length,omitempty"`
	Type      string `json:"type,omitempty"`
	GroupSize int    `json:"groupSize,omitempty"`
	Skip      int    `json:"skip,omitempty"`
	Reverse   bool   `json:"reverse,omitempty"`
	Filename  string `json:"filename,omitempty"`
}

//go:embed roms.json
var romsBytes []byte

var RomDefinitions *Roms

func ParseRoms() error {
	err := json.Unmarshal(romsBytes, &RomDefinitions)
	return err
}
