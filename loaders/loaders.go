//go:build js && wasm
// +build js,wasm

package loaders

import (
	"github.com/Nigel2392/jsext/v2/css"
	"github.com/Nigel2392/jsext/v2/elements"
)

var (
	COLOR_MAIN       = css.COLOR_MAIN
	COLOR_ONE        = css.COLOR_ONE
	COLOR_TWO        = css.COLOR_TWO
	COLOR_THREE      = css.COLOR_THREE
	COLOR_FOUR       = css.COLOR_FOUR
	COLOR_FIVE       = css.COLOR_FIVE
	LOADING_TEXT     = "loading"
	BACKGROUND_COLOR = css.BACKGROUND_COLOR
)

// Function for use in the loader.
//   - Works on its own.
//   - Loader needs to be deleted instead of disabled.
//
// https://www.udemy.com/course/css-animation-transitions-and-transforms-creativity-course/
func LoaderRotatingBlock(idContainer, idLoader string) *elements.Element {
	var container = elements.Div().AttrID(idContainer).AttrStyle("display:none")
	container.Div().AttrID(idLoader)
	container.StyleBlock(`
	#` + idContainer + `{
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: ` + BACKGROUND_COLOR + `;
		display: flex !important;
		justify-content: center !important;
		align-items: center !important;
		perspective: 200px;
	}
	#` + idLoader + `{
		width: 100px;
		height: 100px;
		background-color: ` + COLOR_MAIN + `;
		border-radius: 12px;
		animation: loadingRotatingBlock 2s linear infinite;
	}
	@keyframes loadingRotatingBlock {
		0% { transform: rotateX(0deg) rotateY(0deg); }
		50% { transform: rotateX(180deg) rotateY(0deg); }
		100% { transform: rotateX(180deg) rotateY(180deg); }
	}`)
	return container
}

// Function for use in the loader.
//   - Works on its own.
//   - Loader needs to be deleted instead of disabled.
//
// https://www.udemy.com/course/css-animation-transitions-and-transforms-creativity-course/
func LoaderHexagonRolling(idContainer, idLoader string) *elements.Element {
	var container = elements.Div().AttrID(idContainer).AttrStyle("display:none")
	var line = container.Div().AttrID(idLoader)
	line.Div()

	container.StyleBlock(`#` + idContainer + `{
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: ` + BACKGROUND_COLOR + `;
		display: flex !important;
		justify-content: center !important;
		align-items: center !important;
	}
	#` + idLoader + `{
		width: 300px;
		border-bottom: 4px solid ` + COLOR_MAIN + `;
		position: relative;
		animation: loadingAnimateLine 2s linear infinite;
	}
	#` + idLoader + ` div{
		position: absolute;
		left: 0;
		bottom: 14px;
		width: 50px;
		height: 30px;
		background-color: ` + COLOR_MAIN + `;
		animation: loadingAnimateHexagon 2s linear infinite;
	}
	#` + idLoader + ` div:before{
		content: "";
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: ` + COLOR_MAIN + `;
		transform: rotate(60deg);
	}
	#` + idLoader + ` div:after{
		content: "";
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: ` + COLOR_MAIN + `;
		transform: rotate(-60deg);
	}
	@keyframes loadingAnimateHexagon {
		0% { left: 0; }
		50% { left: calc(100% - 50px); transform: rotate(720deg); }
		100% { left: 0; }
	}
	@keyframes loadingAnimateLine {
		0% { transform: rotate(30deg); }
		25% { transform: rotate(0deg); }
		50% { transform: rotate(-30deg); }
		75% { transform: rotate(0deg); }
		100% { transform: rotate(30deg); }
	}
	`)
	return container
}

// Function for use in the loader.
//   - Works on its own.
//   - Loader needs to be deleted instead of disabled.
//
// https://www.udemy.com/course/css-animation-transitions-and-transforms-creativity-course/
func LoaderMultiRing(idContainer, idLoader string) *elements.Element {
	var container = elements.Div().AttrID(idContainer).AttrStyle("display:none")
	var ld = container.Div().AttrID(idLoader)
	ld.Div()
	ld.Div()
	ld.Div()
	ld.Div()
	var OUTER = " div:nth-child(1)"
	var MIDDLE = " div:nth-child(2)"
	var INNER = " div:nth-child(3)"
	var DOT = " div:nth-child(4)"

	container.StyleBlock(`
	#` + idContainer + `{
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: ` + BACKGROUND_COLOR + `;
		display: flex !important;
		justify-content: center !important;
		align-items: center !important;
	}
	#` + idLoader + `{
		width: 200px;
		height: 200px;
		position: relative;
	}
	#` + idLoader + OUTER + `{
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		border-left: 4px solid ` + COLOR_ONE + `;
		border-right: 4px solid ` + COLOR_ONE + `;
		border-top: 4px solid transparent;
		border-bottom: 4px solid transparent;
		border-radius: 50%;
		animation: loadingRotate 3s linear infinite;
	}
	#` + idLoader + MIDDLE + `{
		position: absolute;
		top: 30px;
		left: 30px;
		right: 30px;
		bottom: 30px;
		border-left: 4px solid ` + COLOR_TWO + `;
		border-right: 4px solid ` + COLOR_TWO + `;
		border-top: 4px solid transparent;
		border-bottom: 4px solid transparent;
		border-radius: 50%;
		animation: loadingRotate 2s linear infinite reverse;
	}
	#` + idLoader + INNER + `{
		position: absolute;
		top: 60px;
		left: 60px;
		right: 60px;
		bottom: 60px;
		border-left: 4px solid ` + COLOR_THREE + `;
		border-right: 4px solid ` + COLOR_THREE + `;
		border-top: 4px solid transparent;
		border-bottom: 4px solid transparent;
		border-radius: 50%;
		animation: loadingRotate 1s linear infinite;
	}
	#` + idLoader + DOT + `{
		position: absolute;
		top: 90px;
		left: 90px;
		right: 90px;
		bottom: 90px;
		background-color: ` + COLOR_MAIN + `;
		border-radius: 50%;
	}
	@keyframes loadingRotate{
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}
	`)
	return container
}

