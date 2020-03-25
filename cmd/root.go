package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	// Used for flags.
	CfgFile string

	Daemon bool

	RootCmd = &cobra.Command{
		Use:   "status",
		Short: "Status Masterlab Socket",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("config file:", CfgFile)
			fmt.Println("Masterlab status is normal")
		},
	}
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start Masterlab Socket",
		Run: func(cmd *cobra.Command, args []string) {
			if Daemon {
				command := exec.Command("masterlab_socket", "start")
				if err := command.Start(); err != nil { // 运行命令
					log.Fatal(err)
				}
				fmt.Printf("Masterlab start, [PID] %d running...\n", command.Process.Pid)
				err := ioutil.WriteFile("gonne.lock", []byte(fmt.Sprintf("%d", command.Process.Pid)), 0666)
				if err != nil { // 运行命令
					log.Fatal(err)
				}
				Daemon = false
				os.Exit(0)

			} else {
				fmt.Println("Masterlab start")
			}
		},
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "MasterlabSocket的版本号",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("MasterlabSocket v1.0")
		},
	}

	stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop Masterlab Socket",
		Run: func(cmd *cobra.Command, args []string) {
			strb, _ := ioutil.ReadFile("gonne.lock")
			fmt.Println("kill gonne.lock: ", string(strb))
			if runtime.GOOS != "windows" {
				command := exec.Command("kill", string(strb))
				if err := command.Start(); err != nil {
					fmt.Println(err)
				}
			} else {
				command := exec.Command("taskkill", "/F", "/T", "/PID", string(strb))
				if err := command.Start(); err != nil {
					fmt.Println(err)
					file, _ := exec.LookPath(os.Args[0])
					println("try kill by name", file)
					exec.Command("taskkill.exe", "/f", "/im",file)
				}
			}
			println("Masterlab Socket stop")
			os.Exit(1)
		},
	}
)


func init() {
	cobra.OnInitialize(initConfig)
	startCmd.Flags().StringVarP(&CfgFile, "config", "c", "", "config file (default is $HOME/config.toml)")
	startCmd.Flags().BoolVarP(&Daemon, "daemon", "d", false, "is daemon?")
	RootCmd.AddCommand(startCmd)
	RootCmd.AddCommand(stopCmd)
	RootCmd.AddCommand(versionCmd)

}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	// Use config file from the default.
	CfgFile = strings.Replace(CfgFile, " ", "",  -1)
	CfgFile = strings.Replace(CfgFile, "\n", "",  -1)
	if CfgFile == "" || CfgFile==`./config.toml` {
		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)
		curPath := filepath.Dir(path)
		CfgFile = fmt.Sprintf("%s%s", curPath, `/config.toml`)
	}
	fmt.Println("cfgFile:", CfgFile)
	fmt.Println("Deamonize:", Daemon)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
