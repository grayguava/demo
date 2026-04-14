package supabase

// Provider implements types.Provider for Supabase.
// To activate:
//  1. Implement Validate and Fetch below
//  2. In providers/registry.go, import this package and uncomment the supabase entry

// import "github.com/grayguava/formseal-sync/providers/types"
// type Provider struct{}
// func (p *Provider) Validate(cfg *types.ProviderConfig) error { ... }
// func (p *Provider) Fetch(cfg *types.ProviderConfig, outputPath string) (int, int, error) { ... }