const source = new EventSource("/events")
const output = document.querySelector("#output")
const autoScroll = document.querySelector("#autoscroll")
output.value = ""
let prefix = ""
source.addEventListener('error', () => {
    addOutput("--- Event source connection error ---")
})
source.addEventListener('open', () => {
    addOutput("--- Event source open ---")
})
source.addEventListener('message', (event) => {
    let data = JSON.parse(event.data)
    let skipLine = false
    if (data.context) {
        switch (data.context) {
            case "connect":
                addOutput("--- Connected:", data.line, "---")
                skipLine = true
                break
            case "error":
                addOutput("-!! ", data.line, "!!-")
                skipLine = true
                break
            case "truncate":
                addOutput("-->", data.line, "<--")
                skipLine = true
                break
            default:
                addOutput("(Unknown context:", context, ")")
        }
    }
    if (data.line && !skipLine) {
        addOutput(data.line)
    }
})
output.addEventListener('scroll',()=>{
    if (output.offsetHeight + output.scrollTop >= output.scrollHeight){
        autoScroll.checked = true
    } else {
        autoScroll.checked = false
    }
})
function addOutput(...line) {
    let start = output.selectionStart
    let end = output.selectionEnd
    output.value += prefix + line.join(' ')
    prefix = "\n"
    if ( start !== end && !autoScroll.checked){
        output.setSelectionRange(start, end)
    }
    maybeScrollDown()
}
function maybeScrollDown() {
    if ( autoScroll.checked ){
        output.scrollTo(0, output.scrollHeight)
    }
}
function zeroPad( something ) {
    return (something+"").padStart(2, "0")
}
document.addEventListener('keydown', (event) => {
    switch(event.key) {
        case "Backspace":
            output.value = ""
            prefix = ""
            break
        case "Enter":
            let today = new Date();
            let date = today.getFullYear() + "-" + zeroPad( today.getMonth() + 1 ) + "-" + zeroPad( today.getDate() )
            let time = zeroPad( today.getHours() ) + ":" + zeroPad( today.getMinutes() ) + ":" + zeroPad( today.getSeconds() )
            addOutput("---", date, time, "---")
            break
        case " ":
            autoScroll.checked = !autoScroll.checked
            if (autoScroll.checked) {
                output.setSelectionRange(0,0)
            }
            maybeScrollDown()
            break
    }
})
