package main

import (
	"github.com/sirupsen/logrus"

	"github.com/tahcohcat/snippety/internal/cli/cobra"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	cobra.Execute()
}
