package polygon

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

// polygonEndpointResolver resolves the flatfiles endpoint from polygon.io.
type polygonEndpointResolver struct{}

func (*polygonEndpointResolver) ResolveEndpoint(
	ctx context.Context,
	params s3.EndpointParameters,
) (smithyendpoints.Endpoint, error) {
	u, err := url.Parse("https://files.polygon.io")
	if err != nil {
		return smithyendpoints.Endpoint{}, err
	}
	return smithyendpoints.Endpoint{
		URI: *u,
	}, nil
}

// DownloadDailyAggs downloads Polygon daily agg files from S3 for
// the given ticker and date range, filters for the ticker, and writes
// out filtered CSV files.
func DownloadDailyAggs(ticker string, startDate, endDate time.Time) error {
	keyId := os.Getenv("AWS_ACCESS_KEY_ID")
	if keyId == "" {
		return fmt.Errorf("missing AWS_ACCESS_KEY_ID")
	}
	secretKeyId := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if secretKeyId == "" {
		return fmt.Errorf("missing AWS_SECRET_ACCESS_KEY")
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("loading AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.EndpointResolverV2 = &polygonEndpointResolver{}
		o.UsePathStyle = true
	})

	downloader := manager.NewDownloader(client)

	ticker = strings.ToUpper(ticker)

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		key := fmt.Sprintf(
			"us_stocks_sip/trades_v1/%04d/%02d/%04d-%02d-%02d.csv.gz",
			d.Year(), d.Month(), d.Year(), d.Month(), d.Day())

		tempFile := fmt.Sprintf("temp_%s.csv.gz", d.Format("2006-01-02"))

		f, err := os.Create(tempFile)
		if err != nil {
			return fmt.Errorf("create temp file: %w", err)
		}

		_, err = downloader.Download(ctx, f, &s3.GetObjectInput{
			Bucket: aws.String("flatfiles"),
			Key:    &key,
		})
		f.Close()

		if err != nil {
			os.Remove(tempFile)
			log.Printf("Skipping %s: %v\n", key, err)
			continue
		}

		if err := filterTickerFromGzip(tempFile, ticker, d); err != nil {
			os.Remove(tempFile)
			return fmt.Errorf("filtering ticker: %w", err)
		}

		os.Remove(tempFile)
		fmt.Printf("Processed %s\n", d.Format("2006-01-02"))
	}

	return nil
}

// filterTickerFromGzip extracts lines for the given ticker from the gzip file
// and writes them (with header) to a plain CSV file named <TICKER>_<YYYY-MM-DD>.csv
func filterTickerFromGzip(gzFile, ticker string, date time.Time) error {
	f, err := os.Open(gzFile)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	scanner := bufio.NewScanner(gz)

	outFile := fmt.Sprintf("%s_%s.csv", ticker, date.Format("2006-01-02"))
	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write header line
	if scanner.Scan() {
		out.WriteString(scanner.Text() + "\n")
	}

	// Write lines that start with ticker,
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, ticker+",") {
			out.WriteString(line + "\n")
		}
	}

	return scanner.Err()
}
