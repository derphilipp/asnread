package main

import (
	"fmt"
	// "image"
	"log"
	"os"
	"regexp"

	"github.com/gen2brain/go-fitz"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: Please provide exactly one PDF file path as input.")
		os.Exit(1)
	}

	pdfPath := os.Args[1]

	// 1. Open PDF document
	doc, err := fitz.New(pdfPath)
	if err != nil {
		log.Fatalf("Error opening PDF: %v", err)
	}
	defer doc.Close()

	if doc.NumPage() == 0 {
		log.Fatal("The PDF has no pages.")
	}

	// 2. Render the first page in memory as an image
	// (roughly equivalent to 200-300 DPI for reliable QR detection)
	img, err := doc.Image(0)
	if err != nil {
		log.Fatalf("Error rendering PDF page: %v", err)
	}

	// 3. Prepare image for QR-code reader
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		log.Fatalf("Error converting image: %v", err)
	}

	// 4. Decode QR code
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		// The library returns an error when no QR code is found.
		fmt.Println("ERROR: No QR code found on the first page.")
		os.Exit(1)
	}

	// 5. Extract ASN number using regex (supports ASN123 and ASN-123)
	re := regexp.MustCompile(`ASN-?([0-9]+)`)
	matches := re.FindStringSubmatch(result.GetText())

	if len(matches) > 1 {
		// Print only the numeric part for downstream tooling.
		fmt.Println(matches[1])
	} else {
		fmt.Printf("ERROR: QR code found (%s), but it does not contain an ASN pattern.\n", result.GetText())
		os.Exit(1)
	}
}