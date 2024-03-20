package main

import (
	"fmt"
	"os"

	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/userns"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var migrateCommand = cli.Command{
	Name:  "live-migrate",
	Usage: "live migrate a running container",
	ArgsUsage: `<container-id>

Where "<container-id>" is the name for the instance of the container to be
checkpointed.`,
	Description: `The migrate command saves the state of the container instance and restore it in the other server.`,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "destination", Value: "", Usage: "ip address for destination server"},
		cli.StringFlag{Name: "image-path", Value: "", Usage: "path for saving criu image files"},
		cli.StringFlag{Name: "work-path", Value: "", Usage: "path for saving work files and logs"},
		cli.StringFlag{Name: "parent-path", Value: "", Usage: "path for previous criu image files in pre-dump"},
		cli.BoolFlag{Name: "receiver", Usage: "set up for recieving criu image files"},
		cli.BoolFlag{Name: "local-test", Usage: "testing live migration for local server"},
		cli.BoolFlag{Name: "leave-running", Usage: "leave the process running after checkpointing"},
		cli.BoolFlag{Name: "tcp-established", Usage: "allow open tcp connections"},
		cli.BoolFlag{Name: "ext-unix-sk", Usage: "allow external unix sockets"},
		cli.BoolFlag{Name: "shell-job", Usage: "allow shell jobs"},
		cli.BoolFlag{Name: "lazy-pages", Usage: "use userfaultfd to lazily restore memory pages"},
		cli.IntFlag{Name: "status-fd", Value: -1, Usage: "criu writes \\0 to this FD once lazy-pages is ready"},
		cli.StringFlag{Name: "page-server", Value: "", Usage: "ADDRESS:PORT of the page server"},
		cli.BoolFlag{Name: "file-locks", Usage: "handle file locks, for safety"},
		cli.BoolFlag{Name: "pre-dump", Usage: "dump container's memory information only, leave the container running after this"},
		cli.StringFlag{Name: "manage-cgroups-mode", Value: "", Usage: "cgroups mode: 'soft' (default), 'full' and 'strict'"},
		cli.StringSliceFlag{Name: "empty-ns", Usage: "create a namespace, but don't restore its properties"},
		cli.BoolFlag{Name: "auto-dedup", Usage: "enable auto deduplication of memory images"},
	},
	Action: func(context *cli.Context) error {

		// fmt.Println("You have successfully called the migrate command!")
		// if err := checkArgs(context, 1, exactArgs); err != nil {
		// 	return err
		// }
		// XXX: Currently this is untested with rootless containers.
		if os.Geteuid() != 0 || userns.RunningInUserNS() {
			logrus.Warn("runc checkpoint is untested with rootless containers")
		}

		options := criuOptions(context)
		if context.Bool("receiver") {
			receive()
			if err := setEmptyNsMask(context, options); err != nil {
				return err
			}
		} else {
			container, err := getContainer(context)
			// newContext := cli.NewContext(context.App, context.Flags(), context.App.Metadata)
			// newContext.SetArgs(context.Args()[2:])
			// container, err := getContainer(newContext)
			if err != nil {
				return err
			}
			status, err := container.Status()
			if err != nil {
				return err
			}
			if status == libcontainer.Created || status == libcontainer.Stopped {
				fatal(fmt.Errorf("Container cannot be checkpointed in %s state", status.String()))
			}

			// options := criuOptions(context)
			options.PreDump = true
			if !(options.LeaveRunning || options.PreDump) {
				// destroy container unless we tell CRIU to keep it
				defer destroy(container)
			}
			// these are the mandatory criu options for a container
			setPageServer(context, options)
			setManageCgroupsMode(context, options)
			if err := setEmptyNsMask(context, options); err != nil {
				return err
			}
			for i := 0; i < 3; i++ {
				container.Checkpoint(options)
				transfer("cr", 1, context.String("destination"), context.String("image-path"))

			}
			options.PreDump = false
			container.Checkpoint(options)
		}

		if context.Bool("local-test") || context.Bool("receiver") {
			status, err := startContainer(context, CT_ACT_RESTORE, options)
			if err != nil {
				return err
			}
			// exit with the container's exit status so any external supervisor is
			// notified of the exit with the correct exit status.
			os.Exit(status)
		} else {
			if context.String("destination") == "" || context.String("image-path") == "" {
				return fmt.Errorf("unable to migrate files with empty destination or empty file path")
			}
			transfer("cr", 2, context.String("destination"), context.String("image-path"))
		}
		return nil
	},
}
