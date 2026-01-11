package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func (app *CLI) newDocsCmd(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:    "docs",
		Short:  "Generate documentation for Kevin CLI",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			docDir := "./docs/cli"
			if err := os.MkdirAll(docDir, 0755); err != nil {
				return err
			}

			if err := doc.GenMarkdownTree(root, docDir); err != nil {
				return err
			}

			fmt.Printf("Documentation generated in %s\n", docDir)
			return nil
		},
	}
}
