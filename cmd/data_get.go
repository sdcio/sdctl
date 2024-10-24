/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sdcio/schema-server/pkg/utils"
	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/prototext"
)

var paths []string
var dataType string
var intended bool
var encoding string

// dataGetCmd represents the get command
var dataGetCmd = &cobra.Command{
	Use:          "get",
	Short:        "get data",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if candidate != "" && intended {
			return fmt.Errorf("cannot set a candidate name and intended store at the same time")
		}
		var dt sdcpb.DataType
		switch strings.ToLower(dataType) {
		case "all":
		case "config":
			dt = sdcpb.DataType_CONFIG
		case "state":
			dt = sdcpb.DataType_STATE
		default:
			return fmt.Errorf("invalid flag value --type %s", dataType)
		}

		var enc sdcpb.Encoding
		switch strings.ToLower(encoding) {
		case "string":
			enc = sdcpb.Encoding_STRING
		case "json":
			enc = sdcpb.Encoding_JSON
		case "json_ietf":
			enc = sdcpb.Encoding_JSON_IETF
		case "proto":
			enc = sdcpb.Encoding_PROTO
		}

		req := &sdcpb.GetDataRequest{
			Name:     datastoreName,
			DataType: dt,
			Encoding: enc,
		}
		for _, p := range paths {
			xp, err := utils.ParsePath(p)
			if err != nil {
				return err
			}
			req.Path = append(req.Path, xp)
		}
		if candidate != "" {
			req.Datastore = &sdcpb.DataStore{
				Type: sdcpb.Type_CANDIDATE,
				Name: candidate,
			}
		}
		if intended {
			req.Datastore = &sdcpb.DataStore{
				Type:     sdcpb.Type_INTENDED,
				Owner:    owner,
				Priority: priority,
			}
		}
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		dataClient, err := createDataClient(ctx, addr)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, "request:")
		fmt.Fprintln(os.Stderr, prototext.Format(req))
		stream, err := dataClient.GetData(ctx, req)
		if err != nil {
			return err
		}
		count := 0
		for {
			rsp, err := stream.Recv()
			if err != nil {
				if strings.Contains(err.Error(), "EOF") {
					break
				}
				return err
			}
			count++

			switch strings.ToLower(format) {
			case "json":
				switch strings.ToLower(encoding) {
				case "json", "json_ietf":
					for _, notifications := range rsp.GetNotification() {
						for _, upd := range notifications.Update {

							var val any
							err = json.Unmarshal(upd.GetValue().GetJsonVal(), &val)
							if err != nil {
								return err
							}

							b, err := json.MarshalIndent(val, "", "  ")
							if err != nil {
								return err
							}
							fmt.Println(string(b))
						}
					}
				default:
					b, err := json.MarshalIndent(rsp, "", "  ")
					if err != nil {
						return err
					}
					fmt.Println(string(b))
				}
			case "flat":
				for _, n := range rsp.GetNotification() {
					for _, upd := range n.GetUpdate() {
						p := utils.ToXPath(upd.GetPath(), false)
						// upd.GetValue()
						fmt.Printf("%s: %s\n", p, upd.GetValue())
					}
				}

			default:
				fmt.Println(prototext.Format(rsp))
			}

		}
		fmt.Fprintln(os.Stderr, "num notifications:", count)

		return nil
	},
}

func init() {
	dataCmd.AddCommand(dataGetCmd)
	dataGetCmd.Flags().StringArrayVarP(&paths, "path", "", []string{}, "get path(s)")
	dataGetCmd.Flags().StringVarP(&dataType, "type", "", "ALL", "data type, one of: ALL, CONFIG, STATE")
	dataGetCmd.Flags().StringVarP(&encoding, "encoding", "", "STRING", "encoding of the returned data: STRING, JSON, JSON_IETF or PROTO")

	// intended store
	dataGetCmd.Flags().BoolVarP(&intended, "intended", "", false, "get data from intended store")
	dataGetCmd.Flags().StringVarP(&owner, "owner", "", "", "intended store owner to query")
	dataGetCmd.Flags().Int32VarP(&priority, "priority", "", 0, "intended store priority")
}
