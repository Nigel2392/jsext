//go:build js && wasm
// +build js,wasm

package jsext

import "syscall/js"

type Style js.Value

func (s Style) Value() js.Value {
	return js.Value(s)
}

func (s Style) JSExt() Value {
	return Value(s)
}

func (s Style) Get(key string) js.Value {
	return s.Value().Get(key)
}

func (s Style) Set(key string, value interface{}) {
	s.Value().Set(key, value)
}

func (s Style) Unset(key string) {
	s.Value().Delete(key)
}

//https://www.w3schools.com/jsref/dom_obj_style.asp

func (s Style) AlignContent(str ...string) string {
	if len(str) > 0 {
		s.Set("alignContent", str[0])
	}
	return s.Get("alignContent").String()
}
func (s Style) AlignItems(str ...string) string {
	if len(str) > 0 {
		s.Set("alignItems", str[0])
	}
	return s.Get("alignItems").String()
}
func (s Style) AlignSelf(str ...string) string {
	if len(str) > 0 {
		s.Set("alignSelf", str[0])
	}
	return s.Get("alignSelf").String()
}
func (s Style) Animation(str ...string) string {
	if len(str) > 0 {
		s.Set("animation", str[0])
	}
	return s.Get("animation").String()
}
func (s Style) AnimationDelay(str ...string) string {
	if len(str) > 0 {
		s.Set("animationDelay", str[0])
	}
	return s.Get("animationDelay").String()
}
func (s Style) AnimationDirection(str ...string) string {
	if len(str) > 0 {
		s.Set("animationDirection", str[0])
	}
	return s.Get("animationDirection").String()
}
func (s Style) AnimationDuration(str ...string) string {
	if len(str) > 0 {
		s.Set("animationDuration", str[0])
	}
	return s.Get("animationDuration").String()
}
func (s Style) AnimationFillMode(str ...string) string {
	if len(str) > 0 {
		s.Set("animationFillMode", str[0])
	}
	return s.Get("animationFillMode").String()
}
func (s Style) AnimationIterationCount(str ...string) string {
	if len(str) > 0 {
		s.Set("animationIterationCount", str[0])
	}
	return s.Get("animationIterationCount").String()
}
func (s Style) AnimationName(str ...string) string {
	if len(str) > 0 {
		s.Set("animationName", str[0])
	}
	return s.Get("animationName").String()
}
func (s Style) AnimationTimingFunction(str ...string) string {
	if len(str) > 0 {
		s.Set("animationTimingFunction", str[0])
	}
	return s.Get("animationTimingFunction").String()
}
func (s Style) AnimationPlayState(str ...string) string {
	if len(str) > 0 {
		s.Set("animationPlayState", str[0])
	}
	return s.Get("animationPlayState").String()
}
func (s Style) Background(str ...string) string {
	if len(str) > 0 {
		s.Set("background", str[0])
	}
	return s.Get("background").String()
}
func (s Style) BackgroundAttachment(str ...string) string {
	if len(str) > 0 {
		s.Set("backgroundAttachment", str[0])
	}
	return s.Get("backgroundAttachment").String()
}
func (s Style) BackgroundColor(str ...string) string {
	if len(str) > 0 {
		s.Set("backgroundColor", str[0])
	}
	return s.Get("backgroundColor").String()
}
func (s Style) BackgroundImage(str ...string) string {
	if len(str) > 0 {
		s.Set("backgroundImage", str[0])
	}
	return s.Get("backgroundImage").String()
}
func (s Style) BackgroundPosition(str ...string) string {
	if len(str) > 0 {
		s.Set("backgroundPosition", str[0])
	}
	return s.Get("backgroundPosition").String()
}
func (s Style) BackgroundRepeat(str ...string) string {
	if len(str) > 0 {
		s.Set("backgroundRepeat", str[0])
	}
	return s.Get("backgroundRepeat").String()
}
func (s Style) BackgroundClip(str ...string) string {
	if len(str) > 0 {
		s.Set("backgroundClip", str[0])
	}
	return s.Get("backgroundClip").String()
}
func (s Style) BackgroundOrigin(str ...string) string {
	if len(str) > 0 {
		s.Set("backgroundOrigin", str[0])
	}
	return s.Get("backgroundOrigin").String()
}
func (s Style) BackgroundSize(str ...string) string {
	if len(str) > 0 {
		s.Set("backgroundSize", str[0])
	}
	return s.Get("backgroundSize").String()
}
func (s Style) BackfaceVisibility(str ...string) string {
	if len(str) > 0 {
		s.Set("backfaceVisibility", str[0])
	}
	return s.Get("backfaceVisibility").String()
}
func (s Style) Border(str ...string) string {
	if len(str) > 0 {
		s.Set("border", str[0])
	}
	return s.Get("border").String()
}
func (s Style) BorderBottom(str ...string) string {
	if len(str) > 0 {
		s.Set("borderBottom", str[0])
	}
	return s.Get("borderBottom").String()
}
func (s Style) BorderBottomColor(str ...string) string {
	if len(str) > 0 {
		s.Set("borderBottomColor", str[0])
	}
	return s.Get("borderBottomColor").String()
}
func (s Style) BorderBottomLeftRadius(str ...string) string {
	if len(str) > 0 {
		s.Set("borderBottomLeftRadius", str[0])
	}
	return s.Get("borderBottomLeftRadius").String()
}
func (s Style) BorderBottomRightRadius(str ...string) string {
	if len(str) > 0 {
		s.Set("borderBottomRightRadius", str[0])
	}
	return s.Get("borderBottomRightRadius").String()
}
func (s Style) BorderBottomStyle(str ...string) string {
	if len(str) > 0 {
		s.Set("borderBottomStyle", str[0])
	}
	return s.Get("borderBottomStyle").String()
}
func (s Style) BorderBottomWidth(str ...string) string {
	if len(str) > 0 {
		s.Set("borderBottomWidth", str[0])
	}
	return s.Get("borderBottomWidth").String()
}
func (s Style) BorderCollapse(str ...string) string {
	if len(str) > 0 {
		s.Set("borderCollapse", str[0])
	}
	return s.Get("borderCollapse").String()
}
func (s Style) BorderColor(str ...string) string {
	if len(str) > 0 {
		s.Set("borderColor", str[0])
	}
	return s.Get("borderColor").String()
}
func (s Style) BorderImage(str ...string) string {
	if len(str) > 0 {
		s.Set("borderImage", str[0])
	}
	return s.Get("borderImage").String()
}
func (s Style) BorderImageOutset(str ...string) string {
	if len(str) > 0 {
		s.Set("borderImageOutset", str[0])
	}
	return s.Get("borderImageOutset").String()
}
func (s Style) BorderImageRepeat(str ...string) string {
	if len(str) > 0 {
		s.Set("borderImageRepeat", str[0])
	}
	return s.Get("borderImageRepeat").String()
}
func (s Style) BorderImageSlice(str ...string) string {
	if len(str) > 0 {
		s.Set("borderImageSlice", str[0])
	}
	return s.Get("borderImageSlice").String()
}
func (s Style) BorderImageSource(str ...string) string {
	if len(str) > 0 {
		s.Set("borderImageSource", str[0])
	}
	return s.Get("borderImageSource").String()
}
func (s Style) BorderImageWidth(str ...string) string {
	if len(str) > 0 {
		s.Set("borderImageWidth", str[0])
	}
	return s.Get("borderImageWidth").String()
}
func (s Style) BorderLeft(str ...string) string {
	if len(str) > 0 {
		s.Set("borderLeft", str[0])
	}
	return s.Get("borderLeft").String()
}
func (s Style) BorderLeftColor(str ...string) string {
	if len(str) > 0 {
		s.Set("borderLeftColor", str[0])
	}
	return s.Get("borderLeftColor").String()
}
func (s Style) BorderLeftStyle(str ...string) string {
	if len(str) > 0 {
		s.Set("borderLeftStyle", str[0])
	}
	return s.Get("borderLeftStyle").String()
}
func (s Style) BorderLeftWidth(str ...string) string {
	if len(str) > 0 {
		s.Set("borderLeftWidth", str[0])
	}
	return s.Get("borderLeftWidth").String()
}
func (s Style) BorderRadius(str ...string) string {
	if len(str) > 0 {
		s.Set("borderRadius", str[0])
	}
	return s.Get("borderRadius").String()
}
func (s Style) BorderRight(str ...string) string {
	if len(str) > 0 {
		s.Set("borderRight", str[0])
	}
	return s.Get("borderRight").String()
}
func (s Style) BorderRightColor(str ...string) string {
	if len(str) > 0 {
		s.Set("borderRightColor", str[0])
	}
	return s.Get("borderRightColor").String()
}
func (s Style) BorderRightStyle(str ...string) string {
	if len(str) > 0 {
		s.Set("borderRightStyle", str[0])
	}
	return s.Get("borderRightStyle").String()
}
func (s Style) BorderRightWidth(str ...string) string {
	if len(str) > 0 {
		s.Set("borderRightWidth", str[0])
	}
	return s.Get("borderRightWidth").String()
}
func (s Style) BorderSpacing(str ...string) string {
	if len(str) > 0 {
		s.Set("borderSpacing", str[0])
	}
	return s.Get("borderSpacing").String()
}
func (s Style) BorderStyle(str ...string) string {
	if len(str) > 0 {
		s.Set("borderStyle", str[0])
	}
	return s.Get("borderStyle").String()
}
func (s Style) BorderTop(str ...string) string {
	if len(str) > 0 {
		s.Set("borderTop", str[0])
	}
	return s.Get("borderTop").String()
}
func (s Style) BorderTopColor(str ...string) string {
	if len(str) > 0 {
		s.Set("borderTopColor", str[0])
	}
	return s.Get("borderTopColor").String()
}
func (s Style) BorderTopLeftRadius(str ...string) string {
	if len(str) > 0 {
		s.Set("borderTopLeftRadius", str[0])
	}
	return s.Get("borderTopLeftRadius").String()
}
func (s Style) BorderTopRightRadius(str ...string) string {
	if len(str) > 0 {
		s.Set("borderTopRightRadius", str[0])
	}
	return s.Get("borderTopRightRadius").String()
}
func (s Style) BorderTopStyle(str ...string) string {
	if len(str) > 0 {
		s.Set("borderTopStyle", str[0])
	}
	return s.Get("borderTopStyle").String()
}
func (s Style) BorderTopWidth(str ...string) string {
	if len(str) > 0 {
		s.Set("borderTopWidth", str[0])
	}
	return s.Get("borderTopWidth").String()
}
func (s Style) BorderWidth(str ...string) string {
	if len(str) > 0 {
		s.Set("borderWidth", str[0])
	}
	return s.Get("borderWidth").String()
}
func (s Style) Bottom(str ...string) string {
	if len(str) > 0 {
		s.Set("bottom", str[0])
	}
	return s.Get("bottom").String()
}
func (s Style) BoxShadow(str ...string) string {
	if len(str) > 0 {
		s.Set("boxShadow", str[0])
	}
	return s.Get("boxShadow").String()
}
func (s Style) BoxSizing(str ...string) string {
	if len(str) > 0 {
		s.Set("boxSizing", str[0])
	}
	return s.Get("boxSizing").String()
}
func (s Style) CaptionSide(str ...string) string {
	if len(str) > 0 {
		s.Set("captionSide", str[0])
	}
	return s.Get("captionSide").String()
}
func (s Style) CaretColor(str ...string) string {
	if len(str) > 0 {
		s.Set("caretColor", str[0])
	}
	return s.Get("caretColor").String()
}
func (s Style) Clear(str ...string) string {
	if len(str) > 0 {
		s.Set("clear", str[0])
	}
	return s.Get("clear").String()
}
func (s Style) Clip(str ...string) string {
	if len(str) > 0 {
		s.Set("clip", str[0])
	}
	return s.Get("clip").String()
}
func (s Style) Color(str ...string) string {
	if len(str) > 0 {
		s.Set("color", str[0])
	}
	return s.Get("color").String()
}
func (s Style) ColumnCount(str ...string) string {
	if len(str) > 0 {
		s.Set("columnCount", str[0])
	}
	return s.Get("columnCount").String()
}
func (s Style) ColumnFill(str ...string) string {
	if len(str) > 0 {
		s.Set("columnFill", str[0])
	}
	return s.Get("columnFill").String()
}
func (s Style) ColumnGap(str ...string) string {
	if len(str) > 0 {
		s.Set("columnGap", str[0])
	}
	return s.Get("columnGap").String()
}
func (s Style) ColumnRule(str ...string) string {
	if len(str) > 0 {
		s.Set("columnRule", str[0])
	}
	return s.Get("columnRule").String()
}
func (s Style) ColumnRuleColor(str ...string) string {
	if len(str) > 0 {
		s.Set("columnRuleColor", str[0])
	}
	return s.Get("columnRuleColor").String()
}
func (s Style) ColumnRuleStyle(str ...string) string {
	if len(str) > 0 {
		s.Set("columnRuleStyle", str[0])
	}
	return s.Get("columnRuleStyle").String()
}
func (s Style) ColumnRuleWidth(str ...string) string {
	if len(str) > 0 {
		s.Set("columnRuleWidth", str[0])
	}
	return s.Get("columnRuleWidth").String()
}
func (s Style) Columns(str ...string) string {
	if len(str) > 0 {
		s.Set("columns", str[0])
	}
	return s.Get("columns").String()
}
func (s Style) ColumnSpan(str ...string) string {
	if len(str) > 0 {
		s.Set("columnSpan", str[0])
	}
	return s.Get("columnSpan").String()
}
func (s Style) ColumnWidth(str ...string) string {
	if len(str) > 0 {
		s.Set("columnWidth", str[0])
	}
	return s.Get("columnWidth").String()
}
func (s Style) CounterIncrement(str ...string) string {
	if len(str) > 0 {
		s.Set("counterIncrement", str[0])
	}
	return s.Get("counterIncrement").String()
}
func (s Style) CounterReset(str ...string) string {
	if len(str) > 0 {
		s.Set("counterReset", str[0])
	}
	return s.Get("counterReset").String()
}
func (s Style) Cursor(str ...string) string {
	if len(str) > 0 {
		s.Set("cursor", str[0])
	}
	return s.Get("cursor").String()
}
func (s Style) Direction(str ...string) string {
	if len(str) > 0 {
		s.Set("direction", str[0])
	}
	return s.Get("direction").String()
}
func (s Style) Display(str ...string) string {
	if len(str) > 0 {
		s.Set("display", str[0])
	}
	return s.Get("display").String()
}
func (s Style) EmptyCells(str ...string) string {
	if len(str) > 0 {
		s.Set("emptyCells", str[0])
	}
	return s.Get("emptyCells").String()
}
func (s Style) Filter(str ...string) string {
	if len(str) > 0 {
		s.Set("filter", str[0])
	}
	return s.Get("filter").String()
}
func (s Style) Flex(str ...string) string {
	if len(str) > 0 {
		s.Set("flex", str[0])
	}
	return s.Get("flex").String()
}
func (s Style) FlexBasis(str ...string) string {
	if len(str) > 0 {
		s.Set("flexBasis", str[0])
	}
	return s.Get("flexBasis").String()
}
func (s Style) FlexDirection(str ...string) string {
	if len(str) > 0 {
		s.Set("flexDirection", str[0])
	}
	return s.Get("flexDirection").String()
}
func (s Style) FlexFlow(str ...string) string {
	if len(str) > 0 {
		s.Set("flexFlow", str[0])
	}
	return s.Get("flexFlow").String()
}
func (s Style) FlexGrow(str ...string) string {
	if len(str) > 0 {
		s.Set("flexGrow", str[0])
	}
	return s.Get("flexGrow").String()
}
func (s Style) FlexShrink(str ...string) string {
	if len(str) > 0 {
		s.Set("flexShrink", str[0])
	}
	return s.Get("flexShrink").String()
}
func (s Style) FlexWrap(str ...string) string {
	if len(str) > 0 {
		s.Set("flexWrap", str[0])
	}
	return s.Get("flexWrap").String()
}
func (s Style) CssFloat(str ...string) string {
	if len(str) > 0 {
		s.Set("cssFloat", str[0])
	}
	return s.Get("cssFloat").String()
}
func (s Style) Font(str ...string) string {
	if len(str) > 0 {
		s.Set("font", str[0])
	}
	return s.Get("font").String()
}
func (s Style) FontFamily(str ...string) string {
	if len(str) > 0 {
		s.Set("fontFamily", str[0])
	}
	return s.Get("fontFamily").String()
}
func (s Style) FontSize(str ...string) string {
	if len(str) > 0 {
		s.Set("fontSize", str[0])
	}
	return s.Get("fontSize").String()
}
func (s Style) FontStyle(str ...string) string {
	if len(str) > 0 {
		s.Set("fontStyle", str[0])
	}
	return s.Get("fontStyle").String()
}
func (s Style) FontVariant(str ...string) string {
	if len(str) > 0 {
		s.Set("fontVariant", str[0])
	}
	return s.Get("fontVariant").String()
}
func (s Style) FontWeight(str ...string) string {
	if len(str) > 0 {
		s.Set("fontWeight", str[0])
	}
	return s.Get("fontWeight").String()
}
func (s Style) FontSizeAdjust(str ...string) string {
	if len(str) > 0 {
		s.Set("fontSizeAdjust", str[0])
	}
	return s.Get("fontSizeAdjust").String()
}
func (s Style) Height(str ...string) string {
	if len(str) > 0 {
		s.Set("height", str[0])
	}
	return s.Get("height").String()
}
func (s Style) Isolation(str ...string) string {
	if len(str) > 0 {
		s.Set("isolation", str[0])
	}
	return s.Get("isolation").String()
}
func (s Style) JustifyContent(str ...string) string {
	if len(str) > 0 {
		s.Set("justifyContent", str[0])
	}
	return s.Get("justifyContent").String()
}
func (s Style) Left(str ...string) string {
	if len(str) > 0 {
		s.Set("left", str[0])
	}
	return s.Get("left").String()
}
func (s Style) LetterSpacing(str ...string) string {
	if len(str) > 0 {
		s.Set("letterSpacing", str[0])
	}
	return s.Get("letterSpacing").String()
}
func (s Style) LineHeight(str ...string) string {
	if len(str) > 0 {
		s.Set("lineHeight", str[0])
	}
	return s.Get("lineHeight").String()
}
func (s Style) ListStyle(str ...string) string {
	if len(str) > 0 {
		s.Set("listStyle", str[0])
	}
	return s.Get("listStyle").String()
}
func (s Style) ListStyleImage(str ...string) string {
	if len(str) > 0 {
		s.Set("listStyleImage", str[0])
	}
	return s.Get("listStyleImage").String()
}
func (s Style) ListStylePosition(str ...string) string {
	if len(str) > 0 {
		s.Set("listStylePosition", str[0])
	}
	return s.Get("listStylePosition").String()
}
func (s Style) ListStyleType(str ...string) string {
	if len(str) > 0 {
		s.Set("listStyleType", str[0])
	}
	return s.Get("listStyleType").String()
}
func (s Style) Margin(str ...string) string {
	if len(str) > 0 {
		s.Set("margin", str[0])
	}
	return s.Get("margin").String()
}
func (s Style) MarginBottom(str ...string) string {
	if len(str) > 0 {
		s.Set("marginBottom", str[0])
	}
	return s.Get("marginBottom").String()
}
func (s Style) MarginLeft(str ...string) string {
	if len(str) > 0 {
		s.Set("marginLeft", str[0])
	}
	return s.Get("marginLeft").String()
}
func (s Style) MarginRight(str ...string) string {
	if len(str) > 0 {
		s.Set("marginRight", str[0])
	}
	return s.Get("marginRight").String()
}
func (s Style) MarginTop(str ...string) string {
	if len(str) > 0 {
		s.Set("marginTop", str[0])
	}
	return s.Get("marginTop").String()
}
func (s Style) MaxHeight(str ...string) string {
	if len(str) > 0 {
		s.Set("maxHeight", str[0])
	}
	return s.Get("maxHeight").String()
}
func (s Style) MaxWidth(str ...string) string {
	if len(str) > 0 {
		s.Set("maxWidth", str[0])
	}
	return s.Get("maxWidth").String()
}
func (s Style) MinHeight(str ...string) string {
	if len(str) > 0 {
		s.Set("minHeight", str[0])
	}
	return s.Get("minHeight").String()
}
func (s Style) MinWidth(str ...string) string {
	if len(str) > 0 {
		s.Set("minWidth", str[0])
	}
	return s.Get("minWidth").String()
}
func (s Style) ObjectFit(str ...string) string {
	if len(str) > 0 {
		s.Set("objectFit", str[0])
	}
	return s.Get("objectFit").String()
}
func (s Style) ObjectPosition(str ...string) string {
	if len(str) > 0 {
		s.Set("objectPosition", str[0])
	}
	return s.Get("objectPosition").String()
}
func (s Style) Opacity(str ...string) string {
	if len(str) > 0 {
		s.Set("opacity", str[0])
	}
	return s.Get("opacity").String()
}
func (s Style) Order(str ...string) string {
	if len(str) > 0 {
		s.Set("order", str[0])
	}
	return s.Get("order").String()
}
func (s Style) Orphans(str ...string) string {
	if len(str) > 0 {
		s.Set("orphans", str[0])
	}
	return s.Get("orphans").String()
}
func (s Style) Outline(str ...string) string {
	if len(str) > 0 {
		s.Set("outline", str[0])
	}
	return s.Get("outline").String()
}
func (s Style) OutlineColor(str ...string) string {
	if len(str) > 0 {
		s.Set("outlineColor", str[0])
	}
	return s.Get("outlineColor").String()
}
func (s Style) OutlineOffset(str ...string) string {
	if len(str) > 0 {
		s.Set("outlineOffset", str[0])
	}
	return s.Get("outlineOffset").String()
}
func (s Style) OutlineStyle(str ...string) string {
	if len(str) > 0 {
		s.Set("outlineStyle", str[0])
	}
	return s.Get("outlineStyle").String()
}
func (s Style) OutlineWidth(str ...string) string {
	if len(str) > 0 {
		s.Set("outlineWidth", str[0])
	}
	return s.Get("outlineWidth").String()
}
func (s Style) Overflow(str ...string) string {
	if len(str) > 0 {
		s.Set("overflow", str[0])
	}
	return s.Get("overflow").String()
}
func (s Style) OverflowX(str ...string) string {
	if len(str) > 0 {
		s.Set("overflowX", str[0])
	}
	return s.Get("overflowX").String()
}
func (s Style) OverflowY(str ...string) string {
	if len(str) > 0 {
		s.Set("overflowY", str[0])
	}
	return s.Get("overflowY").String()
}
func (s Style) Padding(str ...string) string {
	if len(str) > 0 {
		s.Set("padding", str[0])
	}
	return s.Get("padding").String()
}
func (s Style) PaddingBottom(str ...string) string {
	if len(str) > 0 {
		s.Set("paddingBottom", str[0])
	}
	return s.Get("paddingBottom").String()
}
func (s Style) PaddingLeft(str ...string) string {
	if len(str) > 0 {
		s.Set("paddingLeft", str[0])
	}
	return s.Get("paddingLeft").String()
}
func (s Style) PaddingRight(str ...string) string {
	if len(str) > 0 {
		s.Set("paddingRight", str[0])
	}
	return s.Get("paddingRight").String()
}
func (s Style) PaddingTop(str ...string) string {
	if len(str) > 0 {
		s.Set("paddingTop", str[0])
	}
	return s.Get("paddingTop").String()
}
func (s Style) PageBreakAfter(str ...string) string {
	if len(str) > 0 {
		s.Set("pageBreakAfter", str[0])
	}
	return s.Get("pageBreakAfter").String()
}
func (s Style) PageBreakBefore(str ...string) string {
	if len(str) > 0 {
		s.Set("pageBreakBefore", str[0])
	}
	return s.Get("pageBreakBefore").String()
}
func (s Style) PageBreakInside(str ...string) string {
	if len(str) > 0 {
		s.Set("pageBreakInside", str[0])
	}
	return s.Get("pageBreakInside").String()
}
func (s Style) Perspective(str ...string) string {
	if len(str) > 0 {
		s.Set("perspective", str[0])
	}
	return s.Get("perspective").String()
}
func (s Style) PerspectiveOrigin(str ...string) string {
	if len(str) > 0 {
		s.Set("perspectiveOrigin", str[0])
	}
	return s.Get("perspectiveOrigin").String()
}
func (s Style) Position(str ...string) string {
	if len(str) > 0 {
		s.Set("position", str[0])
	}
	return s.Get("position").String()
}
func (s Style) Quotes(str ...string) string {
	if len(str) > 0 {
		s.Set("quotes", str[0])
	}
	return s.Get("quotes").String()
}
func (s Style) Resize(str ...string) string {
	if len(str) > 0 {
		s.Set("resize", str[0])
	}
	return s.Get("resize").String()
}
func (s Style) Right(str ...string) string {
	if len(str) > 0 {
		s.Set("right", str[0])
	}
	return s.Get("right").String()
}
func (s Style) ScrollBehavior(str ...string) string {
	if len(str) > 0 {
		s.Set("scrollBehavior", str[0])
	}
	return s.Get("scrollBehavior").String()
}
func (s Style) TableLayout(str ...string) string {
	if len(str) > 0 {
		s.Set("tableLayout", str[0])
	}
	return s.Get("tableLayout").String()
}
func (s Style) TabSize(str ...string) string {
	if len(str) > 0 {
		s.Set("tabSize", str[0])
	}
	return s.Get("tabSize").String()
}
func (s Style) TextAlign(str ...string) string {
	if len(str) > 0 {
		s.Set("textAlign", str[0])
	}
	return s.Get("textAlign").String()
}
func (s Style) TextAlignLast(str ...string) string {
	if len(str) > 0 {
		s.Set("textAlignLast", str[0])
	}
	return s.Get("textAlignLast").String()
}
func (s Style) TextDecoration(str ...string) string {
	if len(str) > 0 {
		s.Set("textDecoration", str[0])
	}
	return s.Get("textDecoration").String()
}
func (s Style) TextDecorationColor(str ...string) string {
	if len(str) > 0 {
		s.Set("textDecorationColor", str[0])
	}
	return s.Get("textDecorationColor").String()
}
func (s Style) TextDecorationLine(str ...string) string {
	if len(str) > 0 {
		s.Set("textDecorationLine", str[0])
	}
	return s.Get("textDecorationLine").String()
}
func (s Style) TextDecorationStyle(str ...string) string {
	if len(str) > 0 {
		s.Set("textDecorationStyle", str[0])
	}
	return s.Get("textDecorationStyle").String()
}
func (s Style) TextIndent(str ...string) string {
	if len(str) > 0 {
		s.Set("textIndent", str[0])
	}
	return s.Get("textIndent").String()
}
func (s Style) TextOverflow(str ...string) string {
	if len(str) > 0 {
		s.Set("textOverflow", str[0])
	}
	return s.Get("textOverflow").String()
}
func (s Style) TextShadow(str ...string) string {
	if len(str) > 0 {
		s.Set("textShadow", str[0])
	}
	return s.Get("textShadow").String()
}
func (s Style) TextTransform(str ...string) string {
	if len(str) > 0 {
		s.Set("textTransform", str[0])
	}
	return s.Get("textTransform").String()
}
func (s Style) Top(str ...string) string {
	if len(str) > 0 {
		s.Set("top", str[0])
	}
	return s.Get("top").String()
}
func (s Style) Transform(str ...string) string {
	if len(str) > 0 {
		s.Set("transform", str[0])
	}
	return s.Get("transform").String()
}
func (s Style) TransformOrigin(str ...string) string {
	if len(str) > 0 {
		s.Set("transformOrigin", str[0])
	}
	return s.Get("transformOrigin").String()
}
func (s Style) TransformStyle(str ...string) string {
	if len(str) > 0 {
		s.Set("transformStyle", str[0])
	}
	return s.Get("transformStyle").String()
}
func (s Style) Transition(str ...string) string {
	if len(str) > 0 {
		s.Set("transition", str[0])
	}
	return s.Get("transition").String()
}
func (s Style) TransitionProperty(str ...string) string {
	if len(str) > 0 {
		s.Set("transitionProperty", str[0])
	}
	return s.Get("transitionProperty").String()
}
func (s Style) TransitionDuration(str ...string) string {
	if len(str) > 0 {
		s.Set("transitionDuration", str[0])
	}
	return s.Get("transitionDuration").String()
}
func (s Style) TransitionTimingFunction(str ...string) string {
	if len(str) > 0 {
		s.Set("transitionTimingFunction", str[0])
	}
	return s.Get("transitionTimingFunction").String()
}
func (s Style) TransitionDelay(str ...string) string {
	if len(str) > 0 {
		s.Set("transitionDelay", str[0])
	}
	return s.Get("transitionDelay").String()
}
func (s Style) UnicodeBidi(str ...string) string {
	if len(str) > 0 {
		s.Set("unicodeBidi", str[0])
	}
	return s.Get("unicodeBidi").String()
}
func (s Style) UserSelect(str ...string) string {
	if len(str) > 0 {
		s.Set("userSelect", str[0])
	}
	return s.Get("userSelect").String()
}
func (s Style) VerticalAlign(str ...string) string {
	if len(str) > 0 {
		s.Set("verticalAlign", str[0])
	}
	return s.Get("verticalAlign").String()
}
func (s Style) Visibility(str ...string) string {
	if len(str) > 0 {
		s.Set("visibility", str[0])
	}
	return s.Get("visibility").String()
}
func (s Style) WhiteSpace(str ...string) string {
	if len(str) > 0 {
		s.Set("whiteSpace", str[0])
	}
	return s.Get("whiteSpace").String()
}
func (s Style) Width(str ...string) string {
	if len(str) > 0 {
		s.Set("width", str[0])
	}
	return s.Get("width").String()
}
func (s Style) WordBreak(str ...string) string {
	if len(str) > 0 {
		s.Set("wordBreak", str[0])
	}
	return s.Get("wordBreak").String()
}
func (s Style) WordSpacing(str ...string) string {
	if len(str) > 0 {
		s.Set("wordSpacing", str[0])
	}
	return s.Get("wordSpacing").String()
}
func (s Style) WordWrap(str ...string) string {
	if len(str) > 0 {
		s.Set("wordWrap", str[0])
	}
	return s.Get("wordWrap").String()
}
func (s Style) Widows(str ...string) string {
	if len(str) > 0 {
		s.Set("widows", str[0])
	}
	return s.Get("widows").String()
}
func (s Style) ZIndex(str ...string) string {
	if len(str) > 0 {
		s.Set("zIndex", str[0])
	}
	return s.Get("zIndex").String()
}
