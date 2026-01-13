package traceroute

import (
	"github.com/kluzzebass/reflag/preprocessor"
	"github.com/kluzzebass/reflag/preprocessor/urlparse"
)

func init() {
	preprocessor.Register(&Preprocessor{})
}

type Preprocessor struct{}

func (p *Preprocessor) ToolName() string {
	return "traceroute"
}

func (p *Preprocessor) Description() string {
	return "Extract hostname from URLs"
}

func (p *Preprocessor) Preprocess(args []string) []string {
	return urlparse.ProcessArgs(args)
}
