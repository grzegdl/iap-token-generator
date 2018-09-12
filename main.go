package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	iap "github.com/tonglil/iap_curl"
)

var (
	// Configuration options
	refresh     time.Duration
	filename    string
	credentials string
	client      string

	credentialsEnv = "GOOGLE_APPLICATION_CREDENTIALS"
	clientEnv      = "IAP_CLIENT_ID"
)

func init() {
	viper.AutomaticEnv()
	flags := rootCmd.Flags()

	flags.StringP("credentials", "c", "", fmt.Sprintf("The service account JSON credential [%s]", credentialsEnv))
	viper.BindPFlag("credentials", flags.Lookup("credentials"))
	viper.BindEnv("credentials", credentialsEnv)

	flags.StringP("id", "i", "", fmt.Sprintf("The IAP client ID [%s]", clientEnv))
	viper.BindPFlag("client", flags.Lookup("id"))
	viper.BindEnv("client", clientEnv)

	rootCmd.PersistentFlags().DurationVarP(&refresh, "refresh", "r", time.Duration(0), "Refresh the token on a specified interval")
	rootCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "", "Write the token to a file")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:           "iap-token-generator",
	Short:         "Generate a Bearer token for making HTTP requests to IAP-protected apps",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		i, err := iap.New(viper.GetString("credentials"), viper.GetString("client"))
		if err != nil {
			return err
		}

		token, err := i.GetToken()
		if err != nil {
			return err
		}

		if filename != "" {
			err = ioutil.WriteFile(filename, []byte(token), 0644)
			if err != nil {
				return err
			}
			fmt.Printf("Token written to %s\n", filename)
		} else {
			fmt.Println(token)
		}

		if refresh > time.Duration(0) {
			ticker := time.NewTicker(refresh)
			sigs := make(chan os.Signal, 1)
			done := make(chan bool, 1)
			errs := make(chan error, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				for {
					select {
					case <-ticker.C:
						token, err := i.GetToken()
						if err != nil {
							errs <- err
							done <- true
							return
						}

						if filename != "" {
							err = ioutil.WriteFile(filename, []byte(token), 0644)
							if err != nil {
								errs <- err
								done <- true
								return
							}
							fmt.Printf("Token written to %s\n", filename)
						} else {
							fmt.Println(token)
						}
					case sig := <-sigs:
						fmt.Println()
						fmt.Println(sig)
						done <- true
						return

					}
				}
			}()

			// Wait for signal to stop and quit
			<-done
			ticker.Stop()

			select {
			case err := <-errs:
				return err
			default:
				return nil
			}
		}

		return nil
	},
}
