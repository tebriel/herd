package main

import (
	"fmt"

	"github.com/seveas/herd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:                   "list [--oneline] glob [filters] [<+|-> glob [filters]...]",
	Short:                 "Query your datasources and list hosts matching globs and filters",
	Example:               "  herd list *.site2.example.com os=Debian",
	DisableFlagsInUseLine: true,
	RunE:                  runList,
}

func init() {
	listCmd.Flags().Bool("oneline", false, "List all hosts on one line, separated by commas")
	listCmd.Flags().String("separator", ",", "String separating hostnames in --oneline mode")
	listCmd.Flags().StringSlice("attributes", []string{}, "Show not only the names, but also the specified attributes")
	listCmd.Flags().Bool("all-attributes", false, "List hosts with all their attributes")
	listCmd.Flags().Bool("csv", false, "Output in csv format")
	listCmd.Flags().Bool("header", true, "Print attribute names in a header line before printing host data")
	listCmd.Flags().String("template", "", "Template to use for showing hosts")
	listCmd.Flags().StringSlice("count", []string{}, "Show counts for the values of these attributes")
	viper.BindPFlag("OneLine", listCmd.Flags().Lookup("oneline"))
	viper.BindPFlag("Separator", listCmd.Flags().Lookup("separator"))
	viper.BindPFlag("AllAttributes", listCmd.Flags().Lookup("all-attributes"))
	viper.BindPFlag("Attributes", listCmd.Flags().Lookup("attributes"))
	viper.BindPFlag("Csv", listCmd.Flags().Lookup("csv"))
	viper.BindPFlag("Header", listCmd.Flags().Lookup("header"))
	viper.BindPFlag("Template", listCmd.Flags().Lookup("template"))
	viper.BindPFlag("Count", listCmd.Flags().Lookup("count"))
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	splitAt := cmd.ArgsLenAtDash()
	if splitAt != -1 {
		return fmt.Errorf("Command provided, but list mode doesn't support that")
	}
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	engine, err := setupScriptEngine(nil)
	if err != nil {
		return err
	}
	defer engine.End()
	if len(args) == 0 {
		args = append(args, "*")
	}
	if err = engine.ParseCommandLine(args, splitAt); err != nil {
		logrus.Error(err.Error())
		return err
	}
	engine.Execute()
	opts := herd.HostListOptions{
		OneLine:       viper.GetBool("OneLine"),
		Separator:     viper.GetString("Separator"),
		Csv:           viper.GetBool("csv"),
		Attributes:    viper.GetStringSlice("Attributes"),
		AllAttributes: viper.GetBool("AllAttributes"),
		Header:        viper.GetBool("Header"),
		Align:         true,
		Template:      viper.GetString("Template"),
		Count:         viper.GetStringSlice("Count"),
		SortByCount:   !viper.IsSet("Sort"),
	}
	engine.Ui.PrintHostList(engine.Runner.GetHosts(), opts)
	return nil
}
