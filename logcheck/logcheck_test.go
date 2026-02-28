package logcheck_test

import (
	"testing"

	"github.com/f1rezy/loglint/logcheck"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestRule1(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), logcheck.Analyzer, "rule1")
}

func TestRule2(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), logcheck.Analyzer, "rule2")
}

func TestRule3(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), logcheck.Analyzer, "rule3")
}

func TestRule4(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), logcheck.Analyzer, "rule4")
}
