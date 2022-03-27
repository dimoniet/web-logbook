package pdfexport

import (
	"fmt"
	"io"

	"github.com/jung-kurt/gofpdf"
	"github.com/vsimakhin/web-logbook/internal/models"
)

// printA4LogbookHeader prints header
func (l *Logbook) printA4LogbookHeader() {

	l.pdf.SetFillColor(217, 217, 217)
	l.pdf.SetFont(fontBold, "", 8)

	// First header
	l.pdf.SetXY(leftMargin, topMargin)
	x, y := l.pdf.GetXY()
	for i, str := range header1 {
		width := w1[i]
		l.pdf.Rect(x, y-1, width, 5, "FD")
		l.pdf.MultiCell(width, 1, str, "", "C", false)
		x += width
		l.pdf.SetXY(x, y)
	}
	l.pdf.Ln(-1)

	// Second header
	l.pdf.SetXY(leftMargin, topMargin+3)
	x, y = l.pdf.GetXY()
	for i, str := range header2 {
		width := w2[i]
		l.pdf.Rect(x, y-1, width, 12, "FD")
		l.pdf.MultiCell(width, 3, str, "", "C", false)
		x += width
		l.pdf.SetXY(x, y)
	}
	l.pdf.Ln(-1)

	// Header inside header
	l.pdf.SetXY(leftMargin, topMargin+11)
	x, y = l.pdf.GetXY()
	for i, str := range header3 {
		width := w3[i]
		if str != "" {
			l.pdf.Rect(x, y-1, width, 4, "FD")
			l.pdf.MultiCell(width, 2, str, "", "C", false)
		}
		x += width
		l.pdf.SetXY(x, y)
	}
	l.pdf.Ln(-1)

	// Align the logbook body
	_, y = l.pdf.GetXY()
	y += 1
	l.pdf.SetY(y)
}

// printA4LogbookFooter prints footer
func (l *Logbook) printA4LogbookFooter() {

	printTotal := func(totalName string, total models.FlightRecord) {
		l.pdf.SetFillColor(217, 217, 217)
		l.pdf.SetFont(fontBold, "", 8)

		l.pdf.SetX(leftMargin)

		if totalName == "TOTAL THIS PAGE" {
			l.pdf.CellFormat(w4[0], footerRowHeight, "", "LTR", 0, "", true, 0, "")
		} else if totalName == "TOTAL FROM PREVIOUS PAGES" {
			l.pdf.CellFormat(w4[0], footerRowHeight, "", "LR", 0, "", true, 0, "")
		} else {
			l.pdf.CellFormat(w4[0], footerRowHeight, "", "LBR", 0, "", true, 0, "")
		}
		l.pdf.CellFormat(w4[1], footerRowHeight, totalName, "1", 0, "C", true, 0, "")
		l.pdf.CellFormat(w4[2], footerRowHeight, total.Time.SE, "1", 0, "C", true, 0, "")
		l.pdf.CellFormat(w4[3], footerRowHeight, total.Time.ME, "1", 0, "C", true, 0, "")
		l.pdf.CellFormat(w4[4], footerRowHeight, total.Time.MCC, "1", 0, "C", true, 0, "")
		l.pdf.CellFormat(w4[5], footerRowHeight, total.Time.Total, "1", 0, "C", true, 0, "")
		l.pdf.CellFormat(w4[6], footerRowHeight, "", "1", 0, "", true, 0, "")
		l.pdf.CellFormat(w4[7], footerRowHeight, fmt.Sprintf("%d", total.Landings.Day), "1", 0, "C", true, 0, "")
		l.pdf.CellFormat(w4[8], footerRowHeight, fmt.Sprintf("%d", total.Landings.Night), "1", 0, "C", true, 0, "")
		l.pdf.CellFormat(w4[9], footerRowHeight, total.Time.Night, "1", 0, "C", true, 0, "")
		l.pdf.CellFormat(w4[10], footerRowHeight, total.Time.IFR, "1", 0, "C", true, 0, "")
		l.pdf.CellFormat(w4[11], footerRowHeight, total.Time.PIC, "1", 0, "C", true, 0, "")
		l.pdf.CellFormat(w4[12], footerRowHeight, total.Time.CoPilot, "1", 0, "C", true, 0, "")
		l.pdf.CellFormat(w4[13], footerRowHeight, total.Time.Dual, "1", 0, "C", true, 0, "")
		l.pdf.CellFormat(w4[14], footerRowHeight, total.Time.Instructor, "1", 0, "C", true, 0, "")
		l.pdf.CellFormat(w4[15], footerRowHeight, "", "1", 0, "", true, 0, "")
		l.pdf.CellFormat(w4[16], footerRowHeight, total.SIM.Time, "1", 0, "C", true, 0, "")

		l.pdf.SetFont(fontRegular, "", 6)
		if totalName == "TOTAL THIS PAGE" {
			l.pdf.CellFormat(w4[17], footerRowHeight, l.Signature, "LTR", 0, "C", true, 0, "")
		} else if totalName == "TOTAL FROM PREVIOUS PAGES" {
			l.pdf.CellFormat(w4[17], footerRowHeight, "", "LR", 0, "", true, 0, "")
		} else {
			l.pdf.CellFormat(w4[17], footerRowHeight, l.OwnerName, "LBR", 0, "C", true, 0, "")
		}

		l.pdf.Ln(-1)
	}

	printTotal("TOTAL THIS PAGE", totalPage)
	printTotal("TOTAL FROM PREVIOUS PAGES", totalPrevious)
	printTotal("TOTAL TIME", totalTime)
}

