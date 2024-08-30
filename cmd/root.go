package cmd

import (
	"fmt"
	"oasis/pkg/config"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// rootCmd is the root command of backup service
	rootCmd = &cobra.Command{
		Use:   "oasis",
		Short: "oasis service provide cli for making clean block function easier",
		Long:  `oasis service provide cli for making clean block function easier`,
	}
)

var logo = `
                                                                                                                        
     OOOOOOOOO                         AAA                       SSSSSSSSSSSSSSS      IIIIIIIIII        SSSSSSSSSSSSSSS 
   OO:::::::::OO                      A:::A                    SS:::::::::::::::S     I::::::::I      SS:::::::::::::::S
 OO:::::::::::::OO                   A:::::A                  S:::::SSSSSS::::::S     I::::::::I     S:::::SSSSSS::::::S
O:::::::OOO:::::::O                 A:::::::A                 S:::::S     SSSSSSS     II::::::II     S:::::S     SSSSSSS
O::::::O   O::::::O                A:::::::::A                S:::::S                   I::::I       S:::::S            
O:::::O     O:::::O               A:::::A:::::A               S:::::S                   I::::I       S:::::S            
O:::::O     O:::::O              A:::::A A:::::A               S::::SSSS                I::::I        S::::SSSS         
O:::::O     O:::::O             A:::::A   A:::::A               SS::::::SSSSS           I::::I         SS::::::SSSSS    
O:::::O     O:::::O            A:::::A     A:::::A                SSS::::::::SS         I::::I           SSS::::::::SS  
O:::::O     O:::::O           A:::::AAAAAAAAA:::::A                  SSSSSS::::S        I::::I              SSSSSS::::S 
O:::::O     O:::::O          A:::::::::::::::::::::A                      S:::::S       I::::I                   S:::::S
O::::::O   O::::::O         A:::::AAAAAAAAAAAAA:::::A                     S:::::S       I::::I                   S:::::S
O:::::::OOO:::::::O        A:::::A             A:::::A        SSSSSSS     S:::::S     II::::::II     SSSSSSS     S:::::S
 OO:::::::::::::OO        A:::::A               A:::::A       S::::::SSSSSS:::::S     I::::::::I     S::::::SSSSSS:::::S
   OO:::::::::OO         A:::::A                 A:::::A      S:::::::::::::::SS      I::::::::I     S:::::::::::::::SS 
     OOOOOOOOO          AAAAAAA                   AAAAAAA      SSSSSSSSSSSSSSS        IIIIIIIIII      SSSSSSSSSSSSSSS   
                                                                                                                        
Oasis is a powerful CLI tool for generating Golang clean architecture modules. Designed to streamline your development process, Oasis is specially designed for oasis_boilerplate which you can create within this CLI tool.
For more info visit: https://github.com/mhshajib/oasis
`

func init() {
	cobra.OnInitialize(initConfig)
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("Loading configurations")
	if err := config.Init(); err != nil {
		logrus.Warn("Failed to load configuration")
		logrus.Fatal(err)
	}
	logrus.Info("Configurations loaded successfully!")
}
