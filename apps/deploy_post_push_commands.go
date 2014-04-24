package apps

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net/url"
	"os"
)

func (app *App) run_post_push_commands() error {
	for _, command := range app.PostPushCommands {
		err := app.run_post_push_command(command)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *App) run_post_push_command(command string) error {
	if app.config.DryRun {
		fmt.Printf("Run: %s\n - skipped (dry run)\n", command)
		return nil
	}

	fmt.Printf("Run: %s\n", command)

	var (
		req = struct {
			Command string `json:"command"`
			Attach  bool   `json:"attach"`
		}{
			Command: command,
			Attach:  true,
		}

		resp struct {
			AttachURL string `json:"attach_url"`
		}
	)

	err := app.HttpV3("POST", &req, &resp, "/apps/%s/dynos", app.AppName)
	if err != nil {
		fmt.Printf("error=%s", err)
		return err
	}

	err = do_rendervous(resp.AttachURL)
	if err != nil {
		fmt.Printf("error=%s", err)
		return err
	}

	return nil
}

func do_rendervous(rendezvous_url string) error {
	u, err := url.Parse(rendezvous_url)
	if err != nil {
		return err
	}

	cn, err := tls.Dial("tcp", u.Host, nil)
	if err != nil {
		return err
	}
	defer cn.Close()

	br := bufio.NewReader(cn)

	_, err = io.WriteString(cn, u.Path[1:]+"\r\n")
	if err != nil {
		return err
	}

	for {
		_, pre, err := br.ReadLine()
		if err != nil {
			return err
		}
		if !pre {
			break
		}
	}

	if _, err := io.Copy(os.Stdout, br); err != nil {
		return err
	}

	return nil
}
