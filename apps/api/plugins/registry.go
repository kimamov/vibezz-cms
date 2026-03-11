package plugins

import (
	"github.com/vibezz/cms/internal/plugin"
	"github.com/vibezz/cms/plugins/news"
)

// All returns the list of plugins to load.
// To add a custom plugin: implement plugin.Plugin in a subpackage (e.g. plugins/myplugin),
// then import and append it here.
func All() []plugin.Plugin {
	return []plugin.Plugin{
		news.New(),
	}
}
