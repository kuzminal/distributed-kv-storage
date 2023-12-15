package app

import (
	"bytes"
	"context"
	"distapp/internal/store"
	"distapp/models"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/serf/serf"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"math/rand"
	"net/http"
)

type Instance struct {
	storage store.Store
	cluster *serf.Serf
}

func NewInstance(store store.Store, cluster *serf.Serf) *Instance {
	return &Instance{storage: store, cluster: cluster}
}

const MembersToNotify = 2

func (i *Instance) notifyOthers(ctx context.Context, st store.Store, key string) {
	g, ctx := errgroup.WithContext(ctx)
	otherMembers := getOtherMembers(i.cluster)
	if len(otherMembers) <= 2 {
		for _, member := range otherMembers {
			curMember := member
			g.Go(func() error {
				log.Printf("addr is : %v", curMember.Addr.String())
				return notifyMember(ctx, curMember.Addr.String(), st, key)
			})
		}
	} else {
		randIndex := rand.Int() % len(otherMembers)
		for i := 0; i < MembersToNotify; i++ {
			curIndex := i
			g.Go(func() error {
				return notifyMember(
					ctx,
					otherMembers[(randIndex+curIndex)%len(otherMembers)].Addr.String(),
					st,
					key)
			})
		}
	}

	err := g.Wait()
	if err != nil {
		log.Printf("Error when notifying other members: %v", err)
	}
}

func notifyMember(ctx context.Context, addr string, st store.Store, key string) error {
	val, gen := st.GetValue(key)
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(models.NotifyRequest{Gen: gen, Value: val, Key: key})

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%v:8080/notify?notifier=%v", addr, ctx.Value("name")), &buf)
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "Couldn't create request")
	}
	req = req.WithContext(ctx)

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "Couldn't make request")
	}
	return nil
}

func getOtherMembers(cluster *serf.Serf) []serf.Member {
	members := cluster.Members()
	for i := 0; i < len(members); {
		if members[i].Name == cluster.LocalMember().Name || members[i].Status != serf.StatusAlive {
			if i < len(members)-1 {
				members = append(members[:i], members[i+1:]...)
			} else {
				members = members[:i]
			}
		} else {
			i++
		}
	}
	return members
}
