package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"

    "client/gui"
    "client/network"

    "gopkg.in/yaml.v2" // For YAML configuration (you can use JSON or other formats as well)
)

// Config struct for your configuration
type Config struct {
    ServerAddress string `yaml:"server_address"`
    ClientName    string `yaml:"client_name"`
}

func main() {

    reader := bufio.NewReader(os.Stdin)

    // 1. Load or create configuration
    // Get current working directory
    exePath, err := os.Executable()
    if err != nil {
        fmt.Println("Error getting executable path:", err)
        return
    }
    exeDir := filepath.Dir(exePath)

    // Path to the configuration file in the same directory as the binary
    configFilePath := filepath.Join(exeDir, "config.yaml")

    config, err := loadConfig(configFilePath)
    if err != nil {
        fmt.Println("Error loading config:", err)
        return
    }

    // 2. Query start mode
    fmt.Println("Open config.yaml and enter your data here to connect directly to the specified server with 'start'.")
    selectedItem := gui.StartMenu()

    var startMode string

    // Handle the selected item
    if selectedItem != "" {
        fmt.Println("You selected:", selectedItem)
        startMode = selectedItem
        // ... perform actions based on the selected item ...
    } else {
        fmt.Println("No item selected.")
    }

    var serverAddress, clientName string
    if startMode == "Auto Start" {
        // 3a. Use configuration values
        serverAddress = config.ServerAddress
        clientName = config.ClientName
    } else if startMode == "Manual Start" {
        // 3b. Manual input
        serverAddress = network.GetServerAddress(reader)
        clientName = network.GetClientName(reader)
    } else {
        fmt.Println("Invalid start mode.")
        return
    }

    conn, err := network.ConnectToServer(serverAddress)
    if err != nil {
        fmt.Println("Error connecting to server:", err)
        return
    }
    defer conn.Close()

    fmt.Println("Successfully joined.")
    if err := network.SendInitialName(conn, clientName); err != nil {
        fmt.Println("Error sending name:", err)
        return
    }

    // Start GUI
    gui.RunChat(conn, reader)
}

// Function to load the configuration
func loadConfig(filePath string) (*Config, error) {
    config := &Config{}

    // Open configuration file
    file, err := os.Open(filePath)
    if err != nil && !os.IsNotExist(err) { // Error, but not "file not found"
        return nil, err
    }
    defer file.Close()

    // If file exists, read YAML
    if file != nil {
        decoder := yaml.NewDecoder(file)
        if err := decoder.Decode(config); err != nil {
            return nil, err
        }
    } else {
        // If file doesn't exist, set default values and save
        config.ServerAddress = "localhost:8080" // Example default value
        config.ClientName = "New User"         // Example default value
        if err := saveConfig(filePath, config); err != nil {
            return nil, err
        }
    }

    return config, nil
}

// Function to save the configuration
func saveConfig(filePath string, config *Config) error {
    // Create configuration directory if it doesn't exist
    configDir := filepath.Dir(filePath)
    if err := os.MkdirAll(configDir, 0755); err != nil {
        return err
    }

    // Open configuration file (or create if it doesn't exist)
    file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    // Write YAML
    encoder := yaml.NewEncoder(file)
    return encoder.Encode(config)
}
