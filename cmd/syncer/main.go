package main

import (
	"github.com/kjbreil/syncer/endpoint"
	"github.com/kjbreil/syncer/endpoint/settings"
	"github.com/rivo/tview"
	"log/slog"
	"net"
	"time"
)

type state struct {
	app *tview.Application

	endpointOne     *endpoint.Endpoint
	endpointOneData data
	endpointTwo     *endpoint.Endpoint
	endpointTwoData data

	updateFunc []func()
}

type data struct {
	Name string
}

func main() {
	app := tview.NewApplication()
	s := state{
		app:             app,
		endpointOneData: data{},
		endpointTwoData: data{},
	}

	var port = 45012

	peers := []net.TCPAddr{{
		IP:   net.ParseIP("10.0.2.2"),
		Port: port,
	},
	}
	peersTwo := []net.TCPAddr{
		{
			IP:   net.ParseIP("10.0.2.2"),
			Port: port,
		},
		{
			IP:   net.ParseIP("10.0.2.3"),
			Port: port,
		},
	}

	endpointOne, err := endpoint.New(&s.endpointOneData, &settings.Settings{
		Peers: peers,
		Port:  port,
	})
	if err != nil {
		panic(err)
	}
	endpointOneLogInfo := tview.NewTextView()
	endpointOneLogInfo.SetChangedFunc(func() {
		endpointOneLogInfo.ScrollToEnd()
	})
	endpointOne.SetLogger(slog.NewTextHandler(endpointOneLogInfo, nil))
	// endpointOne.SetLogger(slog.NewJSONHandler(io.Discard, nil))
	s.endpointOne = endpointOne

	endpointTwo, err := endpoint.New(&s.endpointTwoData, &settings.Settings{
		Peers:      peersTwo,
		Port:       port,
		AutoUpdate: true,
	})
	if err != nil {
		panic(err)
	}
	endpointTwoLogInfo := tview.NewTextView()
	endpointTwoLogInfo.SetChangedFunc(func() {
		endpointTwoLogInfo.ScrollToEnd()
	})
	endpointTwo.SetLogger(slog.NewTextHandler(endpointTwoLogInfo, nil))
	// endpointTwo.SetLogger(slog.NewJSONHandler(io.Discard, nil))
	s.endpointTwo = endpointTwo

	// logInfo.SetText("this is some text", false)

	grid := tview.NewGrid().
		SetRows(5).
		// SetColumns(2).
		// SetBorders(true).
		AddItem(s.makeEndpointControl("Endpoint One", s.endpointOne), 0, 0, 1, 1, 0, 0, false).
		AddItem(s.makeEndpointControl("Endpoint Two", s.endpointTwo), 0, 1, 1, 1, 0, 0, false).
		AddItem(s.makeEndpointForm(&s.endpointOneData), 1, 0, 1, 1, 0, 0, false).
		AddItem(s.makeEndpointForm(&s.endpointTwoData), 1, 1, 1, 1, 0, 0, false).
		AddItem(endpointOneLogInfo, 2, 0, 1, 1, 0, 0, false).
		AddItem(endpointTwoLogInfo, 2, 1, 1, 1, 0, 0, false)

	go func() {
		time.Sleep(time.Second)
		for {
			<-time.After(time.Second)
			app.QueueUpdateDraw(func() {
				s.update()
			})
		}
	}()

	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (s *state) update() {
	for _, f := range s.updateFunc {
		f()
	}

}

// func (s *state) makeLogData() *tview.TextView {
// 	textView := tview.NewTextView()
// 	textView.Set
// 	return textView
// }

func (s *state) makeEndpointForm(data *data) *tview.Form {
	f := tview.NewForm()
	nameInput := tview.NewInputField().
		SetLabel("Name").
		SetText(data.Name).
		SetChangedFunc(func(text string) {
			data.Name = text
		})

	f.AddFormItem(nameInput)

	s.updateFunc = append(s.updateFunc, func() {
		nameInput.SetText(data.Name)
	})
	return f
}

func (s *state) makeEndpointControl(text string, ep *endpoint.Endpoint) tview.Primitive {
	statusTextView := tview.NewTextView().SetTextAlign(1).SetText("Stopped")

	controlButton := tview.NewButton("Update")
	controlButton.SetDisabled(true)
	controlButton.SetBorder(true)
	controlButton.SetSelectedFunc(func() {
		ep.ClientUpdate()
		s.update()
	})

	startButton := tview.NewButton("Start")
	startButton.SetBorder(true)

	startButton.SetSelectedFunc(func() {
		if ep.Running() {
			ep.Stop()
			ep.Wait()
		} else {
			ep.Run(false)
			c := 0
			for !ep.Running() {
				time.Sleep(100 * time.Millisecond)
				c++
				if c > 10 {
					break
				}
			}
		}
		s.update()
	})

	s.updateFunc = append(s.updateFunc, func() {
		if ep.Running() {
			startButton.SetLabel("Stop")
			if ep.IsServer() {
				statusTextView.SetText("Running as Server")
			} else {
				controlButton.SetDisabled(false)
				statusTextView.SetText("Running as Client")
			}
		} else {
			controlButton.SetDisabled(true)
			startButton.SetLabel("Start")
			statusTextView.SetText("Stopped")
		}
		controlButton.Blur()
		startButton.Blur()
	})

	fv := tview.NewFlex()
	fv.SetTitle(text)
	fv.SetBorder(true)
	fv.AddItem(startButton, 0, 1, false)
	fv.AddItem(controlButton, 0, 1, false)
	fv.AddItem(statusTextView, 0, 1, false)
	return fv
}
