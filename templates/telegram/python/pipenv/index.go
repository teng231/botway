package pipenv

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/abdfnx/botway/constants"
	"github.com/abdfnx/botway/templates"
	"github.com/abdfnx/botway/templates/telegram/python"
	"github.com/abdfnx/looker"
)


func TelegramPythonPipenv(botName string) {
	pythonPath := "python3"

	if runtime.GOOS == "windows" {
		pythonPath = "python"
	}

	_, err := looker.LookPath(pythonPath)
	pipenv, perr := looker.LookPath("pipenv")

	if err != nil {
		fmt.Print(constants.FAIL_BACKGROUND.Render("ERROR"))
		fmt.Println(constants.FAIL_FOREGROUND.Render(" python is not installed"))
	} else if perr != nil {
		fmt.Print(constants.FAIL_BACKGROUND.Render("ERROR"))
		fmt.Println(constants.FAIL_FOREGROUND.Render(" pipenv is not installed"))
	} else {
		resourcesFile := os.WriteFile(filepath.Join(botName, "resources.md"), []byte(python.Resources()), 0644)

		if resourcesFile != nil {
			log.Fatal(resourcesFile)
		}

		pipenvInstall := pipenv + " install python-telegram-bot botway.py pyyaml cryptography PySocks ujson"

		cmd := exec.Command("bash", "-c", pipenvInstall)

		if runtime.GOOS == "windows" {
			cmd = exec.Command("powershell.exe", pipenvInstall)
		}

		cmd.Dir = botName
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()

		if err != nil {
			log.Printf("error: %v\n", err)
		}

		mainFile := os.WriteFile(filepath.Join(botName, "src", "main.py"), []byte(python.MainPyContent()), 0644)
		dockerFile := os.WriteFile(filepath.Join(botName, "Dockerfile"), []byte(DockerfileContent(botName)), 0644)

		if mainFile != nil {
			log.Fatal(mainFile)
		}

		if dockerFile != nil {
			log.Fatal(dockerFile)
		}

		templates.CheckProject(botName, "telegram")
	}
}
