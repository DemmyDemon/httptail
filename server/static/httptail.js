const source = new EventSource("/events")
source.onmessage = (event) => {
    console.log(event.data)
}
console.log(source)
