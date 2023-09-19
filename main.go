// Copyright 2017 The GoReporter Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// GoReporter is a Golang tool that does static analysis, unit testing, code
// review and generate code quality report.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/facebookgo/inject"
	"github.com/hellobchain/goreporter/engine"
	"github.com/hellobchain/goreporter/engine/processbar"
	"github.com/hellobchain/wswlog/wlogging"
)

var logger = wlogging.MustGetLoggerWithoutName()

// Received parameters, you can control some features using:
//
// -p:Specify the relative path of your project(Must Be Relative path),
//    by default, the current path is used
// -r:Specifies the save path for the generated report,
//    by default, the current path is used
// -e:Ignored detection of packages and multiple packages separated by commas.
// -t:Customize the path of the report template, not necessarily using the
//    default report template
// -f:Set the format to generate reports, support text, html and json,not
//    necessarily using the default formate-html.

const VERSION = "v3.0.0"

var (
	version        = flag.Bool("version", false, "print GoReporter version.")
	projectPath    = flag.String("p", "", "path of project.")
	reportPath     = flag.String("r", "", "path of report.")
	exceptPackages = flag.String("e", "", "except packages.")
	templatePath   = flag.String("t", "", "report html template path.")
	reportFormat   = flag.String("f", "", "project report format(text/json/html).")
	name           = flag.String("n", "", "project name.")
	coresOfCPU     = flag.Int("c", -1, "cores of CPU.")
)

func main() {
	flag.Parse()
	if *coresOfCPU != -1 && *coresOfCPU <= runtime.NumCPU() {
		runtime.GOMAXPROCS(*coresOfCPU)
	}
	if *version {
		logger.Infof("GoReporter %s\r\n", VERSION)
		os.Exit(0)
	}

	if *projectPath == "" {
		logger.Fatal("The project path is not specified")
	} else {
		_, err := os.Stat(*projectPath)
		if err != nil {
			logger.Fatal("project path is invalid")
		}
	}

	var templateHtml string
	if *templatePath == "" {
		templateHtml = engine.DefaultTpl
		logger.Info("The template path is not specified,and will use the default template")
	} else {
		if !strings.HasSuffix(*templatePath, ".html") {
			logger.Info("The template file is not a html template")
		}
		fileData, err := ioutil.ReadFile(*templatePath)
		if err != nil {
			logger.Fatal(err)
		} else {
			templateHtml = string(fileData)
		}
	}

	if *reportPath == "" {
		logger.Info("The report path is not specified, and the current path is used by default")
	} else {
		_, err := os.Stat(*reportPath)
		if err != nil {
			logger.Fatal("report path is invalid:", err)
		}
	}

	if *exceptPackages == "" {
		logger.Info("There are no packages that are excepted, review all items of the package")
	}

	synchronizer := &engine.Synchronizer{
		LintersProcessChans:   make(chan int64, 20),
		LintersFinishedSignal: make(chan string, 10),
	}
	syncRW := &sync.RWMutex{}
	waitGW := &engine.WaitGroupWrapper{}

	reporter := engine.NewReporter(*projectPath, *reportPath, *reportFormat, templateHtml, *name)
	strategyCountCode := &engine.StrategyCountCode{}
	strategyCyclo := &engine.StrategyCyclo{}
	strategyDeadCode := &engine.StrategyDeadCode{}
	strategyDependGraph := &engine.StrategyDependGraph{}
	strategyDepth := &engine.StrategyDepth{}
	strategyImportPackages := &engine.StrategyImportPackages{}
	strategyInterfacer := &engine.StrategyInterfacer{}
	strategySimpleCode := &engine.StrategySimpleCode{}
	strategySpellCheck := &engine.StrategySpellCheck{}
	strategyUnitTest := &engine.StrategyUnitTest{}
	strategyLint := &engine.StrategyLint{}
	strategyGoVet := &engine.StrategyGoVet{}
	strategyGoFmt := &engine.StrategyGoFmt{}

	if err := inject.Populate(
		reporter,
		synchronizer,
		strategyCountCode,
		strategyCyclo,
		strategyDeadCode,
		strategyDependGraph,
		strategyDepth,
		strategyImportPackages,
		strategyInterfacer,
		strategySimpleCode,
		strategySpellCheck,
		strategyUnitTest,
		strategyLint,
		strategyGoVet,
		strategyGoFmt,
		syncRW,
		waitGW,
	); err != nil {
		logger.Fatal(err)
	}

	reporter.AddLinters(strategyCountCode, strategyCyclo, strategyDeadCode, strategyDependGraph,
		strategyDepth, strategyImportPackages, strategyInterfacer, strategySimpleCode,
		strategySpellCheck, strategyUnitTest, strategyLint, strategyGoVet, strategyGoFmt)

	go processbar.LinterProcessBar(synchronizer.LintersProcessChans, synchronizer.LintersFinishedSignal)

	if err := reporter.Report(); err != nil {
		logger.Fatal(err)
	}

	if err := reporter.Render(); err != nil {
		logger.Fatal(err)
	}

	logger.Infof(fmt.Sprintf("GoReporter Finished,time consuming %vs", time.Since(reporter.StartTime).Seconds()))
}
