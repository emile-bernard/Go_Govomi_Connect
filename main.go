package main

// import "fmt"

// func main() {
// 	name:= "VS Code"
// 	fmt.Println("Welcome to", name)
// }

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"text/tabwriter"

	"reflect"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/event"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func handleEvent(ref types.ManagedObjectReference, events []types.BaseEvent) (err error) {
	for _, event := range events {
		eventType := reflect.TypeOf(event).String()
		fmt.Printf("Event found of type %s\n", eventType)
	}

	return nil
}

func main() {
	unsecure := true

	// Creating a connection context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parsing URL
	govomihost := os.Getenv("GOVOMIHOST")
	govomiuid := os.Getenv("GOVOMIUID")
	govomipwd := os.Getenv("GOVOMIPWD")
	fmt.Println("host:", govomihost)
	fmt.Println("uid:", govomiuid)

	url, err := url.Parse("https://" + govomiuid + ":" + govomipwd + "@" + govomihost)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	// Connecting to vCenter
	client, err := govmomi.NewClient(ctx, url, unsecure)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	// Selecting default datacenter
	finder := find.NewFinder(client.Client, unsecure)
	dc, err := finder.DefaultDatacenter(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	refs := []types.ManagedObjectReference{dc.Reference()}

	// Setting up the event manager
	eventManager := event.NewManager(client.Client)
	err = eventManager.Events(ctx, refs, 10, false, false, handleEvent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	// vCenter version
	info := client.ServiceContent.About

	fmt.Printf("Info %s\n", info)
	fmt.Printf("Connected to vCenter version %s\n", info.Version)

	// Create a view of Datastore objects
	m := view.NewManager(client.Client)

	v, err := m.CreateContainerView(ctx, client.ServiceContent.RootFolder, []string{"Datastore"}, true)
	if err != nil {
		log.Fatal(err)
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all datastores
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.Datastore.html
	var dss []mo.Datastore
	err = v.Retrieve(ctx, []string{"Datastore"}, []string{"summary"}, &dss)
	if err != nil {
		log.Fatal(err)
	}

	// Print summary per datastore (see also: govc/datastore/info.go)

	tw := tabwriter.NewWriter(os.Stdout, 2, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Name:\tType:\tCapacity:\tFree:\n")

	for _, ds := range dss {
		fmt.Fprintf(tw, "%s\t", ds.Summary.Name)
		fmt.Fprintf(tw, "%s\t", ds.Summary.Type)
		fmt.Fprintf(tw, "%s\t", units.ByteSize(ds.Summary.Capacity))
		fmt.Fprintf(tw, "%s\t", units.ByteSize(ds.Summary.FreeSpace))
		fmt.Fprintf(tw, "\n")
	}

	_ = tw.Flush()
}
