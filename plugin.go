package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"unsafe"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/selectel/fluent-bit-cloudwatch-input-plugin/internal/infra/client"
	"github.com/selectel/fluent-bit-cloudwatch-input-plugin/internal/infra/storage"
	"github.com/selectel/fluent-bit-cloudwatch-input-plugin/internal/infra/storage/sqlite"
	"github.com/selectel/fluent-bit-cloudwatch-input-plugin/internal/model"
	"github.com/selectel/fluent-bit-cloudwatch-input-plugin/internal/plugin"
	"github.com/selectel/fluent-bit-cloudwatch-input-plugin/pkg/fluentbit"
)

type Plugins struct {
	plugins []*plugin.Plugin
	mu      sync.RWMutex
}

func (p *Plugins) Push(plugin *plugin.Plugin) int {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.plugins = append(p.plugins, plugin)

	return len(p.plugins) - 1
}

func (p *Plugins) Pull(idx int) *plugin.Plugin {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.plugins[idx]
}

var plugins = Plugins{
	plugins: make([]*plugin.Plugin, 0),
	mu:      sync.RWMutex{},
}

//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	return fluentbit.FLBPluginRegister(def, "cloudwatch-input", "CloudWatch input plugin")
}

//export FLBPluginInit
func FLBPluginInit(flbPlugin unsafe.Pointer) int {
	region := fluentbit.FLBPluginConfigKey(flbPlugin, "region")
	endpoint := fluentbit.FLBPluginConfigKey(flbPlugin, "endpoint")
	logGroupName := fluentbit.FLBPluginConfigKey(flbPlugin, "log_group_name")
	logStreamName := fluentbit.FLBPluginConfigKey(flbPlugin, "log_stream_name")
	sqlitePath := fluentbit.FLBPluginConfigKey(flbPlugin, "sqlite_path")

	db, err := storage.NewSQLite(sqlitePath)
	if err != nil {
		log.Printf("failed to init sqlite storage: %v\n", err)

		return fluentbit.FLB_ERROR
	}

	ctx := context.Background()
	state := sqlite.NewState(db)

	cloudwatch, err := client.NewCloudwatchClient(ctx, region, endpoint)
	if err != nil {
		log.Printf("failed to init cloudwatch logs client: %v\n", err)

		return fluentbit.FLB_ERROR
	}

	pluginInput := plugin.NewPlugin(region, endpoint, logGroupName, logStreamName, cloudwatch, state)
	pluginIdx := plugins.Push(pluginInput)

	fluentbit.FLBPluginSetContext(flbPlugin, pluginIdx)

	return fluentbit.FLB_OK
}

//export FLBPluginInputCallbackCtx
func FLBPluginInputCallbackCtx(flbContext unsafe.Pointer, data *unsafe.Pointer, size *C.size_t) int {
	pluginIdx, ok := fluentbit.FLBPluginGetContext(flbContext).(int)
	if !ok {
		fmt.Println("failed to convert remote_context from fluent-bit to int")

		return fluentbit.FLB_ERROR
	}

	pluginInput := plugins.Pull(pluginIdx)

	ctx := context.Background()

	nextToken, err := pluginInput.GetNextToken(ctx)
	if err != nil {
		fmt.Printf("failed to get next token from persistent storage: %v\n", err)

		return fluentbit.FLB_ERROR
	}

	events, nextToken, err := pluginInput.GetLogEvents(ctx, nextToken)
	if err != nil {
		fmt.Printf("failed to get log events from cloudwatch: %v\n", err)

		return fluentbit.FLB_RETRY
	}

	if len(events) == 0 {
		if nextToken != "" {
			err = pluginInput.SetNextToken(ctx, nextToken)
			if err != nil {
				fmt.Printf("failed to save next token to persistent storage: %v\n", err)

				return fluentbit.FLB_ERROR
			}
		}

		return fluentbit.FLB_OK
	}

	entry := []interface{}{
		time.Now().UTC().Unix(),
		map[string][]model.Event{"events": events},
	}

	packed, err := msgpack.Marshal(&entry)
	if err != nil {
		fmt.Printf("failed to marshal entry to msgpack: %v\n", err)

		return fluentbit.FLB_ERROR
	}

	length := len(packed)
	*data = C.CBytes(packed)
	*size = C.size_t(length)

	err = pluginInput.SetNextToken(ctx, nextToken)
	if err != nil {
		fmt.Printf("failed to save next token to persistent storage: %v\n", err)

		return fluentbit.FLB_ERROR
	}

	return fluentbit.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	return fluentbit.FLB_OK
}

func main() {
}
