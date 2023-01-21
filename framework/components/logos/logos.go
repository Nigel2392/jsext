package logos

import (
	"github.com/Nigel2392/jsext/framework/elements"
	"github.com/Nigel2392/jsext/framework/helpers"
)

func GoLogo(fontSize string) *elements.Element {
	var hash = helpers.FNVHashString("golang" + fontSize)
	var css = `	.golang-logo-go` + hash + `{
		display: inline-block;
		position:relative;
		font-size: ` + fontSize + `;
		font-weight: 1000;
		color:#00dbbe;
		font-style: italic;
		text-transform: uppercase;
		margin-left:8px;
	}
	.golang-logo-span` + hash + `{
		height: calc(` + fontSize + ` * 2 / 25);
		border-radius: 5px;
		background-color: #00dbbe;
	}
	.golang-logo-span` + hash + `:nth-child(1){
		position: absolute;
		top: 30%;
		left: -8%;
		width: calc(` + fontSize + ` / 6);
	}
	.golang-logo-span` + hash + `:nth-child(2){
		position: absolute;
		top: 48%;
		left: -22%;
		width: calc(` + fontSize + ` * 8 / 30);
		height: calc(` + fontSize + ` / 10);
	}
	.golang-logo-span` + hash + `:nth-child(3){
		position: absolute;
		top: 70%;
		left: -10%;
		width: calc(` + fontSize + ` * 4 / 30);
		height: calc(` + fontSize + ` * 2 / 23);
	}
	`
	var goSpan = elements.Span("Go").AttrClass("golang-logo-go"+hash, "noselect").TextAfter()
	goSpan.Span().AttrClass("golang-logo-span" + hash)
	goSpan.Span().AttrClass("golang-logo-span" + hash)
	goSpan.Span().AttrClass("golang-logo-span" + hash)
	goSpan.StyleBlock(css)
	return goSpan
}

func DjangoLogo(elementSize string) *elements.Element {
	var hash = helpers.FNVHashString("django" + elementSize)
	var css = `	.django-logo-dj` + hash + `{
		display: inline-block;
		position:relative;
		font-size: ` + elementSize + `;
		font-weight: 700;
		background-color: #0C4B33;
		color: #fff;
		font-style: italic;
		text-transform: capitalize;
		line-height: ` + elementSize + `;
		text-align: center;
		height: ` + elementSize + `;
		padding: calc(` + elementSize + ` * 0.125) calc(` + elementSize + ` * 0.25);
		margin: 0 calc(` + elementSize + ` * 0.3);
		border-radius: calc(` + elementSize + ` * 0.25);
		box-shadow: 2px 2px 0 2px rgba(12,75,51, 0.5);
	}
	`
	var djSPan = elements.Span("Django").AttrClass("django-logo-dj"+hash, "noselect").TextAfter()
	djSPan.StyleBlock(css)
	return djSPan
}

func JSLogo(elementSize string) *elements.Element {
	var hash = helpers.FNVHashString("javascript" + elementSize)
	var css = `	.js-logo-js` + hash + `{
		display: inline-block;
		font-size: ` + elementSize + `;
		font-weight: 700;
		background: linear-gradient(to bottom, rgba(0,0,0,0.35), rgba(0,0,0,0)), #F7DF1E;
		color: #000;
		text-transform: uppercase;
		text-align: center;
		padding: calc(` + elementSize + ` * 0.125) calc(` + elementSize + ` * 0.25);
		margin: 0 calc(` + elementSize + ` * 0.25);
		border-radius: calc(` + elementSize + ` * 0.05);
	}
	`
	var jsSpan = elements.Span("JS").AttrClass("js-logo-js"+hash, "noselect").TextAfter()
	jsSpan.StyleBlock(css)
	return jsSpan
}

