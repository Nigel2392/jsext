window.jsextLoaded = function(){
    let jsextInit = new EventTarget();
    jsextInit.On = function(callback){
      this.addEventListener("jsextInit", callback);
    }
    return jsextInit;
}()
  
// Instance of the wasm module when it is ready
let instance = null;
// Path to the WASM module
let wasmFile = "https://raw.githubusercontent.com/Nigel2392/jsext/main/framework/examples/pages/main.wasm";

// Set up the loader
var LoadingText = "loading";
var container_id = "jsext-preload-container";
var inner_container_id = "jsext-preload-inner-container";
var BACKGROUND_COLOR = "rgba(0,0,0,0.5)";
var COLOR_MAIN = "#fff";

// Create the loader
var ldContainer = document.createElement("div");
var innerContainer = document.createElement("div");
var text = document.createElement("div");
var loaderRing = document.createElement("div");

// Set loader attributes
ldContainer.setAttribute("id", container_id);
innerContainer.setAttribute("id", inner_container_id);
text.innerHTML = LoadingText;

// Create the loader style
var style = document.createElement("style");
style.type = "text/css";
style.innerHTML = `
#${container_id}{ position: fixed; top: 0; left: 0; width: 100%; height: 100%; background-color: ${BACKGROUND_COLOR}; display: flex !important; justify-content: center !important; align-items: center !important; }
#${inner_container_id}{ position: relative; width: 200px; height: 200px; }
#${inner_container_id} div:nth-child(1){ position: absolute; top: 0; left: 0; width: 100%; height: 100%; text-align: center; line-height: 200px; color: ${COLOR_MAIN}; font-size: 22px; font-weight: bold; text-transform: uppercase; }
#${inner_container_id} div:nth-child(2){ border-left: 4px solid ${COLOR_MAIN}; border-radius: 50%; width: 100%; height: 100%; animation: loaderRotate 1s linear infinite; }
@keyframes loaderRotate{ 0%{ transform: rotate(0deg); } 100%{ transform: rotate(360deg); } }`;

// Append everything together.
innerContainer.appendChild(text);
innerContainer.appendChild(loaderRing);
ldContainer.appendChild(innerContainer);
ldContainer.appendChild(style);

document.addEventListener("DOMContentLoaded", function (event) {
    // Add the loader to the DOM.
    document.body.appendChild(ldContainer);
    // Define default styles.
    document.body.style.fontFamily = "sans-serif";
    document.body.style.margin = "0";
    document.body.style.padding = "0";
    document.body.style.height = "100vh";

    // Fetch initialize a new instance of our WASM module.
    async function FetchAndInstantiate(url, importObject) {
      return fetch(url)
        .then((response) => response.arrayBuffer())
        .then((bytes) => WebAssembly.instantiate(bytes, importObject))
        .then((results) => results.instance);
    }
    // Initialize a new Go instance.
    let go = new Go();
    // Fetch and instantiate the WASM module.
    let mod = FetchAndInstantiate(wasmFile, go.importObject);

    // When the WASM module is ready, run the Go code.
    // Then, predefine functions for sending and receiving messages,
    // and dispatch an event to let the user know that JSExt is ready.
    window.onload = function () {
        mod.then(function (inst) {
            instance = inst;
            go.run(inst);

            // Loader also gets removed inside the WASM module.
            // if (ldContainer != null && ldContainer != undefined) {
            // ldContainer.remove();
            // }
          
            // Send a message through the messages system.
            jsext.runtime.sendMessage = function(typ, message){
              jsext.runtime.eventEmit("jsextMessages", typ, message);
            };

            // OnMessage takes a callback function that takes the message type
            // and the message itself.
            jsext.runtime.onMessage = function(callBack){
              // Wrap the message function to take the message type,
              // and the message itself.
              var msgCallBack = function(event) {
                  callBack(event.args[0], event.args[1]);
              };
              jsext.runtime.eventOn("jsextMessages", msgCallBack);
            };

            // Dispatch the event to let the user know that JSExt is ready.
            window.jsextLoaded.dispatchEvent(new Event("jsextInit"));
        });
    };
});