// printA4LogbookBody forms and prints the logbook row
func (l *Logbook) printA4LogbookBody(record models.FlightRecord, fill bool) {

	l.pdf.SetFillColor(228, 228, 228)
	l.pdf.SetTextColor(0, 0, 0)
	l.pdf.SetFont(fontRegular, "", 8)

	// 	Data
	l.pdf.SetX(leftMargin)
	l.pdf.CellFormat(w3[0], bodyRowHeight, record.Date, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[1], bodyRowHeight, record.Departure.Place, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[2], bodyRowHeight, record.Departure.Time, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[3], bodyRowHeight, record.Arrival.Place, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[4], bodyRowHeight, record.Arrival.Time, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[5], bodyRowHeight, record.Aircraft.Model, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[6], bodyRowHeight, record.Aircraft.Reg, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[7], bodyRowHeight, record.Time.SE, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[8], bodyRowHeight, record.Time.ME, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[9], bodyRowHeight, record.Time.MCC, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[10], bodyRowHeight, record.Time.Total, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[11], bodyRowHeight, record.PIC, "1", 0, "L", fill, 0, "")
	l.pdf.CellFormat(w3[12], bodyRowHeight, formatLandings(record.Landings.Day), "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[12], bodyRowHeight, formatLandings(record.Landings.Night), "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[14], bodyRowHeight, record.Time.Night, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[15], bodyRowHeight, record.Time.IFR, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[16], bodyRowHeight, record.Time.PIC, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[17], bodyRowHeight, record.Time.CoPilot, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[18], bodyRowHeight, record.Time.Dual, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[19], bodyRowHeight, record.Time.Instructor, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[20], bodyRowHeight, record.SIM.Type, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[21], bodyRowHeight, record.SIM.Time, "1", 0, "C", fill, 0, "")
	l.pdf.CellFormat(w3[22], bodyRowHeight, record.Remarks, "1", 0, "L", fill, 0, "")

	l.pdf.Ln(-1)

	l.pdf.SetX(leftMargin)
}

// ExportA4 creates A4 pdf with logbook in EASA format
func (l *Logbook) ExportA4(flightRecords []models.FlightRecord, w io.Writer) error {
	// init global vars
	l.init(PDFA4)

	// start forming the pdf file
	l.pdf = gofpdf.New("L", "mm", "A4", "")
	l.pdf.SetAutoPageBreak(true, 5)
	l.loadFonts()

	l.pdf.SetLineWidth(.2)

	rowCounter := 0
	pageCounter := 1

	var totalEmpty models.FlightRecord

	l.pdf.AddPage()
	l.printA4LogbookHeader()

	fill := false

	logBookRow := func(item int) {
		rowCounter += 1

		record := flightRecords[item]
		if record.Time.MCC != "" {
			record.Time.ME = ""
		}

		totalPage = models.CalculateTotals(totalPage, record)
		totalTime = models.CalculateTotals(totalTime, record)

		l.printA4LogbookBody(record, fill)

		if rowCounter >= logbookRows {
			l.printA4LogbookFooter()
			l.printPageNumber(pageCounter)

			totalPrevious = totalTime
			totalPage = totalEmpty

			// check for the page breakes to separate logbooks
			if len(l.PageBreaks) > 0 {
				if fmt.Sprintf("%d", pageCounter) == l.PageBreaks[0] {
					l.pdf.AddPage()
					pageCounter = 0

					l.PageBreaks = append(l.PageBreaks[:0], l.PageBreaks[1:]...)
				}
			}

			rowCounter = 0
			pageCounter += 1

			l.pdf.AddPage()
			l.printA4LogbookHeader()
		}
		fill = fillLine(rowCounter, fillRow)

	}

	for i := len(flightRecords) - 1; i >= 0; i-- {
		logBookRow(i)
	}

	// check the last page for the proper format
	var emptyRecord models.FlightRecord
	for i := rowCounter + 1; i <= logbookRows; i++ {
		l.printA4LogbookBody(emptyRecord, fill)
		fill = fillLine(i, fillRow)

	}
	l.printA4LogbookFooter()
	l.printPageNumber(pageCounter)

	err := l.pdf.Output(w)

	return err
}
