"""help command - show help"""

def run():
    print('''formseal-sync - fetch encrypted form submissions

Usage:
  formseal-sync config           Show configuration
  formseal-sync set <key> <val>  Set configuration
  formseal-sync fetch            Fetch ciphertexts
  formseal-sync status           Show stats
  formseal-sync help             Show this help

Config keys:
  provider         - cloudflare or supabase
  cf-token         - Cloudflare API token
  cf-namespace     - Cloudflare KV namespace ID
  sb-url           - Supabase project URL
  sb-key           - Supabase service key
  sb-table         - Supabase table name
  output           - Output folder path

Examples:
  formseal-sync set provider cloudflare
  formseal-sync set cf-token cfun_xxx
  formseal-sync set output C:\\Users\\you\\data
  formseal-sync fetch''')