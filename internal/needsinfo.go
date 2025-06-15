package internal

import (
	"fmt"
)

type Creator struct {
	Program string `json:"program"`
	Version string `json:"version"`
}

type Version struct {
	Creator   Creator `json:"creator"`
	NeedsInfo `json:"needs"`
}

type NeedsInfo map[string]Need

type Need struct {
	// Idea here is to display the thing nicely maybe in the future.
	// But for now we just want to go to the definition
	
	// This should cover most things that we need (I hope)
	Content          string   `json:"content"`
	Docname          string   `json:"docname"`
	ID               string   `json:"id"`
	Lineno           int      `json:"lineno"`
	ReqType          string   `json:"req_type"`
	Safety           string   `json:"safety"`
	SectionName      string   `json:"section_name"`
	Security         string   `json:"security"`
	Status           string   `json:"status"`
	Title            string   `json:"title"`
	Type             string   `json:"type"`
	TypeName         string   `json:"type_name"`
	DocType          string   `json:"doc_type"`
	IsExternal       bool     `json:"is_external"`
	Tags             []string `json:"tags"`
	Approvers        []string `json:"approvers"`
	Hash             string   `json:"hash"`
	Implemented      string   `json:"implemented"` // YES | PARTIAL | NO
	ParentCovered    string   `json:"parent_covered"`
	ParentHasProblem string   `json:"parent_has_problem"`
	Rationale        string   `json:"rationale"`
	ReqCovered       bool     `json:"req_covered"`
	Reviewers        []string `json:"reviewers"`
	SourceCodeLink   []string `json:"source_code_link"`
	TestLink         []string `json:"testlink"`
	TestCovered      bool     `json:"test_covered"`
	// Links
	Realizes    []string `json:"realizes"`
	Links       []string `json:"links"`
	Satisfies   []string `json:"satisfies"`
	Contains    []string `json:"contains"`
	Has         []string `json:"has"`
	Input       []string `json:"input"`
	Output      []string `json:"output"`
	Responsible []string `json:"responsible"`
	ApprovedBy  []string `json:"approved_by"`
	SupportedBy []string `json:"supported_by"`
	Complies    []string `json:"complies"`
	Fulfils     []string `json:"fulfils"`
	Implements  []string `json:"implements"`
	Uses        []string `json:"uses"`
	Includes    []string `json:"includes"`
	IncludedBy  []string `json:"included_by"`
}

type NeedsJsonInfo struct {
	CurrentVersion string             `json:"currentVersion"`
	Project        string             `json:"project"`
	Versions       map[string]Version `json:"versions"`
}

func (n Need) GenerateHoverInfo() string {
	// Type,Status,Implemented
	return fmt.Sprintf("Type: %s\nStatus: %s\nImplemented: %s\n\n %s", n.Type, n.Status, n.Implemented, n.Content)
}
