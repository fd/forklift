package deploy

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/kr/hk/term"
	"io"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func (cmd *Deploy) run_post_push_commands() error {
	for _, command := range cmd.Config.PostPushCommands {
		err := cmd.run_post_push_command(command)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cmd *Deploy) run_post_push_command(command string) error {
	if cmd.DryRun {
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

	err := cmd.Http("POST", &req, &resp, "/apps/%s/dynos", cmd.Config.Name)
	if err != nil {
		return err
	}

	return do_rendervous(resp.AttachURL)
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

	if term.IsTerminal(os.Stdin) && term.IsTerminal(os.Stdout) {
		err = term.MakeRaw(os.Stdin)
		if err != nil {
			return err
		}
		defer term.Restore(os.Stdin)

		sig := make(chan os.Signal)
		signal.Notify(sig, os.Signal(syscall.SIGQUIT), os.Interrupt)
		go func() {
			defer term.Restore(os.Stdin)
			for sg := range sig {
				switch sg {
				case os.Interrupt:
					cn.Write([]byte{3})
				case os.Signal(syscall.SIGQUIT):
					cn.Write([]byte{28})
				default:
					panic("not reached")
				}
			}
		}()
	}

	errc := make(chan error)
	cp := func(a io.Writer, b io.Reader) {
		_, err := io.Copy(a, b)
		errc <- err
	}

	go cp(os.Stdout, br)
	go cp(cn, os.Stdin)
	if err = <-errc; err != nil {
		return err
	}

	return nil
}
