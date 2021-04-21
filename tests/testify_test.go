package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	s := new(GormSuite)
	suite.Run(t, s)
}
