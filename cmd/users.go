package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/tylernix/a0-demo-cli/ansi"
	"github.com/tylernix/a0-demo-cli/users"
	"gopkg.in/auth0.v5/management"
)

var (
	connection = Flag{
		Name:       "Connection",
		LongForm:   "connection",
		ShortForm:  "c",
		Help:       "Name of the database connection this user should be imported in.",
		IsRequired: false,
	}
	filePath = Flag{
		Name:       "File Path",
		LongForm:   "file-path",
		ShortForm:  "f",
		Help:       "Path to file containing user import example.",
		IsRequired: false,
	}
	upsert = Flag{
		Name:       "Upsert",
		LongForm:   "upsert",
		ShortForm:  "u",
		Help:       "When set to false, pre-existing users that match on email address, user ID, or username will fail. When set to true, pre-existing users that match on any of these fields will be updated, but only with upsertable attributes.",
		IsRequired: false,
	}
)

func usersCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Manage resources for users",
		Long:  "Manage resources for users.",
	}

	cmd.AddCommand(importUsersCmd(cli))

	return cmd
}

func importUsersCmd(cli *cli) *cobra.Command {
	var inputs struct {
		Connection   string
		ConnectionId string
		Upsert       bool
	}
	cmd := &cobra.Command{
		Use:   "import",
		Args:  cobra.NoArgs,
		Short: "Import users from schema",
		Long: `Import users from schema. Issues a Create Import Users Job. 
The file size limit for a bulk import is 500KB. You will need to start multiple imports if your data exceeds this size.`,
		Example: `a0-demo users import --connection "Username-Password-Authentication"
a0-demo users import -c "Username-Password-Authentication" --upsert true`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Select from the available connection types
			// Users API currently support database connections
			if err := connection.Select(cmd, &inputs.Connection, cli.connectionPickerOptions(), nil); err != nil {
				return err
			}
			conn, err := cli.api.Connection.ReadByName(inputs.Connection)
			inputs.ConnectionId = *conn.ID

			// The getConnReqUsername returns the value for the requires_username field for the selected connection
			// The result will be used to determine whether to prompt for username
			// conn := cli.getConnReqUsername(auth0.StringValue(&inputs.connection))
			// requireUsername := auth0.BoolValue(conn)
			// if requireUsername {
			// 	if err := userUsername.Ask(cmd, &inputs.Username, nil); err != nil {
			// 		return err
			// 	}
			// 	a.Username = &inputs.Username
			// }

			var (
				userImportOptions = pickerOptions{
					{"Basic", users.BasicExample},
					{"Custom Password Hash", users.CustomPasswordHashExample},
					{"MFA Factors", users.MFAFactors},
				}
			)

			var qs = []*survey.Question{
				{
					Name: "users",
					Prompt: &survey.Select{
						Message: "User Import Examples:",
						Options: userImportOptions.labels(),
					},
					Validate: survey.Required,
				},
			}

			example := ""
			err = survey.Ask(qs, &example)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

			// Convert json to map
			jsonstr := userImportOptions.getValue(example)
			var jsonmap []map[string]interface{}
			json.Unmarshal([]byte(jsonstr), &jsonmap)

			err = ansi.Waiting(func() error {
				return cli.api.Jobs.ImportUsers(&management.Job{
					ConnectionID: &inputs.ConnectionId,
					Users:        jsonmap,
					Upsert:       &inputs.Upsert,
				})
			})

			cli.renderer.Heading("User(s) imported")
			fmt.Println(jsonstr)

			return nil
		},
	}

	connection.RegisterString(cmd, &inputs.Connection, "")
	//filePath.RegisterString(cmd, &inputs.filepath, "")
	upsert.RegisterBool(cmd, &inputs.Upsert, false)

	return cmd
}

func (c *cli) connectionPickerOptions() []string {
	var res []string

	list, err := c.api.Connection.List()
	if err != nil {
		fmt.Println(err)
	}
	for _, conn := range list.Connections {
		if conn.GetStrategy() == "auth0" {
			res = append(res, conn.GetName())
		}
	}
	return res
}

// This is a workaround to get the requires_username field nested inside Options field
func (c *cli) getConnReqUsername(s string) *bool {
	conn, err := c.api.Connection.ReadByName(s)
	if err != nil {
		fmt.Println(err)
	}
	res := fmt.Sprintln(conn.Options)

	opts := &management.ConnectionOptions{}
	if err := json.Unmarshal([]byte(res), &opts); err != nil {
		fmt.Println(err)
	}

	return opts.RequiresUsername
}
