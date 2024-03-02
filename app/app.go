package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/urfave/cli"
	"golang.org/x/sys/windows"
)

func Init() *cli.App {
	app := cli.NewApp()
	app.Name = "Local.Host Symlink"
	app.Usage = "Create symlink for files and directories inside Local.Host document root directory"

	app.Commands = []cli.Command{
		{
			Name:    "directory",
			Aliases: []string{"d"},
			Usage:   "Create symlink for a directory",
			Action:  createDirectorySymlink,
		},
		{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "Create symlink for a file",
			Action:  createFileSymlink,
		},
	}

	return app
}

func createDirectorySymlink(c *cli.Context) error {
	target := c.Args().First()
	dir := filepath.Base(target)
	link := "C:\\local.host\\www\\" + dir

	fmt.Println("Creating symlink for directory ", target, " to ", link)

	if !checkAdmin() {
		runMeElevated()
	}

	command := exec.Command("cmd.exe", "/c", "mklink", "/D", link, target)
	command.Run()

	return nil
}

func createFileSymlink(c *cli.Context) error {
	target := c.Args().First()
	file := filepath.Base(target)
	link := "C:\\local.host\\www\\" + file

	fmt.Println("Creating symlink for file ", target, " to ", link)

	if !checkAdmin() {
		runMeElevated()
	}

	command := exec.Command("cmd.exe", "/c", "mklink", link, target)
	command.Run()

	return nil
}

func runMeElevated() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
}

func checkAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}
