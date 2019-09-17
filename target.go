package target 

import (
        "fmt"
        "os"
        "github.com/spf13/cobra"
        "opendev.org/airship/airshipctl/pkg/environment"
        "opendev.org/airship/airshipctl/cmd/isogen/remotedirect/target/action"
)

// Build target 
func NewRemoteDirectTarget(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
        remotetarget := &cobra.Command{
                Use:   "target",
                Short: "airshppctl remotedirect target <ip:port>",
                Run: func(cmd *cobra.Command, args []string) {
                        targetName, _ := cmd.Flags().GetString("target")
                        if targetName == "" {
                                targetName = "target"
                                fmt.Println("You did not specify the target, please specify target ")
                                os.Exit(1)

                        }
                },
        }

        remotetarget.AddCommand(action.NewRemoteDirectTargetAction(rootSettings))
        return remotetarget
}

