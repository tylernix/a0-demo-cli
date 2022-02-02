package branding

import (
	_ "embed"
)

var (
	//go:embed data/default-template.liquid
	DefaultTemplate string

	//go:embed data/footer-template.liquid
	FooterTemplate string

	//go:embed data/image-template.liquid
	ImageTemplate string
)
