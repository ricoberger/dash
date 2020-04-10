package render

import (
	"context"
	"errors"
	"github.com/ricoberger/dash/pkg/datasource"
	"strconv"
	"time"

	"github.com/ricoberger/dash/pkg/dashboard"
	fLog "github.com/ricoberger/dash/pkg/log"
	"github.com/ricoberger/dash/pkg/render/utils"
	"github.com/ricoberger/dash/pkg/render/widget"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

var (
	ErrNoDashboards = errors.New("no dashboards were provided")
)

func Run(datasources map[string]datasource.Client, dashboards []dashboard.Dashboard, initialInterval, initialRefresh string) error {
	// Check if there was at least one dashboard provided. This is required for the storage implementation, because we
	// choose the first dashboard as the initial one.
	// When the check succeeded we create the storage, which holds the current state of dash.
	if len(dashboards) == 0 {
		return ErrNoDashboards
	}

	storage, err := utils.NewStorage(datasources, dashboards, initialInterval, initialRefresh)
	if err != nil {
		return err
	}

	// Initialize termdash.
	// We create the statusbar, modal and the grid. The initial view shows the statusbar and grid. If an item from the
	// statusbar is selected the modal content will be rendered instead of the grid.
	t, err := termbox.New()
	if err != nil {
		return err
	}
	defer t.Close()

	statusbar, err := widget.NewStatusbar(t.Size().X, storage)
	if err != nil {
		return err
	}
	modal, err := widget.NewModal(storage)
	if err != nil {
		return err
	}

	gridOpts := widget.GridLayout(storage)

	c, err := container.New(t, container.SplitHorizontal(container.Top(container.PlaceWidget(statusbar)), container.Bottom(gridOpts...), container.SplitFixed(1)), container.ID("layout"))
	if err != nil {
		return err
	}

	var modalActive bool
	var previousKey keyboard.Key

	ctx, cancel := context.WithCancel(context.Background())

	ticker := time.NewTicker(storage.GetRefresh())
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				if !modalActive {
					fLog.Debugf("refresh was triggered")
					storage.RefreshInterval()
					gridOpts = widget.GridLayout(storage)
					c.Update("layout", container.SplitHorizontal(container.Top(container.PlaceWidget(statusbar)), container.Bottom(gridOpts...), container.SplitFixed(1)))
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	keyboardSubscriber := func(k *terminalapi.Keyboard) {
		fLog.Debugf("key %s was pressed", k.Key)
		switch k.Key {
		case 'q', keyboard.KeyCtrlC:
			cancel()
		case keyboard.KeyEnter:
			if modalActive {
				modalType, err := modal.Select()
				if err == nil {
					if modalType == widget.ModalTypeDatasource {
						storage.RefreshInterval()
					} else if modalType == widget.ModalTypeDashboard {
						storage.RefreshInterval()
					} else if modalType == widget.ModalTypeVariable {
						storage.RefreshInterval()
					} else if modalType == widget.ModalTypeInterval {
						storage.RefreshInterval()
					} else if modalType == widget.ModalTypeRefresh {
						ticker = time.NewTicker(storage.GetRefresh())
					}

					modalActive = false
					statusbar.Update(t.Size().X)
					gridOpts = widget.GridLayout(storage)
					c.Update("layout", container.SplitHorizontal(container.Top(container.PlaceWidget(statusbar)), container.Bottom(gridOpts...), container.SplitFixed(1)))
				}
			}
		case 'd':
			modalActive = modal.Show(&widget.ModalOptions{Type: widget.ModalTypeDashboard, VariableIndex: 0})
			c.Update("layout", container.SplitHorizontal(container.Top(container.PlaceWidget(statusbar)), container.Bottom(container.PlaceWidget(modal)), container.SplitFixed(1)))
		case 's':
			modalActive = modal.Show(&widget.ModalOptions{Type: widget.ModalTypeDatasource, VariableIndex: 0})
			c.Update("layout", container.SplitHorizontal(container.Top(container.PlaceWidget(statusbar)), container.Bottom(container.PlaceWidget(modal)), container.SplitFixed(1)))
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if modalActive {
				modal.SelectIndex(string(k.Key))
			} else if previousKey == 'v' && k.Key != '0' {
				variableIndex, err := strconv.Atoi(string(k.Key))
				if err == nil && variableIndex <= len(storage.Dashboard().Variables) {
					modalActive = modal.Show(&widget.ModalOptions{Type: widget.ModalTypeVariable, VariableIndex: variableIndex - 1})
					c.Update("layout", container.SplitHorizontal(container.Top(container.PlaceWidget(statusbar)), container.Bottom(container.PlaceWidget(modal)), container.SplitFixed(1)))
				}
			}
		case 'i':
			modalActive = modal.Show(&widget.ModalOptions{Type: widget.ModalTypeInterval, VariableIndex: 0})
			c.Update("layout", container.SplitHorizontal(container.Top(container.PlaceWidget(statusbar)), container.Bottom(container.PlaceWidget(modal)), container.SplitFixed(1)))
		case 'r':
			modalActive = modal.Show(&widget.ModalOptions{Type: widget.ModalTypeRefresh, VariableIndex: 0})
			c.Update("layout", container.SplitHorizontal(container.Top(container.PlaceWidget(statusbar)), container.Bottom(container.PlaceWidget(modal)), container.SplitFixed(1)))
		case keyboard.KeyEsc:
			modalActive = false
			storage.RefreshInterval()
			gridOpts = widget.GridLayout(storage)
			c.Update("layout", container.SplitHorizontal(container.Top(container.PlaceWidget(statusbar)), container.Bottom(gridOpts...), container.SplitFixed(1)))
		}

		previousKey = k.Key
	}

	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(keyboardSubscriber)); err != nil {
		return err
	}

	return nil
}
