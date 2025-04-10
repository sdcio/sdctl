/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	dsutils "github.com/sdcio/data-server/pkg/utils"
	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
	"github.com/sdcio/sdctl/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	transactionId string
)

// dataSetIntentCmd represents the set-intent command
var dataTransactionCmd = &cobra.Command{
	Use:          "transaction",
	Short:        "Transaction based commands",
	SilenceUsage: true,
}

func init() {
	dataCmd.AddCommand(dataTransactionCmd)
	dataTransactionCmd.PersistentFlags().StringVarP(&transactionId, "transaction-id", "", "", "Transaction ID")
	dataTransactionCmd.MarkPersistentFlagRequired("transaction-id")
	dataTransactionCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if transactionId == "" {
			return fmt.Errorf("missing parameter transaction-id")
		}
		if datastoreName == "" {
			return fmt.Errorf("missing parameter datastore")
		}
		return nil
	}
}

type TransactionIntent struct {
	// intent name
	IntentName string `json:"intent-name"`
	// intent priority
	Priority int32 `json:"priority"`
	// list of updates
	Updates []updateDef `json:"updates"`
	// delete indicator
	Delete bool `json:"delete"`
	// delete only from intended store
	// basically keeping the config on the device but unmanaged
	OnlyIntended bool `json:"only-intended"`
}

func LoadTransactionIntent(file string) (*TransactionIntent, error) {
	result := &TransactionIntent{}
	err := utils.JsonUnmarshalStrict(file, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *TransactionIntent) ToSdcpbTransactionIntent() (*sdcpb.TransactionIntent, error) {
	result := &sdcpb.TransactionIntent{
		Intent:   t.IntentName,
		Priority: t.Priority,
		Delete:   t.Delete,
		Orphan:   t.OnlyIntended,
	}

	for _, upd := range t.Updates {
		p, err := dsutils.ParsePath(upd.Path)
		if err != nil {
			return nil, err
		}
		bb, err := json.Marshal(upd.Value)
		if err != nil {
			return nil, err
		}
		result.Update = append(result.Update, &sdcpb.Update{
			Path: p,
			Value: &sdcpb.TypedValue{
				Value: &sdcpb.TypedValue_JsonVal{JsonVal: bb},
			},
		})
	}

	return result, nil
}

type updateDef struct {
	Path  string `json:"path,omitempty"`
	Value any    `json:"value,omitempty"`
}
