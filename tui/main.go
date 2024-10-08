package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Command struct {
	Name            *string
	RomSetName      *string
	ZipFilepath     *string
	BinFilepath     *string
	MraFilepath     *string
	DiffZipFilepath *string
}

func processCommand(command string) Command {
	switch command {
	case "Decrypt":
		romName := RomPicker{}.New()
		romZipFilepath := DecryptPicker{}.New()
		return Command{&command, &romName, &romZipFilepath, nil, nil, nil}
	case "Encrypt":
		romName := RomPicker{}.New()
		romZipFilepath, binFilepath := EncryptPicker{}.New()
		return Command{&command, &romName, &romZipFilepath, &binFilepath, nil, nil}
	case "Patch":
		romName := RomPicker{}.New()
		romZipFilepath := Filepicker{}.New("Choose the ROM .zip to patch:", []string{".zip"})
		mraFilepath := Filepicker{}.New("Choose the .mra file to patch with:", []string{".mra"})
		return Command{&command, &romName, &romZipFilepath, nil, &mraFilepath, nil}
	case "Diff":
		romName := RomPicker{}.New()
		romZipFilepath := Filepicker{}.New("Choose the first ROM .zip to compare:", []string{".zip"})
		diffZipFilepath := Filepicker{}.New("Choose the second ROM.zip to compare:", []string{".zip"})
		return Command{&command, &romName, &romZipFilepath, nil, nil, &diffZipFilepath}
	}
	return Command{}
}

func StartTui() Command {
	m := InitializeCommandPicker()
	n, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	commandStr := n.(CommandPickerModel).choice
	return processCommand(commandStr)
}
