package main

import (
	"fmt"
	//	"image/color"
	//	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/wilrad/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	//	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	a = app.New()
	w = a.NewWindow("Flughafenverwaltung")
	//	bindings = make([]ExternalString)
	db, err = sql.Open("sqlite3", "/home/wilrad/sqlite/flug.db")
	flüge   = selectFlüge()

	header  = []string{"flugid", "flugnr", "von", "nach", "fluglinie", "flugzeug"}
	fields  = util.Array2string(header)
	ohneKey = util.Array2string(header[1:])
)

type flug struct {
	flugid    string
	flugnr    string
	von       string
	nach      string
	fluglinie string
	flugzeug  string
}

func init() {
	util.CheckErr(err)

}

func makeUI() fyne.CanvasObject {
	var flug flug
	ind := 0

	flug = flüge[ind]
	bFlugId := binding.BindString(&flug.flugid)
	bFlugnr := binding.BindString(&flug.flugnr)
	bVon := binding.BindString(&flug.von)
	bNach := binding.BindString(&flug.nach)
	bFluglinie := binding.BindString(&flug.fluglinie)
	bFlugzeug := binding.BindString(&flug.flugzeug)

	reload := func() {
		err = bFlugId.Reload()
		util.CheckErr(err)
		err = bFlugnr.Reload()
		util.CheckErr(err)
		err = bVon.Reload()
		util.CheckErr(err)
		err = bNach.Reload()
		util.CheckErr(err)
		err = bFluglinie.Reload()
		util.CheckErr(err)
		err = bFlugzeug.Reload()
		util.CheckErr(err)
	}
	wFlugId := widget.NewEntryWithData(bFlugId)
	wFlugId.Validator = validation.NewRegexp(`^[1-9]?[0-9]{0,6}$`, "nur max 6 Ziffern")
	wFlugnr := widget.NewEntryWithData(bFlugnr)
	wFlugnr.Validator = validation.NewRegexp(`^[A-Z]{2}[0-9]{4}$`, "Form AZ9999")
	wVon := widget.NewEntryWithData(bVon)
	wVon.Validator = validation.NewRegexp(`^[1-9]?[0-9]{0,6}$`, "nur max 6 Ziffern")
	wNach := widget.NewEntryWithData(bNach)
	wNach.Validator = validation.NewRegexp(`^[1-9]?[0-9]{0,6}$`, "nur max 6 Ziffern")
	wFluglinie := widget.NewEntryWithData(bFluglinie)
	wFlugzeug := widget.NewEntryWithData(bFlugzeug)
	disable := func() {
		wFlugId.Disable()
		wFlugnr.Disable()
		wVon.Disable()
		wNach.Disable()
		wFlugzeug.Disable()
		wFluglinie.Disable()
	}
	enable := func() {
		wFlugId.Enable()
		wFlugnr.Enable()
		wVon.Enable()
		wNach.Enable()
		wFlugzeug.Enable()
		wFluglinie.Enable()
	}
	disable()

	var (
		form = widget.NewForm(
			widget.NewFormItem("FlugId", wFlugId),
			widget.NewFormItem("FlugNr", wFlugnr),
			widget.NewFormItem("von", wVon),
			widget.NewFormItem("nach", wNach),
			widget.NewFormItem("Fluglinie", wFluglinie),
			widget.NewFormItem("Flugzeug", wFlugzeug),
			widget.NewFormItem("", widget.NewButton("first", func() { ind = 0; flug = flüge[ind]; reload() })),
			widget.NewFormItem("", widget.NewButton("next", func() {
				if ind < len(flüge)-1 {
					ind += 1
				}
				flug = flüge[ind]
				reload()
			})),
			widget.NewFormItem("", widget.NewButton("prev", func() {
				if ind > 0 {
					ind -= 1
				}
				flug = flüge[ind]
				reload()
			})),
			widget.NewFormItem("", widget.NewButton("last", func() { ind = len(flüge) - 1; flug = flüge[ind]; reload() })),
			widget.NewFormItem("", widget.NewButton("Edit", func() { enable() })),
			widget.NewFormItem("", widget.NewButton("Update", func() { flüge[ind] = flug; disable() })),
			widget.NewFormItem("", widget.NewButton("End", func() { w.Close() })),
		)
	)
	//	form.OnCancel = func() { fmt.Println("Cancelled") }
	//	form.OnSubmit = func() { fmt.Println("Submitted"); w.Close() }

	top := container.New(layout.NewHBoxLayout())
	left := container.New(layout.NewVBoxLayout(), form)
	bottom := container.New(layout.NewHBoxLayout(),
		widget.NewButton("delete", func() { flüge = append(flüge[:ind], flüge[ind+1:]...); flug = flüge[ind]; reload() }))
	fmt.Println("%T", bottom)
	return container.New(layout.NewBorderLayout(top, bottom, left, nil), top, left, bottom)
}

func selectFlüge() []flug {
	q := "select " + fields + " from flug"
	rows, err := db.Query(q)
	util.CheckErr(err)

	var flüge []flug
	for i := 0; rows.Next(); i++ {
		var f flug
		err := rows.Scan(&f.flugid, &f.flugnr, &f.von, &f.nach, &f.fluglinie, &f.flugzeug)
		util.CheckErr(err)

		flüge = append(flüge, f)
	}
	//	fmt.Println(flüge)
	return flüge
}

func main() {

	w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("File",
		fyne.NewMenuItem("Open Directory...", func() {
			fmt.Println("menu")
		}))))

	a.Settings().SetTheme(theme.LightTheme())
	w.Resize(fyne.NewSize(800, 600))
	w.SetContent(makeUI())
	w.ShowAndRun()

	err = db.Close()
	util.CheckErr(err)
}
