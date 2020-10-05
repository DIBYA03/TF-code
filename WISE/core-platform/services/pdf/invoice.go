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

	"github.com/jung-kurt/gofpdf"
)

type Invoice struct {
	BusinessName  string
	BusinessPhone string
	InvoiceNo     string
	IssueDate     string
	ContactName   string
	ContactEmail  string
	Notes         string
	Amount        string
	WisePhone     string
	AccountNumber string
	RoutingNumber string
}

type InvoiceService interface {
	// Generate Invoice
	GenerateInvoice() (*string, error)
}

func NewInvoiceService(invoice Invoice) InvoiceService {
	return &invoice
}

func (invoice *Invoice) GenerateInvoice() (*string, error) {

	pdf := gofpdf.New("P", "mm", LetterPaper, "")
	pdf.AddPage()
	pdf.SetFont(ArialFont, "B", FontSize24)

	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 5, invoice.BusinessName, "", 0, LeftAlign, false, 0, "")

	pdf.SetFont(ArialFont, "", FontSize24)
	pdf.CellFormat(0, 5, "Invoice", "", 1, RightAlign, false, 0, "")

	pdf.Ln(-1)
	pdf.SetFont(ArialFont, "", FontSize10)
	pdf.SetAlpha(0.6, NormalBlend)
	pdf.CellFormat(30, 5, "Phone", "", 0, LeftAlign, false, 0, "")

	pdf.SetAlpha(1.0, NormalBlend)
	pdf.CellFormat(106, 5, invoice.BusinessPhone, "", 0, LeftAlign, false, 0, "")

	pdf.SetFont(ArialFont, "", FontSize10)
	pdf.SetAlpha(0.6, NormalBlend)
	pdf.CellFormat(27, 5, "Invoice number", "", 0, RightAlign, false, 0, "")

	pdf.SetAlpha(1.0, NormalBlend)
	pdf.CellFormat(33, 5, invoice.InvoiceNo, "", 1, RightAlign, false, 0, "")

	pdf.SetFont(ArialFont, "", FontSize10)
	pdf.SetAlpha(0.6, NormalBlend)
	pdf.CellFormat(30, 5, "Account Number", "", 0, LeftAlign, false, 0, "")

	pdf.SetAlpha(1.0, NormalBlend)
	pdf.CellFormat(30, 5, invoice.AccountNumber, "", 0, LeftAlign, false, 0, "")

	pdf.SetAlpha(0.6, NormalBlend)
	pdf.CellFormat(100, 5, "Date of issue", "", 0, RightAlign, false, 0, "")

	pdf.SetAlpha(1.0, NormalBlend)
	pdf.CellFormat(36, 5, invoice.IssueDate, "", 1, RightAlign, false, 0, "")

	pdf.SetFont(ArialFont, "", FontSize10)
	pdf.SetAlpha(0.6, NormalBlend)
	pdf.CellFormat(30, 5, "Routing Number", "", 0, LeftAlign, false, 0, "")

	pdf.SetAlpha(1.0, NormalBlend)
	pdf.CellFormat(166, 5, invoice.RoutingNumber, "", 1, LeftAlign, false, 0, "")

	pdf.Ln(-1)
	pdf.Ln(-1)

	pdf.SetAlpha(0.6, NormalBlend)
	pdf.CellFormat(0, 5, "Bill to", "", 1, LeftAlign, false, 0, "")

	pdf.SetAlpha(1.0, NormalBlend)
	pdf.CellFormat(0, 5, invoice.ContactName, "", 1, LeftAlign, false, 0, "")
	pdf.CellFormat(0, 5, invoice.ContactEmail, "", 1, LeftAlign, false, 0, "")

	pdf.Ln(-1)
	pdf.Ln(-1)

	pdf.SetAlpha(0.6, NormalBlend)
	pdf.CellFormat(10, 5, "", "", 0, RightAlign, false, 0, "")
	pdf.CellFormat(DescriptionCellWidth, 5, "Description", "", 0, LeftAlign, false, 0, "")
	pdf.CellFormat(AmountCellWidth, 5, "Amount", "", 1, RightAlign, false, 0, "")

	pdf.SetAlpha(1.0, NormalBlend)
	pdf.SetDrawColor(170, 183, 196)
	pdf.Line(12, 77, 204.5, 77)
	pdf.Ln(2.5)

	renderInvoiceNotes(pdf, invoice.Notes, invoice.Amount)

	pdf.Ln(155)

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

func renderInvoiceNotes(pdf *gofpdf.Fpdf, notes string, amount string) {
	startY := 77.5
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
