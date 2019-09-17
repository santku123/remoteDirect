package remotedirect

import (
        "github.com/spf13/cobra"
        "opendev.org/airship/airshipctl/pkg/environment"
         "opendev.org/airship/airshipctl/cmd/isogen/remotedirect/target"
)


// Added a command to send action to remote BMC 
func NewRemotedirectCommand(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
        remotedirectActionCmd := &cobra.Command{
                Use:   "remotedirect",
                Short: "airshipctl remotedirect",
                Run: func(cmd *cobra.Command, args []string) {
                        name, _ := cmd.Flags().GetString("name")
                        if name == "" {
                                name = "remotedirect"
                        }

                },
        }

        remotedirectActionCmd.AddCommand(target.NewRemoteDirectTarget(rootSettings))
        return remotedirectActionCmd 
}

