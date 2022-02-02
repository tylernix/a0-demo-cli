package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tylernix/a0-demo-cli/auth0"
	"github.com/tylernix/a0-demo-cli/display"
	"gopkg.in/auth0.v5/management"
)

var cfgFile string

// Viper uses the mapstructure package under the hood for unmarshaling values
// We use the mapstructure tags to specify the name of each config field.
// Fields must be capitalized for viper.Unmarshal to work
type params struct {
	Domain        string `mapstructure:"AUTH0_DOMAIN"`
	Client_id     string `mapstructure:"AUTH0_CLIENT_ID"`
	Client_secret string `mapstructure:"AUTH0_CLIENT_SECRET"`
}

var A0config params

// cli provides all the foundational things for all the commands in the CLI,
// specifically:
//
// 1. A management API instance (e.g. go-auth0/auth0)
// 2. A renderer (which provides ansi, coloring, etc).
//
// In addition, it stores a reference to all the flags passed, e.g.:
//
// 1. --format
// 2. --tenant
// 3. --debug
//
type cli struct {
	// core primitives exposed to command builders.
	api *auth0.API
	//authenticator *auth.Authenticator
	renderer *display.Renderer
	//tracker       *analytics.Tracker
	// set of flags which are user specified.
	debug   bool
	tenant  string
	format  string
	force   bool
	noInput bool
	noColor bool

	// config state management.
	//initOnce sync.Once
	errOnce error
	path    string
	//config   config
}

const (
	envPrefix = "AUTH0"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cli := &cli{
		renderer: display.NewRenderer(),
	}

	rootCmd := buildRootCmd(cli)

	addPersistentFlags(rootCmd, cli)
	addSubcommands(rootCmd, cli)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
func buildRootCmd(cli *cli) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "a0-demo",
		Short: "CLI tool to help with giving a compelling customer demo in Auth0",
		Long: `CLI tool to help with giving a compelling customer demo in Auth0.
	
			A longer description that spans multiple lines and likely contains
			examples and usage of using your application. For example:
			
			Cobra is a CLI library for Go that empowers applications.
			This application is a tool to generate the needed files
			to quickly create a Cobra application.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			//return initializeConfig(cmd)
			return cli.setup()
		},
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()

			// Print the final resolved value from binding cobra flags and viper config
			fmt.Fprintln(out, "Your Auth0 domain is:", A0config.Domain)
			fmt.Fprintln(out, "Your client_id is:", A0config.Client_id)

		},
	}
	return rootCmd
}

func addPersistentFlags(rootCmd *cobra.Command, cli *cli) {
	rootCmd.PersistentFlags().StringVar(&A0config.Domain, "AUTH0_DOMAIN", "", "Set Auth0 Tenant Domain")
	rootCmd.PersistentFlags().StringVar(&A0config.Client_id, "AUTH0_CLIENT_ID", "", "Set Auth0 Client Id")
	rootCmd.PersistentFlags().StringVar(&A0config.Client_secret, "AUTH0_CLIENT_SECRET", "", "Set Auth0 Client Secret")

	// rootCmd.PersistentFlags().StringVar(&cli.tenant,
	// 	"tenant", cli.config.DefaultTenant, "Specific tenant to use.")

	// rootCmd.PersistentFlags().BoolVar(&cli.debug,
	// 	"debug", false, "Enable debug mode.")

	// rootCmd.PersistentFlags().StringVar(&cli.format,
	// 	"format", "", "Command output format. Options: json.")

	// rootCmd.PersistentFlags().BoolVar(&cli.force,
	// 	"force", false, "Skip confirmation.")

	// rootCmd.PersistentFlags().BoolVar(&cli.noInput,
	// 	"no-input", false, "Disable interactivity.")

	// rootCmd.PersistentFlags().BoolVar(&cli.noColor,
	// 	"no-color", false, "Disable colors.")

}

func addSubcommands(rootCmd *cobra.Command, cli *cli) {
	// order of the comamnds here matters
	// so add new commands in a place that reflect its relevance or relation with other commands:

	// rootCmd.AddCommand(loginCmd(cli))
	// rootCmd.AddCommand(logoutCmd(cli))
	rootCmd.AddCommand(configCmd(cli))
	// rootCmd.AddCommand(tenantsCmd(cli))
	// rootCmd.AddCommand(appsCmd(cli))
	// rootCmd.AddCommand(usersCmd(cli))
	// rootCmd.AddCommand(rulesCmd(cli))
	// rootCmd.AddCommand(actionsCmd(cli))
	// rootCmd.AddCommand(apisCmd(cli))
	// rootCmd.AddCommand(rolesCmd(cli))
	// rootCmd.AddCommand(organizationsCmd(cli))
	rootCmd.AddCommand(brandingCmd(cli))
	// rootCmd.AddCommand(ipsCmd(cli))
	// rootCmd.AddCommand(quickstartsCmd(cli))
	// rootCmd.AddCommand(testCmd(cli))
	// rootCmd.AddCommand(logsCmd(cli))

	// keep completion at the bottom:
	//rootCmd.AddCommand(completionCmd(cli))

}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	v := viper.New()
	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name "config" (without extension).
		v.AddConfigPath(home)
		v.AddConfigPath(".")
		v.SetConfigName("local")
		v.SetConfigType("env")
	}

	v.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err == nil {
		//fmt.Fprintln(os.Stderr, "Using config file:", v.ConfigFileUsed())
	}

	v.SetEnvPrefix(envPrefix)

	//bindFlags(rootCmd, v)

	if err := v.Unmarshal(&A0config); err != nil {
		fmt.Println(err)
	}
}

func (c *cli) setup() error {
	var (
		m   *management.Management
		err error
	)
	//fmt.Println(A0config)
	if A0config.Client_id != "" && A0config.Client_secret != "" {
		m, err = management.New(A0config.Domain, management.WithClientCredentials(A0config.Client_id, A0config.Client_secret))
	}
	if err != nil {
		return err
	}

	c.api = auth0.NewAPI(m)

	return nil
}