// Function for use in the loader.
//   - Works on its own.
//   - Loader needs to be deleted instead of disabled.
//
// https://www.udemy.com/course/css-animation-transitions-and-transforms-creativity-course/
func LoaderRing(idContainer, idLoader string) *elements.Element {
	var container = elements.Div().AttrID(idContainer).AttrStyle("display:none")
	var innerContainer = container.Div().AttrID(idLoader)
	innerContainer.Div(LOADING_TEXT)
	innerContainer.Div()
	container.StyleBlock(`#` + idContainer + `{
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: ` + BACKGROUND_COLOR + `;
		display: flex !important;
		justify-content: center !important;
		align-items: center !important;
	}
	#` + idLoader + `{
		position: relative;
		width: 200px;
		height: 200px;
	}
	#` + idLoader + ` div:nth-child(1){
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		text-align: center;
		line-height: 200px;
		color: ` + COLOR_MAIN + `;
		font-size: 22px;
		font-weight: bold;
		text-transform: uppercase;
	}
	#` + idLoader + ` div:nth-child(2){
		border-left: 4px solid ` + COLOR_MAIN + `;
		border-radius: 50%;
		width: 100%;
		height: 100%;
		animation: loaderRotate 1s linear infinite;
	}
	@keyframes loaderRotate{
		0%{ transform: rotate(0deg); }
		100%{ transform: rotate(360deg); }
	}`)
	return container
}

// Function for use in the loader.
//   - Works on its own.
//   - Loader needs to be deleted instead of disabled.
//
// https://www.udemy.com/course/css-animation-transitions-and-transforms-creativity-course/
func LoaderQuadSquares(idContainer, idLoader string) *elements.Element {
	var container = elements.Div().AttrID(idContainer).AttrStyle("display:none")
	var innerContainer = container.Div().AttrID(idLoader)
	innerContainer.Span()
	innerContainer.Span()
	innerContainer.Span()
	innerContainer.Span()
	container.StyleBlock(`#` + idContainer + `{
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: ` + BACKGROUND_COLOR + `;
		display: flex !important;
		justify-content: center !important;
		align-items: center !important;
	}
	#` + idLoader + `{
		width: 100px;
		height: 100px;
		position: relative;
		animation: loaderQuadSquares 1s linear infinite;
		transform: rotate(45deg);
	}
	#` + idLoader + ` span{
		position: absolute;
		width: 50px;
		height: 50px;
		animation: loaderQuadSquaresRotate 1s linear infinite;
	}
	#` + idLoader + ` span:nth-child(1){
		background-color: ` + COLOR_ONE + `;
		top: 0;
		left: 0;
	}
	#` + idLoader + ` span:nth-child(2){
		background-color: ` + COLOR_TWO + `;
		top: 0;
		right: 0;
	}
	#` + idLoader + ` span:nth-child(3){
		background-color: ` + COLOR_FOUR + `;
		bottom: 0;
		left: 0;
	}
	#` + idLoader + ` span:nth-child(4){
		background-color: ` + COLOR_THREE + `;
		bottom: 0;
		right: 0;
	}
	@keyframes loaderQuadSquares{
		0%{ width: 100px; height: 100px; }
		10%{ width: 100px; height: 100px; }
		50%{ width: 150px; height: 150px; }
		90%{ width: 100px; height: 100px; }
		100%{ width: 100px; height: 100px; }
	}
	@keyframes loaderQuadSquaresRotate{
		0%{ transform: rotate(0deg) }
		10%{ transform: rotate(0deg) }
		50%{ transform: rotate(90deg) }
		90%{ transform: rotate(90deg) }
		100%{ transform: rotate(90deg) }
	}
	`)
	return container
}
