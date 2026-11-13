package shared

import "encoding/json"

func Marshal(v any) ([]byte, error) { return json.MarshalIndent(v, "", "  ") }
