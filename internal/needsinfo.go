package internal

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

// "current_version": "0.1",
// "project": "Score Docs-as-Code",
// "versions": {
//   "0.1": {
//     "creator": {
//       "program": "sphinx_needs",
//       "version": "5.1.0"
//     },
//     "needs": {
//       "feat_req__example__some_title": {
//         "content": "With this requirement we can check if the removal of the prefix is working correctly.\nIt should remove id_prefix (SCORE _) as it's defined inside the BUILD file and remove it before it checks the leftover value\nagainst the allowed defined regex in the metamodel\nNote: The ID is different here as the 'folder structure' is as well",
//         "docname": "how-to-integrate/example/index",
//         "external_css": "external_link",
//         "id": "feat_req__example__some_title",
//         "layout": "score",
//         "lineno": 38,
//         "reqtype": "Process",
//         "safety": "ASIL_D",
//         "satisfies": [
//           "SCORE_stkh_req__overall_goals__reuse_of_app_soft"
//         ],
//         "section_name": "Example",
//         "sections": [
//           "Example"
//         ],
//         "security": "YES",
//         "status": "invalid",
//         "title": "Some Title",
//         "type": "feat_req",
//         "type_name": "Feature Requirement"
//       },
