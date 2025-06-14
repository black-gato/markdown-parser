/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/black-gato/markdown-parser/internal"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	filesFlag = "files"
	tagFlag   = "tag"
)

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunParse(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.Flags().StringSliceP(filesFlag, "f", []string{}, "passing markdown files that you want to parse")

	parseCmd.Flags().StringSliceP(tagFlag, "", []string{}, "search for tag in a file ie. [[Hello]]")

	parseCmd.MarkFlagRequired(filesFlag)

	parseCmd.MarkFlagRequired(tagFlag)

}

func RunParse(cmd *cobra.Command, args []string) error {
	var files []string
	inputFiles, err := cmd.Flags().GetStringSlice(filesFlag)
	if err != nil {
		return err
	}
	for _, f := range inputFiles {

		exstention := filepath.Ext(f)

		if strings.ToLower(exstention) != ".md" {
			logger.Printf("This file is not a markdown file %s\n", f)
			continue
		}

		_, err := os.Stat(f)
		if err != nil {
			logger.Printf("This file is not real %s\n %v", f, err)
			continue
		}
		files = append(files, f)
	}

	tag, err := cmd.Flags().GetStringSlice(tagFlag)
	if err != nil {
		return err
	}

	reference, err := internal.Parse(files, tag)
	if err != nil {
		logger.Fatal(err)
		return err
	}
	fmt.Println(reference)
	return nil
}
