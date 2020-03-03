/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"flag"
	"os"

	"github.com/spf13/pflag"

	"k8s.io/klog"

	"fmt"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd"
	"k8s.io/utils/trace"
	"time"
)

// Run creates and executes new kubeadm command
func Run() error {
	klog.InitFlags(nil)
	trace.InitFlags(nil)
	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	pflag.Set("logtostderr", "true")
	// We do not want these flags to show up in --help
	// These MarkHidden calls must be after the lines above
	pflag.CommandLine.MarkHidden("version")
	pflag.CommandLine.MarkHidden("log-flush-frequency")
	pflag.CommandLine.MarkHidden("alsologtostderr")
	pflag.CommandLine.MarkHidden("log-backtrace-at")
	pflag.CommandLine.MarkHidden("log-dir")
	pflag.CommandLine.MarkHidden("logtostderr")
	pflag.CommandLine.MarkHidden("stderrthreshold")
	pflag.CommandLine.MarkHidden("vmodule")

	cmd := cmd.NewKubeadmCommand(os.Stdin, os.Stdout, os.Stderr)
	ret := cmd.Execute()

	initTrace := trace.New("New trace", trace.Field{"name", os.Args[0]})
	fmt.Println("Sleep 1 Second")
	time.Sleep(time.Second)
	initTrace.Step("trace sleep 1")
	initTrace.Log()
	return ret
}
