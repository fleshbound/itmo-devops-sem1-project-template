package port

import (
	"context"
	"io"
)

type IParserService interface {
	ParseZip(ctx context.Context, filename string, fileBytes []byte) ([]CreateProductParam, error)
	ParseCsv(ctx context.Context, r io.Reader) ([]CreateProductParam, error)
	GenerateCsv(products []CreateProductParam) ([]byte, error)
	CreateZipFile(csvData []byte) ([]byte, error)
}
