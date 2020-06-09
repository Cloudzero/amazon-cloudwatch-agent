package serialization

import (
	"fmt"

	"github.com/aws/amazon-cloudwatch-agent/tool/data"
	"github.com/aws/amazon-cloudwatch-agent/tool/processors"
	"github.com/aws/amazon-cloudwatch-agent/tool/processors/ssm"
	"github.com/aws/amazon-cloudwatch-agent/tool/runtime"
	"github.com/aws/amazon-cloudwatch-agent/tool/util"
)

var Processor processors.Processor = &processor{}

type processor struct{}

func (p *processor) Process(ctx *runtime.Context, config *data.Config) {
	_, resultMap := config.ToMap(ctx)
	byteArray := util.SerializeResultMapToJsonByteArray(resultMap)
	util.SaveResultByteArrayToJsonFile(byteArray)
	fmt.Printf("Current config as follows:\n"+
		"%s\n"+
		"Please check the above content of the config.\n"+
		"The config file is also located at %s.\n"+
		"Edit it manually if needed.\n",
		string(byteArray),
		util.ConfigFilePath())
}

func (p *processor) NextProcessor(ctx *runtime.Context, config *data.Config) interface{} {
	return ssm.Processor
}