func WasmLogo(elementSize string) *elements.Element {
	var hash = helpers.FNVHashString("wasm" + elementSize)
	var css = `
	.wasm-logo-wa` + hash + `{
		position: relative;
		background: linear-gradient(to right, rgba(0,0,0,0.35), rgba(0,0,0,0)), #654ff0;
		font-size: ` + elementSize + `;
		text-align: right;
		overflow: hidden;
		padding: calc(` + elementSize + ` * 0.125) calc(` + elementSize + ` * 0.25);
		margin: 0 calc(` + elementSize + ` * 0.25);
		border-radius: calc(` + elementSize + ` * 0.05);
	}
	.wasm-logo-wa` + hash + `:before{
		content: "WA";
		position: absolute;
		top: 0;
		left: 50%;
		width: calc(` + elementSize + ` * 0.4);
		height: calc(` + elementSize + ` * 0.4);
		background: #fff;
		border-radius: 50%;
		transform: translate(-50%, -50%);
		clip-path: polygon(100% 50%, 100% 70%, 80% 90%, 50% 100%, 20% 90%, 0% 70%, 0 50%);}
	`
	var waContainer = elements.Span("WA").AttrClass("wasm-logo-wa"+hash, "noselect")
	waContainer.StyleBlock(css)
	return waContainer
}

func NginxLogo(elementSize string) *elements.Element {
	var hash = helpers.FNVHashString("nginx" + elementSize)
	var css = `
	.nginx-logo-n` + hash + `{
		display: inline-block;
		position: relative;
		background-color: #009639;
		font-size: ` + elementSize + `;
		text-align: center;
		padding: calc(` + elementSize + ` * 0.25) calc(` + elementSize + ` * 0.25);
		margin: 0 calc(` + elementSize + ` * 0.15);
		clip-path: polygon(50% 5%, 100% 25%, 100% 75%, 50% 95%, 0 75%, 0 25%);
	}
	`
	var nContainer = elements.Span("N").AttrClass("nginx-logo-n"+hash, "noselect")
	nContainer.StyleBlock(css)
	return nContainer
}

func SQLLogo(elementSize string) *elements.Element {
	var hash = helpers.FNVHashString("sql" + elementSize)
	var css = `
	.sql-logo` + hash + `{
		margin: 0 calc(` + elementSize + ` * 0.15);
		white-space: nowrap;
	}
	.sql-logo-s` + hash + `{
		position: relative;
		display: inline-block;
		width: ` + elementSize + `;
		height: ` + elementSize + `;
		margin-right: calc(` + elementSize + ` * 0.15);
	}
	.sql-logo-disk` + hash + `{
		position: absolute;
		height: calc(` + elementSize + ` * 0.5);
		background-color: #FFB300;
		width: ` + elementSize + `;
	}
	.sql-logo-disk` + hash + `:nth-child(1){
		top: 0;
		left: 0;
		clip-path: ellipse(50% 30% at 50% 50%);
	}
	.sql-logo-disk` + hash + `:nth-child(2){
		top: calc(` + elementSize + ` * 0.35);
		left: 0;
		clip-path: polygon(0% 0%, 25% 15%, 50% 20%, 75% 15%, 100% 0, 100% 80%, 80% 95%, 50% 100%, 20% 95%, 0 80%);
		border-radius: 2px 2px 0 0;
	}
	.sql-logo-disk` + hash + `:nth-child(3){
		top: calc(` + elementSize + ` * 0.8);
		left: 0;
		clip-path: polygon(0% 0%, 25% 15%, 50% 20%, 75% 15%, 100% 0, 100% 80%, 80% 95%, 50% 100%, 20% 95%, 0 80%);
		border-radius: 2px 2px 0 0;
	}
	.sql-logo-text` + hash + `{
		font-size: calc(` + elementSize + ` / 2);
		font-weight: 700;
		color: #fff;
		font-family: "Roboto", sans-serif;
		text-transform: uppercase;
		text-align: center;
		background-color: #FFB300;
		border-radius: calc(` + elementSize + ` * 0.05);
		padding: calc(` + elementSize + ` * 0.05) calc(` + elementSize + ` * 0.05);
	}
	`
	var mainContainer = elements.Span().AttrClass("sql-logo"+hash, "noselect")
	var sContainer = mainContainer.Span().AttrClass("sql-logo-s" + hash)
	sContainer.Span().AttrClass("sql-logo-disk" + hash)
	sContainer.Span().AttrClass("sql-logo-disk" + hash)
	sContainer.Span().AttrClass("sql-logo-disk" + hash)
	mainContainer.Span("SQL").AttrClass("sql-logo-text" + hash)
	mainContainer.StyleBlock(css)
	return mainContainer
}
