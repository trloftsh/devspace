package flags

import (
	"github.com/loft-sh/devspace/pkg/devspace/config/loader"
	flag "github.com/spf13/pflag"
)

// GlobalFlags is the flags that contains the global flags
type GlobalFlags struct {
	Silent bool
	NoWarn bool
	Debug  bool

	OverrideName             string
	Namespace                string
	KubeContext              string
	Profiles                 []string
	ProfileRefresh           bool
	DisableProfileActivation bool
	SwitchContext            bool
	ConfigPath               string
	Vars                     []string

	InactivityTimeout int

	Flags *flag.FlagSet
}

// ToConfigOptions converts the globalFlags into config options
func (gf *GlobalFlags) ToConfigOptions() *loader.ConfigOptions {
	profiles := []string{}
	profiles = append(profiles, gf.Profiles...)
	return &loader.ConfigOptions{
		OverrideName:             gf.OverrideName,
		Profiles:                 profiles,
		ProfileRefresh:           gf.ProfileRefresh,
		DisableProfileActivation: gf.DisableProfileActivation,
		Vars:                     gf.Vars,
	}
}

// SetGlobalFlags applies the global flags
func SetGlobalFlags(flags *flag.FlagSet) *GlobalFlags {
	globalFlags := &GlobalFlags{
		Vars:  []string{},
		Flags: flags,
	}

	flags.StringVar(&globalFlags.OverrideName, "override-name", "", "If specified will override the devspace.yaml name")
	flags.BoolVar(&globalFlags.NoWarn, "no-warn", false, "If true does not show any warning when deploying into a different namespace or kube-context than before")
	flags.BoolVar(&globalFlags.Debug, "debug", false, "Prints the stack trace if an error occurs")
	flags.BoolVar(&globalFlags.Silent, "silent", false, "Run in silent mode and prevents any devspace log output except panics & fatals")

	flags.StringVar(&globalFlags.ConfigPath, "config", "", "The devspace config file to use")
	flags.StringSliceVarP(&globalFlags.Profiles, "profile", "p", []string{}, "The DevSpace profiles to apply. Multiple profiles are applied in the order they are specified")
	flags.BoolVar(&globalFlags.ProfileRefresh, "profile-refresh", false, "If true will pull and re-download profile parent sources")
	flags.BoolVar(&globalFlags.DisableProfileActivation, "disable-profile-activation", false, "If true will ignore all profile activations")
	flags.BoolVarP(&globalFlags.SwitchContext, "switch-context", "s", false, "Switches and uses the last kube context and namespace that was used to deploy the DevSpace project")
	flags.StringVarP(&globalFlags.Namespace, "namespace", "n", "", "The kubernetes namespace to use")
	flags.StringVar(&globalFlags.KubeContext, "kube-context", "", "The kubernetes context to use")
	flags.StringSliceVar(&globalFlags.Vars, "var", []string{}, "Variables to override during execution (e.g. --var=MYVAR=MYVALUE)")

	flags.IntVar(&globalFlags.InactivityTimeout, "inactivity-timeout", 0, "Minutes the current user is inactive (no mouse or keyboard interaction) until DevSpace will exit automatically. 0 to disable. Only supported on windows and mac operating systems")
	return globalFlags
}
