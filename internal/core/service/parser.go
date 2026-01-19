package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"supermarket/internal/core/port"
	"time"
)

type ParserService struct{}

func NewParserService() *ParserService {
	return &ParserService{}
}

func (p *ParserService) ParseZip(ctx context.Context, filename string, fileBytes []byte) ([]port.CreateProductParam, error) {
	// Create zip reader from bytes
	zipReader, err := zip.NewReader(bytes.NewReader(fileBytes), int64(len(fileBytes)))
	if err != nil {
		return nil, err
	}

	// Find CSV file in the zip archive
	var csvFile *zip.File
	for _, f := range zipReader.File {
		if strings.HasSuffix(strings.ToLower(f.Name), ".csv") {
			csvFile = f
			break
		}
	}
	if csvFile == nil {
		return nil, nil
	}

	// Open the CSV file
	rc, err := csvFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	// Parse CSV data
	return p.ParseCsv(ctx, rc)
}

func (p *ParserService) ParseCsv(ctx context.Context, r io.Reader) ([]port.CreateProductParam, error) {
	reader := csv.NewReader(r)

	// Read and validate header
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Ensure we have proper columns
	if len(headers) < 5 {
		return nil, nil
	}

	var products []port.CreateProductParam

	// Read records
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // Skip problematic lines
		}

		if len(record) < 5 {
			continue
		}

		// Parse product
		product, err := p.parseCsvRecord(record)
		if err != nil {
			continue // Skip invalid records
		}

		products = append(products, product)
	}

	return products, nil
}

func (p *ParserService) parseCsvRecord(record []string) (port.CreateProductParam, error) {
	var product port.CreateProductParam

	// Parse ID
	id, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		return product, err
	}
	product.Id = id

	// Parse created date
	createdDate, err := time.Parse("2006-01-02", record[1])
	if err != nil {
		createdDate = time.Now()
	}
	product.Create_date = createdDate

	// Parse name
	product.Name = record[2]

	// Parse category
	product.Category = record[3]

	// Parse price
	price, err := strconv.ParseFloat(record[4], 64)
	if err != nil {
		return product, err
	}
	product.PriceCents = int64(price * 100)
	product.PriceStr = strconv.FormatFloat(price, 'f', 2, 64)

	return product, nil
}

func (p *ParserService) GenerateCsv(products []port.CreateProductParam) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	err := writer.Write([]string{"id", "create_date", "name", "category", "price"})
	if err != nil {
		return nil, err
	}

	// Write data rows
	for _, product := range products {
		dateStr := product.Create_date.Format("2006-01-02")
		priceStr := strconv.FormatFloat(float64(product.PriceCents)/100.0, 'f', 2, 64)

		err := writer.Write([]string{
			strconv.FormatInt(product.Id, 10),
			dateStr,
			product.Name,
			product.Category,
			priceStr,
		})
		if err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (p *ParserService) CreateZipFile(csvData []byte) ([]byte, error) {
	var buf bytes.Buffer

	// Create zip writer
	zipWriter := zip.NewWriter(&buf)

	// Create CSV file in zip
	csvWriter, err := zipWriter.Create("data.csv")
	if err != nil {
		return nil, err
	}

	// Write CSV data
	_, err = csvWriter.Write(csvData)
	if err != nil {
		return nil, err
	}

	// Close zip writer (IMPORTANT!)
	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
