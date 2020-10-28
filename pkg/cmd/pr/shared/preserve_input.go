package shared

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/cli/cli/pkg/iostreams"
)

func PreserveInput(io *iostreams.IOStreams, ims *IssueMetadataState, defs Defaults, doPreserve *bool) func() {
	return func() {
		if ims.Body == defs.Body && ims.Title == defs.Title {
			return
		}

		if !*doPreserve {
			return
		}

		out := io.ErrOut

		// this extra newline guards against appending to the end of a survey line
		fmt.Fprintln(out)

		data, err := json.Marshal(ims)
		if err != nil {
			fmt.Fprintf(out, "failed to save input to file: %s\n", err)
			fmt.Fprintln(out, "would have saved:")
			fmt.Fprintf(out, "%v\n", ims)
			return
		}

		dumpPath := fmt.Sprintf("/tmp/gh-dump-%x.json", time.Now().UnixNano())

		err = ioutil.WriteFile(dumpPath, data, 0660)
		if err != nil {
			fmt.Fprintf(out, "failed to save input to file: %s\n", err)
			fmt.Fprintln(out, "would have saved:")
			fmt.Fprintln(out, string(data))
			return
		}

		cs := io.ColorScheme()

		fmt.Fprintf(out, "%s operation failed. input saved to: %s\n", cs.FailureIcon(), dumpPath)
	}
}
