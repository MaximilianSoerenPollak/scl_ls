package internal

type Creator struct {
	Program string `json:"program"`
	Version string `json:"version"`
}

type Version struct {
	Creator Creator                `json:"creator"`
	Needs   map[string]Requirement `json:"needs"`
}

type Requirement struct {
	// Idea here is to display the thing nicely maybe in the future.
	// But for now we just want to go to the definition
	Content     string   `json:"content"`
	Docname     string   `json:"docname"`
	ID          string   `json:"id"`
	Lineno      int      `json:"lineno"`
	ReqType     string   `json:"reqType"`
	Safety      string   `json:"safety"`
	Satisfies   []string `json:"satisfies"`
	SectionName string   `json:"sectionName"`
	Security    string   `json:"security"`
	Status      string   `json:"status"`
	Title       string   `json:"title"`
	Type        string   `json:"type"`
	TypeName    string   `json:"typeName"`
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
