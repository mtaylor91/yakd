package system

import (
	"context"

	"github.com/sirupsen/logrus"
)

var Local *LocalSystem

func init() {
	Local = &LocalSystem{context.Background(), logrus.StandardLogger()}
}
