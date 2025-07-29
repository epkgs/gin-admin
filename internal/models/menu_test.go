package models

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/menus.json
var rawMenus []byte

func TestMarshal(t *testing.T) {
	var menus Menus
	err := json.Unmarshal(rawMenus, &menus)
	assert.NoError(t, err)

	fmt.Printf("menus len: %d \n", len(menus))

	byts, err := json.MarshalIndent(menus, "", "  ")

	fmt.Printf("menus: %v \n", string(byts))

	assert.NoError(t, err)
}
