package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"regexp"
	"fmt"
	"encoding/json"

	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/go/pkg/beam/options/gcpopts"
	"github.com/apache/beam/sdks/go/pkg/beam/io/pubsubio"
	"github.com/apache/beam/sdks/go/pkg/beam/x/beamx"
	"github.com/apache/beam/sdks/go/pkg/beam/core/util/stringx"
	//"github.com/apache/beam/sdks/go/pkg/beam/x/debug"
	"github.com/apache/beam/sdks/go/pkg/beam/transforms/stats"


	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"
)

var (
	// input = flag.String("input", "gs://apache-beam-samples/shakespeare/kinglear.txt", "File(s) to read.")
	output = flag.String("output", "gs://result-dataproc-bucket/count", "Output file (required).")
	topic = flag.String("topic", "messages-normalized", "File(s) to read.")
	subscription = flag.String("subscription", "messages-normalized-sub-dataproc", "File(s) to read.")
)


var (
	wordRE  = regexp.MustCompile(`[a-zA-Z]+('[a-z])?`)
	empty   = beam.NewCounter("extract", "emptyLines")
	lineLen = beam.NewDistribution("extract", "lineLenDistro")
)

// extractFn is a DoFn that emits the words in a given line.
func extractFn(ctx context.Context, line string, emit func(string)) {
	var msg libmetier.MessageSocial
	log.Printf("%s", line)
	err := json.Unmarshal([]byte(line), &msg)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("%s", msg.Data)
		lineLen.Update(ctx, int64(len(msg.Data)))
		if len(strings.TrimSpace(msg.Data)) == 0 {
			empty.Inc(ctx, 1)
		}
		for _, word := range wordRE.FindAllString(msg.Data, -1) {
			emit(word)
		}
	}

	
}

// formatFn is a DoFn that formats a word and its count as a string.
func formatFn(w string, c int) string {
	return fmt.Sprintf("%s: %v", w, c)
}

// Concept #4: A composite PTransform is a Go function that adds
// transformations to a given pipeline. It is run at construction time and
// works on PCollections as values. For monitoring purposes, the pipeline
// allows scoped naming for composite transforms. The difference between a
// composite transform and a construction helper function is solely in whether
// a scoped name is used.
//
// For example, the CountWords function is a custom composite transform that
// bundles two transforms (ParDo and Count) as a reusable function.

// CountWords is a composite transform that counts the words of a PCollection
// of lines. It expects a PCollection of type string and returns a PCollection
// of type KV<string,int>. The Beam type checker enforces these constraints
// during pipeline construction.
func CountWords(s beam.Scope, lines beam.PCollection) beam.PCollection {
	s = s.Scope("CountWords")

	// Convert lines of text into individual words.
	col := beam.ParDo(s, extractFn, lines)

	// Count the number of times each word occurs.
	return stats.Count(s, col)
}


func main() {
	// If beamx or Go flags are used, flags must be parsed first.
	flag.Parse()
	ctx := context.Background()
	// beam.Init() is an initialization hook that must be called on startup. On
	// distributed runners, it is used to intercept control.
	beam.Init()

	// Input validation is done as usual. Note that it must be after Init().
	if *output == "" {
		log.Fatal("No output provided")
	}

	// Concepts #3 and #4: The pipeline uses the named transform and DoFn.
	p := beam.NewPipeline()
	s := p.Root()

	//lines := textio.Read(s, *input)
	project := gcpopts.GetProject(ctx)
	msg := pubsubio.Read(s, project, *topic, &pubsubio.ReadOptions{
		Subscription: *subscription,
	})

	str := beam.ParDo(s, stringx.FromBytes, msg)
	counted := CountWords(s, str)
	formatted := beam.ParDo(s, formatFn, counted)
	//cap := beam.ParDo(s, strings.ToUpper, str)
	//debug.Print(s, cap)

	textio.Write(s, *output, formatted)

	// Concept #1: The beamx.Run convenience wrapper allows a number of
	// pre-defined runners to be used via the --runner flag.
	if err := beamx.Run(context.Background(), p); err != nil {
		log.Fatalf("Failed to execute job: %v", err)
	}
}
