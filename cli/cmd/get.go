package cmd

import (
	"fmt"
	"rommi/brain"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get triggerword | services | actions [SERVICENAME] | sentences [SERVICENAME [ACTIONNAME]]",
	Short: "Tells Rommi a command.",
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := brain.New()
		if err != nil {
			return err
		}
		b.Run()
		if len(args) == 0 {
			return cmd.Usage()
		}
		switch args[0] {
		case "triggerword":
			// fmt.Println("Getting TriggerWord")
			triggerWord := b.GetTriggerWord()
			// fmt.Printf("TriggerWord is '%s'", triggerWord)
			fmt.Println(triggerWord)
		case "services":
			// fmt.Println("Getting Services")
			services := b.GetServices()
			// fmt.Println("Available Services are")
			for _, service := range services {
				fmt.Printf("%s\n", service)
			}
		case "actions":
			// fmt.Println("Getting Actions")
			if len(args) < 2 {
				services := b.GetServices()
				// fmt.Println("Available Actions are")
				for _, service := range services {
					actions := b.GetServiceActions(service)
					for _, action := range actions {
						fmt.Printf("%s\n", action)
					}
				}
			} else {
				actions := b.GetServiceActions(args[1])
				// fmt.Println("Available Actions are")
				for _, action := range actions {
					fmt.Printf("%s\n", action)
				}
			}
		case "sentences":
			// fmt.Println("Getting Sentences")
			if len(args) < 2 {
				// fmt.Println("Available Sentences are")
				sentences := b.GetSentences()
				for _, s := range sentences {
					fmt.Printf("%s\n", s)
				}
			} else if len(args) == 2 {
				service := args[1]
				// fmt.Println("Available Sentences are")
				sentences := b.GetServiceSentences(service)
				for _, s := range sentences {
					fmt.Printf("%s\n", s)
				}
			} else {
				service := args[1]
				action := args[2]
				sentences := b.GetActionSentences(service, action)
				// fmt.Println("Available Sentences are")
				for _, s := range sentences {
					fmt.Printf("%s\n", s)
				}
			}
		default:
			return cmd.Usage()
		}
		return nil
	},
}
