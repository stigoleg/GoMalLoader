package loader

type LoaderConfig struct {
	Mode          string
	Source        string
	Path          string
	URL           string
	AESKey        string
	Obfuscated    bool
	TargetProcess string
	SelfDelete    bool
}

type ShellcodeLoader interface {
	Run(cfg LoaderConfig) error
}

type RemoteInjector interface {
	Inject(cfg LoaderConfig) error
}

type ReflectiveLoader interface {
	Load(cfg LoaderConfig) error
}
