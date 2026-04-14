package providers

import (
	"github.com/grayguava/formseal-sync/providers/cloudflare"
	"github.com/grayguava/formseal-sync/providers/types"
)

// registry maps provider names to implementations.
// To add Supabase: import the package and uncomment the entry.
var registry = map[string]types.Provider{
	"cloudflare": &cloudflare.Provider{},
	// "supabase": &supabase.Provider{},
}

// Get returns the provider for the given name, or nil if unknown.
func Get(name string) types.Provider {
	return registry[name]
}

// List returns all registered provider names.
func List() []string {
	names := make([]string, 0, len(registry))
	for k := range registry {
		names = append(names, k)
	}
	return names
}