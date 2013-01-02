package main

import (
	"fmt"
	"time"
	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/spaghetty/udev"
)

var USB_VENDOR_ID_LOGITECH string = "046d"
var USB_DEVICE_ID_UNIFYING_RECEIVER string = "c52b"
var USB_DEVICE_ID_UNIFYING_RECEIVER_2 string = "c532"

var magic_sequence = [...]uint8{0x10, 0xFF, 0x80, 0xB2, 0x01, 0x00, 0x00}

var MainWindow *gtk.Window
var SplashWindow *gtk.Window
var Store *gtk.ListStore

func SearchValid() []udev.Device {
	u := udev.NewUdev()
        defer u.Unref()
	
	fmt.Println("cool")
        e := u.NewEnumerate()
        defer e.Unref()
	fmt.Println("cool1")

        e.AddMatchSubsystem("hidraw")
        err := e.ScanDevices()
	if err!=nil {
		fmt.Println(err)
	}
	devices := make(map[udev.DevNum]udev.Device)
	
        for device := e.First(); !device.IsNil(); device = device.Next() {
                path := device.Name()
                dev := u.DeviceFromSysPath(path)
		dev = dev.ParentWithSubsystemDevType("usb", "usb_device")
		if ( dev.SysAttrValue("idVendor")==USB_VENDOR_ID_LOGITECH && 
			( dev.SysAttrValue("idProduct")==USB_DEVICE_ID_UNIFYING_RECEIVER || 
			dev.SysAttrValue("idProduct")==USB_DEVICE_ID_UNIFYING_RECEIVER_2)) {
			fmt.Println(dev.DevType(), dev.DevNum(), dev.SysAttrValue("product"))
			if _, ok :=devices[dev.DevNum()]; !ok {
				devices[dev.DevNum()] = dev
			}
		}
		
	}
	var dlist []udev.Device;
	for _, v := range devices {
		dlist = append(dlist, v)
	}
	time.Sleep(1000 * time.Millisecond)
	return dlist
}

func unlinked_main() {
	devices := SearchValid()
	fmt.Println("ciao")
	gdk.ThreadsEnter()
	for _, x := range devices {
		var iter gtk.TreeIter
		Store.Append(&iter)
		Store.Set(&iter, x.DevNode(), x.SysAttrValue("product"))
		//tmp := gtk.NewRadioButtonWithLabel(bgroup, x.DevNode())
		//bgroup = tmp.GetGroup()
		//bbbox.Add(tmp)
	}
	SplashWindow.Hide()
	MainWindow.ShowAll()
	gdk.ThreadsLeave()
}

func buildList(vbox *gtk.VBox) {
	frame := gtk.NewFrame("Device List")
	framebox := gtk.NewVBox(false, 1)
        frame.Add(framebox)
	vbox.Add(frame)	
	Store = gtk.NewListStore(glib.G_TYPE_STRING, glib.G_TYPE_STRING)
	treeview := gtk.NewTreeView()
	framebox.Add(treeview)
	treeview.SetModel(Store)
	treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Device", gtk.NewCellRendererText(), "text", 0))
	treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Name", gtk.NewCellRendererText(), "text", 1))
	controls := gtk.NewHBox(false,2)
	start := gtk.NewButtonWithLabel("Start Sync")
	start.Clicked(func() {
		var iter gtk.TreeIter
		var device glib.GValue
		treeview.GetSelection().GetSelected(&iter)
		Store.GetValue(&iter, 0, &device)
		fmt.Println(device.GetString())
	})
	controls.Add(start)
	framebox.Add(controls)
}

func buildGUI() {
	MainWindow = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	MainWindow.SetPosition(gtk.WIN_POS_CENTER)
	MainWindow.SetTitle("GTK Go!")
	MainWindow.SetIconName("gtk-dialog-info")
	MainWindow.Connect("destroy", func(ctx *glib.CallbackContext) {
		println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "foo")
	MainWindow.SetSizeRequest(600, 300)
	vbox := gtk.NewVBox(false, 10)
	buildList(vbox)
	MainWindow.Add(vbox)
}

func buildSplash() {
	SplashWindow = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	SplashWindow.SetTypeHint(gdk.WINDOW_TYPE_HINT_SPLASHSCREEN)
	SplashWindow.SetSizeRequest(600,300)
}

func main() {
	gtk.Init(nil)
	gdk.ThreadsInit()
	buildSplash()
	buildGUI()
	SplashWindow.ShowAll()
	go unlinked_main()
	gtk.Main()
}