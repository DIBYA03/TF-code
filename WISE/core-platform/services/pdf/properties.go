/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package pdf

import "strings"

const (
	FontSize24 = 24
	FontSize16 = 16
	FontSize10 = 10
	FontSize8  = 8

	DescriptionCellWidth = 150
	AmountCellWidth      = 25

	LetterPaper = "Letter"
	ArialFont   = "Arial"

	NormalBlend = "Normal"

	LeftAlign  = "L"
	RightAlign = "R"
)

func GetNotesLineItems(notes string) []string {
	s := strings.Split(notes, "\n")

	return s

}
