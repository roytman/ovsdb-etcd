package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/roytman/ovsdb-etcd/pkg/codegenerator"
)


var (
	rootCmd = &cobra.Command{
		Use:   "generator",
		Short: "A code generator for OVSDB schema types",
		Long: `generator creates goLang code for tables defined by the OVSDB schema.`,
		Run: func(cmd *cobra.Command, args []string) {
			codegenerator.Run()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&codegenerator.SchemaFile, "schemaFile", "s", "", "input schema file")
	rootCmd.MarkPersistentFlagRequired("schemaFile")
	rootCmd.PersistentFlags().StringVarP(&codegenerator.OutputDir, "outputDir", "o", ".",  "output directory")
	rootCmd.PersistentFlags().StringVarP(&codegenerator.PkgName, "packageName", "p", "", "the package of generated files, default is database name")
}

func initConfig() {

	viper.AutomaticEnv()

}

func main() {
	Execute()
}
