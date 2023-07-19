package dom

type NodeType int

const (
	NodeTypeElement NodeType = iota + 1
	NodeTypeAttribute
	NodeTypeText
	NodeTypeCDATASection
	NodeTypeProcessingInstruction
	NodeTypeComment
	NodeTypeDocument
	NodeTypeDocumentType
	NodeTypeDocumentFragment
)

var AllNodeTypes = []NodeType{
	NodeTypeElement,
	NodeTypeAttribute,
	NodeTypeText,
	NodeTypeCDATASection,
	NodeTypeProcessingInstruction,
	NodeTypeComment,
	NodeTypeDocument,
	NodeTypeDocumentType,
	NodeTypeDocumentFragment,
}
