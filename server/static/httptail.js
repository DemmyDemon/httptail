const source = new EventSource("/events")
const output = document.querySelector("#output")
const autoScroll = document.querySelector("#autoscroll")
output.value = ""
let prefix = ""
source.addEventListener('error', () => {
    output.value += prefix + "--- Event source connection error ---"
    maybeScrollDown(output)
})
source.addEventListener('open', () => {
    output.value += prefix + "--- Event source opened ---"
    prefix = "\n"
    maybeScrollDown(output)
})
source.addEventListener('message', (event) => {
    let data = JSON.parse(event.data)
    // console.log(data)
    let skipLine = false
    if (data.context) {
        switch (data.context) {
            case "connect":
                output.value += prefix + "--- Connected: " + data.line + " ---"
                skipLine = true
                break
            default:
                output.value += prefix + "(Unknown context: " + data.context + ")"
        }
    }
    if (data.line && !skipLine) {
        output.value += prefix + data.line
    }
    prefix = "\n"
    maybeScrollDown(output)
})
output.addEventListener('scroll',()=>{
    if (output.offsetHeight + output.scrollTop >= output.scrollHeight){
        autoScroll.checked = true
    } else {
        autoScroll.checked = false
    }
})
function maybeScrollDown( elmt ) {
    if ( autoScroll.checked ){
        elmt.scrollTo(0, elmt.scrollHeight)
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
            output.value += prefix + "---- " + date + " " + time + " ----"
            break
    }
})
