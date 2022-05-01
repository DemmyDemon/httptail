const source = new EventSource("/events")
const output = document.getElementById("output")
let prefix = ""
source.onmessage = (event) => {
    console.log(event.data)
    output.append(prefix + event.data)
    prefix = "\n"
    output.scrollTo(0, output.scrollHeight)
}
console.log(source)
