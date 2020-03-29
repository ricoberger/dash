package render

import (
	"errors"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ricoberger/dash/pkg/dashboard"
	"github.com/ricoberger/dash/pkg/render/utils"
	"github.com/ricoberger/dash/pkg/render/widget"

	ui "github.com/gizak/termui/v3"
)

var (
	ErrNoDashboards = errors.New("no dashboards were provided")
)

func Run(dashboards []dashboard.Dashboard, initialInterval, initialRefresh string) error {
	if len(dashboards) == 0 {
		return ErrNoDashboards
	}

	storage, err := utils.NewStorage(dashboards, initialInterval, initialRefresh)
	if err != nil {
		return err
	}

	// Initialize termui.
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	var modalActive bool
	var previousKey string

	termWidth, termHeight := ui.TerminalDimensions()
	statusbar := widget.NewStatusbar(termWidth, termHeight, storage)
	grid := widget.NewGrid(termWidth, termHeight, storage)
	grid.Refresh()
	modal := widget.NewModal(termWidth, termHeight, storage)

	ui.Render(statusbar, grid, modal)

	sigTerm := make(chan os.Signal, 2)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(storage.GetRefresh())

	for {
		select {
		case <-sigTerm:
			return nil
		case <-ticker.C:
			storage.RefreshInterval()
			grid = widget.NewGrid(termWidth, termHeight, storage)
			ui.Render(statusbar, grid, modal)
		case e := <-ui.PollEvents():
			switch e.ID {
			case "q", "<C-c>":
				return nil
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				statusbar.SetRect(0, 0, payload.Width, 1)
				grid.SetRect(0, 1, payload.Width, payload.Height)
				modal.SetDimensions(payload.Width, payload.Height)
				ui.Render(statusbar, grid, modal)
			case "j", "<Down>":
				if modalActive {
					modal.ScrollDown()
					ui.Render(statusbar, grid, modal)
				}
			case "k", "<Up>":
				if modalActive {
					modal.ScrollUp()
					ui.Render(statusbar, grid, modal)
				}
			case "<Enter>":
				if modalActive {
					modalType := modal.Select()
					if modalType == widget.ModalTypeDashboard {
						storage.RefreshInterval()
						grid = widget.NewGrid(termWidth, termHeight, storage)
					} else if modalType == widget.ModalTypeVariable {
						storage.RefreshInterval()
						grid = widget.NewGrid(termWidth, termHeight, storage)
					} else if modalType == widget.ModalTypeInterval {
						storage.RefreshInterval()
						grid = widget.NewGrid(termWidth, termHeight, storage)
					} else if modalType == widget.ModalTypeRefresh {
						ticker = time.NewTicker(storage.GetRefresh())
					}

					modal.Hide()
					modalActive = false
					ui.Render(statusbar, grid, modal)
				}
			case "<Escape>":
				if modalActive {
					modal.Hide()
					modalActive = false
					ui.Render(statusbar, grid, modal)
				}
			case "d":
				modalActive = modal.Show(&widget.ModalOptions{Type: widget.ModalTypeDashboard, VariableIndex: 0})
				ui.Clear()
				ui.Render(statusbar, grid, modal)
			case "1", "2", "3", "4", "5", "6", "7", "8", "9":
				if previousKey == "v" {
					variableIndex, err := strconv.Atoi(e.ID)
					if err == nil {
						modalActive = modal.Show(&widget.ModalOptions{Type: widget.ModalTypeVariable, VariableIndex: variableIndex})
						ui.Render(statusbar, grid, modal)
					}
				}
			case "i":
				modalActive = modal.Show(&widget.ModalOptions{Type: widget.ModalTypeInterval, VariableIndex: 0})
				ui.Render(statusbar, grid, modal)
			case "r":
				modalActive = modal.Show(&widget.ModalOptions{Type: widget.ModalTypeRefresh, VariableIndex: 0})
				ui.Render(statusbar, grid, modal)
			}

			previousKey = e.ID
		}
	}
}
