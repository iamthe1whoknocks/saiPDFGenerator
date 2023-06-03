package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

const (
	chromedpLib = "chromedp"
)

// convert html from link to pdf depending on chosen library
func (is *InternalService) convert(library, link string) (pdfLink string, err error) {
	switch library {
	case chromedpLib:
		file, err := is.chromedpConvert(link)
		if err != nil {
			return "", fmt.Errorf("convert - chromedpConvert - %w", err)
		}
		s3Link, err := is.s3Upload(file)
		if err != nil {
			return "", fmt.Errorf("convert - chromedpConvert - %w", err)
		}
		return s3Link, nil

	default:
		file, err := is.chromedpConvert(link)
		if err != nil {
			return "", fmt.Errorf("convert - chromedpConvert - %w", err)
		}
		s3Link, err := is.s3Upload(file)
		if err != nil {
			return "", fmt.Errorf("convert - chromedpConvert - %w", err)
		}
		return s3Link, nil
	}

}

// convert html into pdf using chromedp
func (is *InternalService) chromedpConvert(link string) (fileData []byte, err error) {
	taskCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()
	var pdfBuffer []byte
	if err := chromedp.Run(taskCtx, pdfGrabber(link, "body", &pdfBuffer)); err != nil {
		is.Logger.Error("handlers - chromedpConvert - chromedp.Run", zap.Error(err))
		return nil, err
	}

	// if err := ioutil.WriteFile("test.pdf", pdfBuffer, 0644); err != nil {
	// 	log.Fatal(err)
	// }
	return pdfBuffer, nil
}

// func to convert html into pdf
func pdfGrabber(url string, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		emulation.SetUserAgentOverride("WebScraper 1.0"),
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return fmt.Errorf("pdfGrabber - page.PrintToPDF")
			}
			*res = buf
			return nil
		}),
	}
}
