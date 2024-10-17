package main

import (
	"distapp/internal/app"
	"distapp/internal/router"
	"distapp/internal/store"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/serf/serf"
	"github.com/pkg/errors"
)

func main() {
	cluster, err := setupCluster(
		os.Getenv("ADVERTISE_ADDR"),
		os.Getenv("CLUSTER_ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer cluster.Leave()
	storage := store.NewInMemory()
	i := app.NewInstance(storage, cluster)
	r := router.NewRouter(i)
	// go func() {
	// 	if err := http.ListenAndServe(":8080", r); err != nil {
	// 		log.Println(err)
	// 	}
	// }()

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Println(err)
	}
	// ctx := context.Background()
	// if name, err := os.Hostname(); err == nil {
	// 	ctx = context.WithValue(ctx, "name", name)
	// }

	// debugDataPrinterTicker := time.Tick(time.Second * 5)
	// for {
	// 	select {
	// 	case <-debugDataPrinterTicker:
	// 		log.Printf("Members: %v\n", cluster.Members())

	// 		/*curVal, curGen := storage.GetValue("id")
	// 		log.Printf("State: Val: %v Gen: %v\n", curVal, curGen)*/
	// 	}
	// }
}

func setupCluster(advertiseAddr string, clusterAddr string) (*serf.Serf, error) {
	conf := serf.DefaultConfig()
	conf.Init()
	conf.MemberlistConfig.AdvertiseAddr = advertiseAddr

	cluster, err := serf.Create(conf)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't create cluster")
	}

	_, err = cluster.Join([]string{clusterAddr}, true)
	if err != nil {
		log.Printf("Couldn't join cluster, starting own: %v\n", err)
	}

	return cluster, nil
}
