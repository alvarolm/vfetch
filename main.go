package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
	flag.Parse()

	fetchNames := flag.Args()

	if configPath == "" {
		fmt.Printf("Using default config path since none was provided\n")
		fmt.Printf("To specify a custom path, use: vfetch -config <path>\n")
		configPath = "vfetch-config.json"
	}
	fmt.Printf("Config path: %s\n", configPath)

	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := ValidateConfig(config); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	fetchItems, err := FilterFetchItems(config, fetchNames)
	if err != nil {
		log.Fatalf("Failed to filter fetch items: %v", err)
	}

	for i, fetchItem := range fetchItems {
		fmt.Printf("Processing item %d: %s\n", i+1, fetchItem.Name)

		if err := ProcessFetchItem(config, fetchItem); err != nil {
			log.Fatalf("Failed to process fetch item %s: %v", fetchItem.Name, err)
		}

		fmt.Printf("Successfully processed: %s\n", fetchItem.Name)
	}

	fmt.Println("All items processed successfully")
}
