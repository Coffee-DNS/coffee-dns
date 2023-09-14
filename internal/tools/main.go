package main

import (
	_ "github.com/client9/misspell"
	_ "github.com/mgechev/revive"
	_ "github.com/securego/gosec/v2"
	_ "github.com/securego/gosec/v2/report/sarif" // required by gosec
	_ "github.com/securego/gosec/v2/report/text"  // required by gosec
	_ "github.com/uw-labs/lichen"
	_ "golang.org/x/tools/cmd/goimports"
	_ "honnef.co/go/tools/cmd/staticcheck"
)

func main() {

}
