package app

import (
	"context"
	"fmt"
	"os"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/spinner"
	"github.com/ksckaan1/templ-iconify/internal/core/domain"
	"github.com/ksckaan1/templ-iconify/internal/core/service"
	"github.com/ksckaan1/templ-iconify/internal/core/tui"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type RootCmd struct {
	cmd         *cobra.Command
	iconService *service.IconService
	saveDir     string
	workerCount int
	workChan    chan *domain.Icon

	program *tea.Program
}

func NewRootCmd(iconService *service.IconService) *RootCmd {
	return &RootCmd{
		cmd: &cobra.Command{
			Use:     "templ-iconify",
			Short:   "Download icons from Iconify",
			Long:    `Download icons from Iconify and generate templates for them.`,
			Example: `templ-iconify "mdi:home" "solar:*" "mdi:home-*" "*:*" -o ./icons/`,
			Args:    cobra.MinimumNArgs(1),
		},
		iconService: iconService,
		saveDir:     "./icons/",
		workerCount: 10,
		workChan:    make(chan *domain.Icon, 1),
	}
}

func (r *RootCmd) Run(ctx context.Context) error {
	r.cmd.RunE = r.runE
	r.cmd.Flags().StringVarP(&r.saveDir, "out", "o", r.saveDir, "Output directory")
	r.cmd.Flags().IntVarP(&r.workerCount, "worker", "w", r.workerCount, "Worker count")
	err := r.cmd.ExecuteContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *RootCmd) runE(cmd *cobra.Command, args []string) error {
	foundIcons := make([]*domain.Icon, 0)

	err := spinner.New().
		Title("Finding icons...").
		Context(cmd.Context()).
		ActionWithErr(func(ctx context.Context) error {
			icons, err := r.iconService.FindIcons(ctx, args...)
			if err != nil {
				return err
			}
			foundIcons = icons
			return nil
		}).
		Run()
	if err != nil {
		return err
	}

	model := tui.NewDownloadModel(len(foundIcons))

	r.program = tea.NewProgram(model)
	if err != nil {
		return fmt.Errorf("tea.NewProgram: %w", err)
	}

	var eg errgroup.Group

	eg.Go(func() error {
		<-model.OnStart
		wg := new(sync.WaitGroup)
		for range r.workerCount {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := r.worker(cmd.Context())
				if err != nil {
					r.program.Quit()
				}
			}()
		}
		go func() {
			for _, icon := range foundIcons {
				r.workChan <- icon
			}
			close(r.workChan)
		}()
		wg.Wait()
		r.program.Quit()
		return nil
	})

	eg.Go(func() error {
		_, err = r.program.Run()
		if err != nil {
			return fmt.Errorf("program.Run: %w", err)
		}
		os.Exit(0)
		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (r *RootCmd) worker(ctx context.Context) error {
	for work := range r.workChan {
		err := r.iconService.DownloadIcon(ctx, work, r.saveDir)
		if err != nil {
			return err
		}
		r.program.Send(tui.SaveMsg{Icon: work})
	}
	return nil
}
