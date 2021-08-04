package ginjet

// Options for JetRender
type RenderOptions struct {
	TemplateDir string
	ContentType string
}


// Default options
func DefaultOptions() *RenderOptions {
	return &RenderOptions{
		TemplateDir: "./views",
		ContentType: "text/html; charset=utf-8",
	}
}
