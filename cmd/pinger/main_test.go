package main

import (
	"bytes"
	"errors"
	"log"
	"testing"
	"github.com/stretchr/testify/suite"
)

type PingerTests struct {
	suite.Suite
}

func TestPinger(t *testing.T) {
	suite.Run(t, &PingerTests{})
}

func (s *PingerTests) Test_handleError() {
	var logs bytes.Buffer
	mockLogger := log.New(&logs, "mock", log.LstdFlags)
	mockError := errors.New("___ mock error")
	handleError(mockError, mockLogger)
	s.Contains(logs.String(), "___ mock error")
}
