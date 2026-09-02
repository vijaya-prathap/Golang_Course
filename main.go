package main

import "fmt"

type ServerConfig struct {
    Host      string
    Port      int
    Debug     bool
    Tags      []string
    Metadata  map[string]string
}

func main() {
    // No fields set explicitly - Go still gives every field a usable zero value
    // instead of leaving it undefined/garbage, unlike C.
    var cfg ServerConfig

    fmt.Printf("Host:     %q\n", cfg.Host)
    fmt.Printf("Port:     %d\n", cfg.Port)
    fmt.Printf("Debug:    %v\n", cfg.Debug)
    fmt.Printf("Tags:     %v (nil? %v)\n", cfg.Tags, cfg.Tags == nil)
    fmt.Printf("Metadata: %v (nil? %v)\n", cfg.Metadata, cfg.Metadata == nil)

    // Zero values mean the struct is immediately safe to use even before
    // any explicit initialization - e.g. reading from a nil map is fine.
    fmt.Printf("Lookup on nil map: %q\n", cfg.Metadata["region"])

    // Apply defaults only where the zero value isn't good enough.
    if cfg.Host == "" {
        cfg.Host = "localhost"
    }
    if cfg.Port == 0 {
        cfg.Port = 8080
    }
    fmt.Printf("After defaults -> Host: %s, Port: %d\n", cfg.Host, cfg.Port)
}