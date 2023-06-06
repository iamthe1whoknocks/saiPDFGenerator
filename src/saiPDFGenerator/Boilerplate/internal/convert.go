package internal

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sync"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

const (
	chromedpLib    = "chromedp"
	wkhtmltopdfLib = "wkhtmltopdf"
)

// convert html from link to pdf depending on chosen library
func (is *InternalService) convert(library string, html []byte) (result interface{}, err error) {
loop:
	switch library {
	case chromedpLib:
		file, err := is.chromedpConvert(html)
		if err != nil {
			return "", fmt.Errorf("convert - chromedpConvert - %w", err)
		}

		// save file to check without s3
		// if err := ioutil.WriteFile("chromedp.pdf", file, 0644); err != nil {
		// 	return "", fmt.Errorf("convert - chromedpConvert - ioutil.WriteFile - %w", err)
		// }

		s3Link, err := is.s3Upload(file)
		if err != nil {
			return file, fmt.Errorf("convert - s3Upload - %w", err)
		}
		return s3Link, nil
	case wkhtmltopdfLib:
		file, err := is.wkhtmltopdfConvert(html)
		if err != nil {
			return "", fmt.Errorf("convert - wkhtmltopdfConvert - %w", err)
		}
		// save file to check without s3
		// if err := ioutil.WriteFile("wkhtmltopdf.pdf", file, 0644); err != nil {
		// 	return "", fmt.Errorf("convert - chromedpConvert - ioutil.WriteFile - %w", err)
		// }

		s3Link, err := is.s3Upload(file)
		if err != nil {
			return file, fmt.Errorf("convert - s3Upload - %w", err)
		}
		return s3Link, nil

	default:
		is.Logger.Debug("handlers - convert - unknown library - go to default library", zap.String("default library", is.Context.GetConfig("specific.default_convert_library", "chromedp").(string)))
		library = is.Context.GetConfig("specific.default_convert_library", "chromedp").(string)
		goto loop
	}
}

// convert html into pdf using chromedp
func (is *InternalService) chromedpConvert(html []byte) (fileData []byte, err error) {
	taskCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()
	var pdfBuffer []byte
	if err := chromedp.Run(taskCtx, pdfGrabberFromFile(html, "body", &pdfBuffer)); err != nil {
		is.Logger.Error("handlers - chromedpConvert - chromedp.Run", zap.Error(err))
		return nil, err
	}

	// if err := ioutil.WriteFile("test.pdf", pdfBuffer, 0644); err != nil {
	// 	log.Fatal(err)
	// }
	return pdfBuffer, nil
}

// func to convert html into pdf
func pdfGrabberFromFile(html []byte, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		actionLoadHTMLContent(html),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

// actionLoadHTMLContent load html from byte slice
func actionLoadHTMLContent(data []byte) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		isSet, isSetLock := false, sync.Mutex{}
		listenerCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		chromedp.ListenTarget(listenerCtx, func(ev interface{}) {
			if _, ok := ev.(*page.EventLoadEventFired); ok {
				// stop listener
				cancel()

				isSetLock.Lock()
				isSet = true
				isSetLock.Unlock()
			}
		})

		frameTree, err := page.GetFrameTree().Do(ctx)
		if err != nil {
			return err
		}

		if err := page.SetDocumentContent(frameTree.Frame.ID, string(data)).Do(ctx); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-listenerCtx.Done():
			isSetLock.Lock()
			defer isSetLock.Unlock()
			if isSet {
				return nil
			}
			return listenerCtx.Err()
		}
	}
}

func (is *InternalService) wkhtmltopdfConvert(data []byte) (output []byte, err error) {

	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)
	pdfg.Grayscale.Set(false)

	// Create a new input page from an URL
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(data))

	// Set options for this page
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
}
