package main

import (
    "fmt"
    "github.com/kisielk/raven-go/raven"
    "os"
    "errors"
    "strings"
    "flag"
    "io/ioutil"
    "platform"
)

type DataMap map[string]string

func (tags DataMap) Set(value string) error {
    tokens := strings.SplitN(value, "=", 2)

    fmt.Printf("%s", tokens)

    if len(tokens) != 2 {
        return errors.New("Unparseable tag format")
    } else {
        tags[tokens[0]] = tokens[1]
        return nil
    }
}

func (tags DataMap) String() string {
    var s []string
    for k, v := range tags {
        s = append(s, fmt.Sprintf("%s=%s", k, v))
    }
    return strings.Join(s, ", ")
}

func main() {
    var dsn string

    // How does this return an error ?
    defaultHostname, _ := os.Hostname()

    project   := flag.String("project", "", "The project name to use for the event")
    timestamp := flag.String("timestamp", "", "The (iso8601) timestamp to use for the event")
    level     := flag.String("level", "error", "The logging level to send")
    logger    := flag.String("logger", "root", "The logger name to use for the event")
    hostname  := flag.String("hostname", defaultHostname, "The host to report for")
    //exception := flag.String("exception", "", "An exception (gdb format) to report for")

    tags := make(DataMap)
    flag.Var(tags, "tag", "List of tags of the form name=value")

    modules := make(DataMap)
    flag.Var(modules, "module", "Add a module of the form name=value")

    extras := make(DataMap)
    flag.Var(extras, "extra", "Add an extra of the form name=value")

    flag.Parse()

    dsn = flag.Arg(0)

    if dsn == "" {
        dsn = os.Getenv("SENTRY_DSN")
    }

    if dsn == "" {
        fmt.Printf("Error: No configuration detected!\n")
        fmt.Printf("You must either pass a DSN to the command, or set the SENTRY_DSN environment variable\n")
        os.Exit(1)
    }

    client, err := raven.NewClient(dsn)

    if err != nil {
        fmt.Printf("could not connect: %v", dsn)
        os.Exit(2)
    }

    bytes, err := ioutil.ReadAll(os.Stdin)

    if err != nil {
        panic("Unable to read from stdin ?")
    }

    plat := platform.GetPlatform()

    event := raven.Event{
        Message: string(bytes),
        ServerName: *hostname,
        Project: *project,
        Level: *level,
        Timestamp: *timestamp,
        Logger: *logger,
        Tags: tags,
        Modules: modules,
        Extra: extras,
        Platform: fmt.Sprintf("%s - %s", plat.OSName, plat.Release),
    }

    sentryErr := client.Capture(&event)

    if sentryErr != nil {
        fmt.Printf("failed: %v\n", err)
        os.Exit(3)
    }
}
