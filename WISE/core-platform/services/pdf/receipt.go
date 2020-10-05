/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package pdf

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

type Receipt struct {
	BusinessName  string
	BusinessPhone string
	InvoiceNo     *string
	ReceiptNo     string
	PaidDate      string
	ContactName   *string
	ContactEmail  *string
	Notes         string
	Amount        string
	WisePhone     string
	CardBrand     string
	CardLast4     string
}

type ReceiptService interface {
	// Generate Invoice
	GenerateReceipt() (*string, error)
}

func NewReceiptService(receipt Receipt) ReceiptService {
	return &receipt
}

func (receipt *Receipt) GenerateReceipt() (*string, error) {
	pdf := gofpdf.New("P", "mm", LetterPaper, "")
	pdf.AddPage()
	pdf.SetFont(ArialFont, "B", FontSize24)

	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 5, receipt.BusinessName, "", 0, LeftAlign, false, 0, "")

	pdf.SetFont(ArialFont, "", FontSize24)
	pdf.CellFormat(0, 5, "Receipt", "", 1, RightAlign, false, 0, "")

	pdf.Ln(-1)
	pdf.SetFont(ArialFont, "", FontSize10)
	pdf.SetAlpha(0.6, NormalBlend)
	pdf.CellFormat(113, 5, receipt.BusinessPhone, "", 0, LeftAlign, false, 0, "")

	pdf.SetFont(ArialFont, "", FontSize10)
	pdf.SetAlpha(0.6, NormalBlend)
	pdf.CellFormat(37, 5, "Receipt number", "", 0, RightAlign, false, 0, "")

	pdf.SetAlpha(1.0, NormalBlend)
	pdf.CellFormat(46, 5, receipt.ReceiptNo, "", 1, RightAlign, false, 0, "")

	var lines = 0.0
	if receipt.InvoiceNo != nil {
		pdf.SetAlpha(0.6, NormalBlend)
		pdf.CellFormat(150, 5, "Invoice number", "", 0, RightAlign, false, 0, "")

		pdf.SetAlpha(1.0, NormalBlend)
		pdf.CellFormat(46, 5, *receipt.InvoiceNo, "", 1, RightAlign, false, 0, "")

		lines = lines + 5
	}

	pdf.SetAlpha(0.6, NormalBlend)
	pdf.CellFormat(150, 5, "Date paid", "", 0, RightAlign, false, 0, "")

	pdf.SetAlpha(1.0, NormalBlend)
	pdf.CellFormat(46, 5, receipt.PaidDate, "", 1, RightAlign, false, 0, "")

	pdf.SetAlpha(0.6, NormalBlend)
	pdf.CellFormat(150, 5, "Payment method", "", 0, RightAlign, false, 0, "")

	pdf.SetAlpha(1.0, NormalBlend)
	pdf.CellFormat(46, 5, strings.ToUpper(receipt.CardBrand)+" - "+receipt.CardLast4, "", 1, RightAlign, false, 0, "")

	pdf.Ln(-1)
	pdf.Ln(-1)
	pdf.Ln(-1)

	if receipt.ContactName != nil || receipt.ContactEmail != nil {
		pdf.SetAlpha(0.6, NormalBlend)
		pdf.CellFormat(0, 5, "Bill to", "", 1, LeftAlign, false, 0, "")
		lines = lines + 5
	}

	pdf.SetAlpha(1.0, NormalBlend)
	if receipt.ContactName != nil {
		pdf.CellFormat(0, 5, *receipt.ContactName, "", 1, LeftAlign, false, 0, "")
		lines = lines + 5
	}

	if receipt.ContactEmail != nil {
		pdf.CellFormat(0, 5, *receipt.ContactEmail, "", 1, LeftAlign, false, 0, "")
		lines = lines + 5
	}

	pdf.Ln(-1)
	pdf.Ln(-1)

	pdf.SetFont(ArialFont, "", FontSize24)
	pdf.SetAlpha(1.0, NormalBlend)
	pdf.CellFormat(0, 5, "$"+receipt.Amount+" paid on "+receipt.PaidDate, "", 1, LeftAlign, false, 0, "")

	pdf.Ln(-1)
	pdf.Ln(-1)

	pdf.SetFont(ArialFont, "", FontSize10)
	pdf.SetAlpha(0.6, NormalBlend)
	pdf.CellFormat(10, 5, "", "", 0, RightAlign, false, 0, "")
	pdf.CellFormat(DescriptionCellWidth, 5, "Description", "", 0, LeftAlign, false, 0, "")
	pdf.CellFormat(AmountCellWidth, 5, "Amount", "", 1, RightAlign, false, 0, "")

	pdf.SetAlpha(1.0, NormalBlend)
	pdf.SetDrawColor(170, 183, 196)
	pdf.Line(12, (82 + lines), 204.5, (82 + lines))
	pdf.Ln(2.5)

	renderReceiptNotes(pdf, lines, receipt.Notes, receipt.Amount)

	pdf.Ln(135)

	pdf.SetAlpha(0.6, NormalBlend)
	pdf.SetFont(ArialFont, "", FontSize8)
	pdf.CellFormat(0, 5, "Copyright 2019 Wise Company", "", 1, LeftAlign, false, 0, "")
	pdf.CellFormat(0, 5, "Wise Company banking services provided by BBVA USA.", "", 1, LeftAlign, false, 0, "")
	pdf.CellFormat(0, 5, "The Wise Company Visa Card is issued pursuant to license from Visa U.S.A", "", 1, LeftAlign, false, 0, "")

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	err := pdf.Output(w)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	w.Flush()

	log.Println(b.Len())

	readBuf, err := ioutil.ReadAll(&b)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	sEnc := base64.StdEncoding.EncodeToString(readBuf)

	return &sEnc, nil
}

func renderReceiptNotes(pdf *gofpdf.Fpdf, lines float64, notes string, amount string) {

	startY := 82.5 + lines
	increment := 10.0

	s := GetNotesLineItems(notes)

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	for i, v := range s {
		pdf.SetFillColor(246, 249, 252)
		pdf.Rect(12, startY, 193, 10, "F")

		pdf.SetAlpha(1.0, NormalBlend)
		pdf.CellFormat(10, 10, "", "", 0, LeftAlign, false, 0, "")

		startPoint := 1
		if i == 0 {
			startPoint = 0
		}
		pdf.CellFormat(DescriptionCellWidth, 10, tr(v),
			"", startPoint, LeftAlign, true, 0, "")

		if i == 0 {
			pdf.CellFormat(AmountCellWidth, 10, "$"+amount,
				"", 1, RightAlign, true, 0, "")
		}

		startY = startY + increment

	}

	pdf.SetAlpha(1.0, NormalBlend)
	pdf.SetDrawColor(170, 183, 196)
	pdf.Line(12, startY, 204.5, startY)

}
