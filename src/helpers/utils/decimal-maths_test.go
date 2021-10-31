package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type inputs struct {
	amount   float64
	incoming float64
}

func TestCalcSum(t *testing.T) {

	tests := []struct {
		name string
		given inputs
		expected float64
	}{
		{"Sum 1 + 2", inputs{1, 2}, 3},
		{"Sum 0 + 2", inputs{0, 2}, 2},
		{"Sum 0.5 + 0.2", inputs{0.5, 0.2}, 0.7},
		{"Sum 1 + 0.5", inputs{1, 0.5}, 1.5},
		{"Sum 0.3 + 0.001", inputs{0.3, 0.001}, 0.301},
	}

	asserts := assert.New(t)
	for _, test := range tests {
		asserts.Equal(test.expected, DecimalMaths().Sum(test.given.amount, test.given.incoming), test.name)
	}
}

func TestCalcSub(t *testing.T) {

	tests := []struct {
		name string
		given inputs
		expected float64
	}{
		{"Sub 1 - 2", inputs{1, 2}, -1},
		{"Sub 0 - 0", inputs{0, 0}, 0},
		{"Sub 3 - 0.7", inputs{3, 0.7}, 2.3},
		{"Sub 0.7 - 0.5", inputs{0.7, 0.5}, 0.2},
		{"Sub 0.3 - 0.001", inputs{0.3, 0.001}, 0.299},
	}

	asserts := assert.New(t)
	for _, test := range tests {
		asserts.Equal(test.expected, DecimalMaths().Sub(test.given.amount, test.given.incoming), test.name)
	}
}

func TestCalcMul(t *testing.T) {

	tests := []struct {
		name string
		given inputs
		expected float64
	}{
		{"Mul 1 * 2", inputs{1, 2}, 2},
		{"Mul 0 * 0", inputs{0, 0}, 0},
		{"Mul 2 * 0.7", inputs{2, 0.7}, 1.4},
		{"Mul 0.7 * 0.5", inputs{0.7, 0.5}, 0.35},
		{"Mul 0.3 * 0.001", inputs{0.5, 0.002}, 0.001},
	}

	asserts := assert.New(t)
	for _, test := range tests {
		asserts.Equal(test.expected, DecimalMaths().Mul(test.given.amount, test.given.incoming), test.name)
	}
}