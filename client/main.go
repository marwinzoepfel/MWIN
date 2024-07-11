package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"client/gui"
	"client/network"

	"gopkg.in/yaml.v2" // Für YAML-Konfiguration (du kannst auch JSON oder andere Formate verwenden)
)

// Config struct für deine Konfiguration
type Config struct {
	ServerAddress string `yaml:"server_address"`
	ClientName    string `yaml:"client_name"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// 1. Konfiguration laden oder erstellen
	configFilePath := "config.yaml" // Dateiname im aktuellen Verzeichnis
	config, err := loadConfig(configFilePath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// 2. Startmodus abfragen
	fmt.Print("Startmodus (start, start manual): ")
	startMode, _ := reader.ReadString('\n')
	startMode = startMode[:len(startMode)-1] // Zeilenumbruch entfernen

	var serverAddress, clientName string
	if startMode == "start" {
		// 3a. Konfigurationswerte verwenden
		serverAddress = config.ServerAddress
		clientName = config.ClientName
	} else if startMode == "start manual" {
		// 3b. Manuelle Eingabe
		serverAddress = network.GetServerAddress(reader)
		clientName = network.GetClientName(reader)
	} else {
		fmt.Println("Ungültiger Startmodus.")
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

// Funktion zum Laden der Konfiguration
func loadConfig(filePath string) (*Config, error) {
	config := &Config{}

	// Konfigurationsdatei öffnen
	file, err := os.Open(filePath)
	if err != nil && !os.IsNotExist(err) { // Fehler, aber nicht "Datei nicht gefunden"
		return nil, err
	}
	defer file.Close()

	// Wenn Datei vorhanden, YAML einlesen
	if file != nil {
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			return nil, err
		}
	} else {
		// Wenn Datei nicht vorhanden, Standardwerte setzen und speichern
		config.ServerAddress = "localhost:8080" // Beispiel-Standardwert
		config.ClientName = "Neuer Benutzer"    // Beispiel-Standardwert
		if err := saveConfig(filePath, config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// Funktion zum Speichern der Konfiguration
func saveConfig(filePath string, config *Config) error {
	// Konfigurationsverzeichnis erstellen, falls nicht vorhanden
	configDir := filepath.Dir(filePath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Konfigurationsdatei öffnen (oder erstellen, falls nicht vorhanden)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// YAML schreiben
	encoder := yaml.NewEncoder(file)
	return encoder.Encode(config)
}
