package cmd

import (
	"fmt"

	"github.com/tylernix/a0-demo-cli/ansi"
	"github.com/tylernix/a0-demo-cli/branding"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"gopkg.in/auth0.v5/management"
)

// brandingCmd represents the branding command
func brandingCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branding",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("branding called")
			res, apiErr := cli.api.Branding.Read()
			if apiErr != nil {
				fmt.Println(apiErr)
			}
			fmt.Println(res)
		},
	}

	cmd.AddCommand(showBrandingCmd(cli))
	//cmd.AddCommand(updateBrandingCmd(cli))
	cmd.AddCommand(templateCmd(cli))
	return cmd
}

func templateCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "templates",
		Short: "Manage custom page templates",
		Long:  `Manage custom [page templates](https://auth0.com/docs/universal-login/new-experience/universal-login-page-templates). This requires a custom domain to be configured for the tenant.`,
	}

	//cmd.SetUsageTemplate(resourceUsageTemplate())
	cmd.AddCommand(showBrandingTemplateCmd(cli))
	cmd.AddCommand(updateBrandingTemplateCmd(cli))
	cmd.AddCommand(deleteBrandingTemplateCmd(cli))
	return cmd
}

func showBrandingCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Args:    cobra.NoArgs,
		Short:   "Display the custom branding settings for Universal Login",
		Long:    "Display the custom branding settings for Universal Login.",
		Example: "auth0 branding show",
		RunE: func(cmd *cobra.Command, args []string) error {
			var branding *management.Branding // Load app by id

			if err := ansi.Waiting(func() error {
				var err error
				branding, err = cli.api.Branding.Read()
				return err
			}); err != nil {
				return fmt.Errorf("Unable to load branding settings due to an unexpected error: %w", err)
			}

			cli.renderer.BrandingShow(branding)

			return nil
		},
	}

	return cmd
}

func showBrandingTemplateCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Args:    cobra.NoArgs,
		Short:   "Display the custom template for Universal Login",
		Long:    "Display the custom template for Universal Login.",
		Example: "auth0 branding templates show",
		RunE: func(cmd *cobra.Command, args []string) error {
			var template *management.BrandingUniversalLogin // Load app by id

			if err := ansi.Waiting(func() error {
				var err error
				template, err = cli.api.Branding.UniversalLogin()
				return err
			}); err != nil {
				return fmt.Errorf("Unable to load the Universal Login template due to an unexpected error: %w", err)
			}

			cli.renderer.Heading("template")
			fmt.Println(*template.Body)

			return nil
		},
	}

	return cmd
}

func updateBrandingTemplateCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Args:    cobra.NoArgs,
		Short:   "Update the custom template for Universal Login",
		Long:    "Update the custom template for Universal Login.",
		Example: "auth0 branding templates update",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				customTemplateOptions = pickerOptions{
					{"Basic", branding.DefaultTemplate},
					{"Login box + image", branding.ImageTemplate},
					{"Page footers", branding.FooterTemplate},
				}
			)

			var qs = []*survey.Question{
				{
					Name: "template",
					Prompt: &survey.Select{
						Message: "Template:",
						Options: customTemplateOptions.labels(),
					},
					Validate: survey.Required,
				},
			}

			template := ""
			var htmlbody string

			err := survey.Ask(qs, &template)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

			htmlbody = customTemplateOptions.getValue(template)

			// editor := false
			// prompt := &survey.Confirm{
			// 	Message: "Open template in an editor?",
			// }
			// survey.AskOne(prompt, &editor)
			// if editor != false {
			// 	prompt := &survey.Editor{
			// 		Message:  "Edit Auth0 Page Template",
			// 		FileName: "*.liquid",
			// 	}
			// 	survey.AskOne(prompt, &htmlbody)
			// }

			err = ansi.Waiting(func() error {
				return cli.api.Branding.SetUniversalLogin(&management.BrandingUniversalLogin{
					Body: &htmlbody,
				})
			})

			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func deleteBrandingTemplateCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Args:    cobra.NoArgs,
		Short:   "Delete the custom template for Universal Login",
		Long:    "Delete the custom template for Universal Login.",
		Example: "auth0 branding templates delete",
		RunE: func(cmd *cobra.Command, args []string) error {

			err := ansi.Waiting(func() error {
				return cli.api.Branding.DeleteUniversalLogin()
			})

			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
