package notify

import (
	"fmt"
	"os/exec"
)

var command = exec.Command

func (n *Notify) Send() error {
	notifyCmdName := "terminal-notifier"

	notifyCmd, err := exec.LookPath(notifyCmdName)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("(%s)%s", n.severity, n.title)

	notifyCommand := command(notifyCmd, "-title", title, "-message", n.message)

	return notifyCommand.Run()
}
