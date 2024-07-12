#!/bin/bash

# Name des Go-Projekts (passe es an deinen Projektnamen an)
PROJECT_NAME="MWIN-Chat-Server"

# Ausgabeverzeichnis für die Binaries
OUTPUT_DIR="bin"

# Plattformen und Architekturen
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

# Erstelle das Ausgabeverzeichnis, falls es nicht existiert
mkdir -p "$OUTPUT_DIR"

# Erstelle die Binaries für jede Plattform
for PLATFORM in "${PLATFORMS[@]}"; do
    OS_ARCH=(${PLATFORM//\// })
    GOOS=${OS_ARCH[0]}
    GOARCH=${OS_ARCH[1]}

    echo "Erstelle $PROJECT_NAME für $GOOS/$GOARCH..."

    # Setze Umgebungsvariablen für Go
    export GOOS=$GOOS
    export GOARCH=$GOARCH

    # Erstelle das Binary
    go build -o "$OUTPUT_DIR/$PROJECT_NAME-$GOOS-$GOARCH"

    # Optionale Schritte:
    # - Zippe das Binary für Linux-Plattformen
    if [[ $GOOS == "linux" ]]; then
        zip "$OUTPUT_DIR/$PROJECT_NAME-$GOOS-$GOARCH.zip" "$OUTPUT_DIR/$PROJECT_NAME-$GOOS-$GOARCH"
    fi
done

echo "Build abgeschlossen! Binaries befinden sich im Ordner '$OUTPUT_DIR'."
