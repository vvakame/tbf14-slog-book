//go:build tools
// +build tools

package main

// from https://github.com/golang/go/issues/25922#issuecomment-412992431

import (
	_ "github.com/vvakame/ptproc/cmd/ptproc"
)
