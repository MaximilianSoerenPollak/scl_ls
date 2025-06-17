package internal

import (
	"fmt"
	"sclls/lsp"
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
	Content          string      `json:"content,omitempty"`
	Docname          string      `json:"docname,omitempty"`
	ID               string      `json:"id,omitempty"`
	Lineno           int         `json:"lineno,omitempty"`
	ReqType          string      `json:"req_type,omitempty"`
	Safety           string      `json:"safety,omitempty"`
	SectionName      string      `json:"section_name,omitempty"`
	Security         string      `json:"security,omitempty"`
	Status           string      `json:"status,omitempty"`
	Title            string      `json:"title,omitempty"`
	Type             string      `json:"type,omitempty"`
	TypeName         string      `json:"type_name,omitempty"`
	DocType          string      `json:"doc_type,omitempty"`
	IsExternal       bool        `json:"is_external,omitempty"`
	Tags             StringSlice `json:"tags,omitempty"`
	Approvers        StringSlice `json:"approvers,omitempty"`
	Hash             string      `json:"hash,omitempty"`
	Implemented      string      `json:"implemented,omitempty"` // YES | PARTIAL | NO
	ParentCovered    string      `json:"parent_covered,omitempty"`
	ParentHasProblem string      `json:"parent_has_problem,omitempty"`
	Rationale        string      `json:"rationale,omitempty"`
	ReqCovered       bool        `json:"req_covered,omitempty"`
	Reviewers        StringSlice `json:"reviewers,omitempty"`
	SourceCodeLink   StringSlice `json:"source_code_link,omitempty"`
	TestLink         StringSlice `json:"testlink,omitempty"`
	TestCovered      bool        `json:"test_covered,omitempty"`
	// Links
	Realizes    StringSlice `json:"realizes,omitempty"`
	Links       StringSlice `json:"links,omitempty"`
	Satisfies   StringSlice `json:"satisfies,omitempty"`
	Contains    StringSlice `json:"contains,omitempty"`
	Has         StringSlice `json:"has,omitempty"`
	Input       StringSlice `json:"input,omitempty"`
	Output      StringSlice `json:"output,omitempty"`
	Responsible StringSlice `json:"responsible,omitempty"`
	ApprovedBy  StringSlice `json:"approved_by,omitempty"`
	SupportedBy StringSlice `json:"supported_by,omitempty"`
	Complies    StringSlice `json:"complies,omitempty"`
	Fulfils     StringSlice `json:"fulfils,omitempty"`
	Implements  StringSlice `json:"implements,omitempty"`
	Uses        StringSlice `json:"uses,omitempty"`
	Includes    StringSlice `json:"includes,omitempty"`
	IncludedBy  StringSlice `json:"included_by,omitempty"`
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

func (n Need) GenerateCompletionInfo() lsp.CompletionItem {
	return lsp.CompletionItem{
		Label:         n.ID,
		Detail:        fmt.Sprintf("Type: %s\nStatus: %s\nImplemented: %s\n\n %s", n.Type, n.Status, n.Implemented),
		Documentation: n.Content,
	}
}
