package project

import (
	"fmt"

	"github.com/openshift/odo/pkg/log"
	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/spf13/cobra"

	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const getRecommendedCommandName = "get"

var (
	getExample = ktemplates.Examples(`
	# Get the active project
	%[1]s 
	`)

	getLongDesc = ktemplates.LongDesc(`Get the active project`)

	getShortDesc = `Get the active project`
)

// ProjectGetOptions encapsulates the options for the odo project get command
type ProjectGetOptions struct {

	// if supplied then only print the project name
	projectShortFlag bool

	// generic context options common to all commands
	*genericclioptions.Context
}

// NewProjectGetOptions creates a ProjectGetOptions instance
func NewProjectGetOptions() *ProjectGetOptions {
	return &ProjectGetOptions{}
}

// Complete completes ProjectGetOptions after they've been created
func (pgo *ProjectGetOptions) Complete(name string, cmd *cobra.Command, args []string) (err error) {
	pgo.Context = genericclioptions.NewContext(cmd)
	return
}

// Validate validates the parameters of the ProjectGetOptions
func (pgo *ProjectGetOptions) Validate() (err error) {
	return
}

// Run runs the project get command
func (pgo *ProjectGetOptions) Run() (err error) {
	project := pgo.Context.Project

	if pgo.projectShortFlag {
		fmt.Print(project)
	} else {
		log.Infof("The current project is: %v", project)
	}

	return
}

// NewCmdProjectGet creates the project get command
func NewCmdProjectGet(name, fullName string) *cobra.Command {
	o := NewProjectGetOptions()

	projectGetCmd := &cobra.Command{
		Use:     name,
		Short:   getShortDesc,
		Long:    getLongDesc,
		Example: fmt.Sprintf(getExample, fullName),
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	projectGetCmd.Flags().BoolVarP(&o.projectShortFlag, "short", "q", false, "If true, display only the project name")

	return projectGetCmd
}
