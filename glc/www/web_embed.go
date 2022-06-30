package www

import "embed"

var (
	//go:embed web/dist/*
	Static embed.FS
)
