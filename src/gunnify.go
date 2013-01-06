package main

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/spaghetty/udev"
)

var USB_VENDOR_ID_LOGITECH string = "046d"
var USB_DEVICE_ID_UNIFYING_RECEIVER string = "c52b"
var USB_DEVICE_ID_UNIFYING_RECEIVER_2 string = "c532"

var magic_sequence = []byte{0x10, 0xFF, 0x80, 0xB2, 0x01, 0x00, 0x00}

type Gui struct {
	MainWindow *gtk.Window
	SplashWindow *gtk.Window
	Store *gtk.ListStore
	Status *gtk.Statusbar
	Start *gtk.Button
	Recheck *gtk.Button
}

var MainGui Gui;

func SearchValid() []udev.Device {
	u := udev.NewUdev()
        defer u.Unref()

        e := u.NewEnumerate()
        defer e.Unref()

        e.AddMatchSubsystem("hidraw")
        err := e.ScanDevices()
	if err!=nil {
		fmt.Println(err)
	}
	devices := make(map[udev.DevNum]udev.Device)
	
        for device := e.First(); !device.IsNil(); device = device.Next() {
                path := device.Name()
		//fmt.Println(path)
                dev := u.DeviceFromSysPath(path)
		dev = dev.ParentWithSubsystemDevType("usb", "usb_device")
		if ( dev.SysAttrValue("idVendor")==USB_VENDOR_ID_LOGITECH && 
			( dev.SysAttrValue("idProduct")==USB_DEVICE_ID_UNIFYING_RECEIVER || 
			dev.SysAttrValue("idProduct")==USB_DEVICE_ID_UNIFYING_RECEIVER_2)) {
			fmt.Println(dev.DevType(), dev.DevNum(), dev.SysAttrValue("product"))
			if _, ok :=devices[dev.DevNum()]; !ok {
				devices[dev.DevNum()] = dev
				fmt.Println(dev.SeqNum(), dev.DevNode())
			}
		}
		
	}
	var dlist []udev.Device;
	for _, v := range devices {
		dlist = append(dlist, v)
	}
	time.Sleep(1 * time.Millisecond)
	return dlist
}

func doWrite(d string) bool {
	f, err := os.OpenFile(d, os.O_WRONLY|os.O_SYNC,500)
	if err!=nil {
		MainGui.notify("Some error opening the file")
		return false
	}
	defer f.Close()
	_, err = f.Write(magic_sequence)
	if err!=nil {
		MainGui.notify("Some error writing to the file, "+ err.Error())
		return false
	}
	MainGui.notify("Ready to sync try to turn off and then on the desidered device")
	return true
}

func (g *Gui)buildList(vbox *gtk.VBox) {
	frame := gtk.NewFrame("Device List")
	framebox := gtk.NewVBox(false, 1)
        frame.Add(framebox)
	vbox.Add(frame)
	g.Status = gtk.NewStatusbar()
	vbox.PackStart(g.Status,false,false,0)
	g.Store = gtk.NewListStore(glib.G_TYPE_STRING, glib.G_TYPE_STRING)
	treeview := gtk.NewTreeView()
	framebox.Add(treeview)
	treeview.SetModel(g.Store)
	treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Device", gtk.NewCellRendererText(), "text", 0))
	treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Name", gtk.NewCellRendererText(), "text", 1))
	treeview.GetSelection().SetMode(gtk.SELECTION_SINGLE)
	controls := gtk.NewHBox(true,0)
	g.Start = gtk.NewButtonWithLabel("Start Sync")
	g.Start.Clicked(func() {
		var iter gtk.TreeIter
		var device glib.GValue
		selection := treeview.GetSelection()
		if selection.CountSelectedRows() > 0 {
			selection.GetSelected(&iter)
			g.Store.GetValue(&iter, 0, &device)
			MainGui.notify("Start Writing On: "+device.GetString())
			doWrite(device.GetString())
		} else {
			MainGui.notify("No Active Selection")
		}
	})
	controls.Add(g.Start)
	g.Recheck = gtk.NewButtonWithLabel("Rescan")
	g.Recheck.Clicked(func() {
		devices := SearchValid()
		MainGui.Store.Clear()
		for _, x := range devices {
			MainGui.appendItem("/dev/hidraw"+strconv.FormatUint(x.SeqNum(),10), x.SysAttrValue("product"))
		}
		MainGui.notify("Scanning Done")
	})
	controls.Add(g.Recheck)
	framebox.PackStart(controls, false, false,0)
}

func (g *Gui)notify(s string) {
	g.Status.Push(0, s)
	time.AfterFunc(1000*time.Millisecond, func () {
		g.Status.Pop(0)
	})
}

func (g *Gui)buildGUI() {
	g.MainWindow = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	g.MainWindow.SetPosition(gtk.WIN_POS_CENTER)
	g.MainWindow.SetTitle("Gunnify")
	g.MainWindow.SetIconName("gtk-dialog-info")
	g.MainWindow.Connect("destroy", func(ctx *glib.CallbackContext) {
		println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "foo")
	g.MainWindow.SetSizeRequest(600, 300)
	vbox := gtk.NewVBox(false, 0)
	g.buildList(vbox)
	g.MainWindow.Add(vbox)
}

func (g *Gui)buildSplash() {
	g.SplashWindow = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	g.SplashWindow.SetTypeHint(gdk.WINDOW_TYPE_HINT_SPLASHSCREEN)
	g.SplashWindow.SetSizeRequest(600,300)
}

func (g *Gui)appendItem(name, descr string) {
	var iter gtk.TreeIter
	MainGui.Store.Append(&iter)
	MainGui.Store.Set(&iter, name, descr)
}

func unlinked_main() {
	devices := SearchValid()
	for _, x := range devices {
		MainGui.appendItem("/dev/hidraw"+strconv.FormatUint(x.SeqNum(),10), x.SysAttrValue("product"))
	}
	/////////// THIS IS JUST FOR DEBUG PURPOSE //////////////
	//MainGui.appendItem("prova1", "prova 123")
	//MainGui.appendItem("prova2", "Prova 123")
	//MainGui.appendItem("prova3", "prova 123")
	//MainGui.appendItem("prova4", "prova 1234")
	///////////////////////////////////////////////////////
	gdk.ThreadsEnter()
	MainGui.SplashWindow.Hide()
	//ctx := MainGui.Status.GetContextId("prova 123")
	MainGui.Status.Push(0, "ready for operate")
	MainGui.MainWindow.ShowAll()
	gdk.ThreadsLeave()
}

func main() {
	gtk.Init(nil)
	gdk.ThreadsInit()
	MainGui.buildSplash()
	MainGui.buildGUI()
	MainGui.SplashWindow.ShowAll()
	go unlinked_main()
	gtk.Main()
}