package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/rdrahul/termcap/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	// Used for flags.
	cfgFile string
	closed  bool
)

// recordCmd handles terminal recording.
var termcap = &cobra.Command{
	Use:   Application,
	Short: fmt.Sprintf("%s -/ - version=%s (%s)\r\n", Application, Version, runtime.Version()),
	Long:  "A tool for sharing live recording of your terminal",
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := startRecording(); err != nil {
			utils.Er(err)
			return err
		}
		return nil
	},
}

func startRecording() error {

	//get the shell
	shell := utils.GetShell()
	fmt.Printf(shell)
	//launch this shell with attached pseudoterminal(pty)
	recTerminal := exec.Command(shell)
	ptyFile, err := pty.Start(recTerminal)
	if err != nil {
		return err
	}

	//wait for the command
	go func() {
		defer func() {
			ptyFile.Close()
		}()
		recTerminal.Wait()
		closed = true
	}()

	fmt.Print("Recording started!\r\n")

	// Make the terminal to raw mode.
	//In raw mode, characters are directly read from and written to the device without any translation or interpretation by the operating system.
	origState, err := terminal.MakeRaw(int(os.Stdin.Fd()))

	defer func() {
		_ = terminal.Restore(int(os.Stdin.Fd()), origState)
	}()

	//trapping the resize event
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINFO)

	go func() {
		for range ch {
			println("jer")
			if err := pty.InheritSize(os.Stdin, ptyFile); err != nil {
				log.Printf("error resizing pty: %v", err)
			}
		}
	}()
	ch <- syscall.SIGINFO

	//read from stdin and put it to pty
	readChannel := make(chan string)
	go func(read chan string) {
		//create a input buffer
		inputBuffer := make([]byte, 4096)
		for {
			chars, err := os.Stdin.Read(inputBuffer)
			if err != nil {
				if err == io.EOF {
					err = nil
				} else {
					fmt.Println("channeler read error")
					utils.Er(err)
				}
				break
			}
			if !closed {

				if _, err = ptyFile.Write(inputBuffer[0:chars]); err != nil {

					log.Println(err)
					break
				}
			} else {

				read <- string(inputBuffer[0:chars])

			}

		}
	}(readChannel)

	recordTime := time.Now()
	bufout := make([]byte, 4096)
	// lines := []string{}
	for {

		nr, err := ptyFile.Read(bufout)
		if err != nil {
			println("Sdfs")
			break
		}
		// tstamp := int64(time.Since(recordTime).Nanoseconds() / int64(time.Millisecond))
		line := string(bufout[0:nr])

		// Write to STDOUT
		if _, err = os.Stdout.WriteString(line); err != nil {

			log.Println(err)
			break
		}

		recordTime = time.Now()
	}

	_ = terminal.Restore(int(os.Stdin.Fd()), origState)

	// on exit provide options
	HandleExit(readChannel)

	fmt.Printf("Record Time => %s", recordTime)
	return nil

}

//HandleExit : handles the exit part
func HandleExit(channel chan string) {

loop:
	for {
		select {
		case stdin, ok := <-channel:
			if !ok {
				break loop
			} else {
				text := strings.Replace(stdin, "\n", "", -1)
				println(text)
			}
		case <-time.After(1 * time.Second):
			fmt.Println("\r\n[bold]Ten second timeout hit, exiting...")
			break loop
		}
	}
}

//Execute : runs the cobra cli
func Execute() {
	if err := termcap.Execute(); err != nil {
		utils.Er(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	termcap.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.termcap.yaml)")
	viper.SetDefault("license", "apache")

}
