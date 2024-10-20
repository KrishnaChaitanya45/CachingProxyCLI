package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cacheproxy",
	Short: "CacheProxy is a CLI tool for caching HTTP requests",
	Long: `CacheProxy is a CLI tool that allows you to run a caching proxy server,
where you can cache HTTP requests and responses locally and fetch them from the cache
instead of querying the origin server repeatedly.`,
	// The action when root command is called directly
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Fatalf("ERROR WHILE GETTING PORT %s", err)
		}
		origin, err := cmd.Flags().GetString("origin")
		if err != nil {
			log.Fatalf("ERROR WHILE GETTING ORIGIN %s", err)
		}
		clearCache, err := cmd.Flags().GetBool("clear-cache")
		if err != nil {
			log.Fatalf("ERROR WHILE GETTING CLEAR CACHE %s", err)
		}
		cache := make(map[string]string)
		requestHandler := RequestHandler{origin: origin, cache: cache}
		if clearCache {
			requestHandler.ClearCache()
		} else {
			server := http.Server{
				Addr:    fmt.Sprintf(":%s", port),
				Handler: &requestHandler,
			}
			fmt.Printf("🦾 PROXY SERVER IS LISTENING ON %s", port)
			server.ListenAndServe()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
	RootCmd.PersistentFlags().StringP("port", "p", "3000", "Port on which the server will run")
	RootCmd.PersistentFlags().StringP("origin", "o", "http://dummyjson.com", "Origin server URL")
	RootCmd.PersistentFlags().BoolP("clear-cache", "c", false, "Clear the cache on startup")
	RootCmd.MarkFlagRequired("port")
	RootCmd.MarkFlagRequired("origin")
}
