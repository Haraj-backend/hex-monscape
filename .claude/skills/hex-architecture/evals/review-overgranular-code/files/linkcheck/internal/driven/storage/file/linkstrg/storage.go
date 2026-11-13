package linkstrg

// Storage is a file-backed adapter. (body omitted for brevity: ~40 lines of
// os.ReadFile / json encoding)
type Storage struct {
	path string
}

type Config struct {
	Path string `validate:"nonzero"`
}

func New(cfg Config) (*Storage, error) {
	return &Storage{path: cfg.Path}, nil
}
