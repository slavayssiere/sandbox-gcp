package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"
)

// User-provided constants.
const (
	columnFamilyName = "ms"
	columnNameData   = "data"
	columnNameUser   = "user"
	columnNameSource = "source"
	columnNameDate   = "time"
)

// sliceContains reports whether the provided string is present in the given slice of strings.
func sliceContains(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}

func bigtableClient(ctx context.Context) bigtable.Client {
	adminClient, err := bigtable.NewAdminClient(ctx, *projectid, *instanceid)
	if err != nil {
		log.Fatalf("Could not create admin client: %v", err)
	}

	tblInfo, err := adminClient.TableInfo(ctx, *tableid)
	if err != nil {
		log.Fatalf("Could not read info for table %s: %v", *tableid, err)
	}

	if !sliceContains(tblInfo.Families, columnFamilyName) {
		if err := adminClient.CreateColumnFamily(ctx, *tableid, columnFamilyName); err != nil {
			log.Fatalf("Could not create column family %s: %v", columnFamilyName, err)
		}
	}

	client, err := bigtable.NewClient(ctx, *projectid, *instanceid)
	if err != nil {
		log.Fatalf("Could not create data operations client: %v", err)
	}

	return *client
}

func (s server) writeMessage(ctx context.Context, mess libmetier.MessageSocial) {

	tbl := s.bt.Open(*tableid)

	// Mutation Way
	mut := bigtable.NewMutation()
	mut.Set(columnFamilyName, columnNameData, bigtable.Now(), []byte(mess.Data))
	mut.Set(columnFamilyName, columnNameUser, bigtable.Now(), []byte(mess.User))
	mut.Set(columnFamilyName, columnNameSource, bigtable.Now(), []byte(mess.Source))

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(mess.Date.Nanosecond()))
	mut.Set(columnFamilyName, columnNameDate, bigtable.Now(), b)

	// Read pubsub attribute key to determine BT row key
	var key = fmt.Sprintf("%s%s", mess.Source, mess.Date)

	if err := tbl.Apply(ctx, key, mut); err != nil {
		fmt.Println(err)
	}
}

func (s server) readMessage(ctx context.Context) libmetier.MessageSocial {

	var mess libmetier.MessageSocial

	tbl := s.bt.Open(*tableid)

	rowKey := "test"

	row, err := tbl.ReadRow(ctx, rowKey)
	if err != nil {
		log.Fatalf("Could not read row with key %s: %v", rowKey, err)
	}
	log.Printf("Row key: %s\n", rowKey)
	mess.Data = string(row[columnFamilyName][0].Value)
	mess.User = string(row[columnFamilyName][1].Value)
	mess.Source = string(row[columnFamilyName][2].Value)
	mess.Date = time.Unix(0, int64(binary.LittleEndian.Uint64(row[columnFamilyName][3].Value)))
	log.Println("Data:", mess.Data)
	log.Println("Source:", mess.Source)
	log.Println("User:", mess.User)
	log.Println("Date:", mess.Date)

	return mess
}

func (s server) writeMessages(ctx context.Context) {

	for {
		mess := <-s.messages
		s.writeMessage(ctx, mess)
	}
}
