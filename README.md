# EventLogger

An HTTP server that accepts POST requests (JSON) from multiple clients' websites. Each request forms part of
a struct (for that particular visitor) and each stage of construction(representing an event) is logged to the terminal up until the struct is fully constructed. Below is the Event struct constructed.

### Event Struct
```go
type Event struct {
	WebsiteUrl         string
	SessionId          string
	ResizeFrom         Dimension
	ResizeTo           Dimension
	CopyAndPaste       map[string]bool // map[fieldId]true
	FormCompletionTime int // Seconds
}

type Dimension struct {
	Width  string
	Height string
}
```


## Event options
The following are the possible events a frontend client can capture and post.
  - if the screen resizes, the before and after dimensions
  - copy & paste for each field in a form
  - time taken, in seconds, from the first character being typed to submitting the form
  
### Example JSON Requests
```javascript
{
  "eventType": "copyAndPaste",
  "websiteUrl": "https://stripe.com",
  "sessionId": "123123-123123-123123123",
  "pasted": true,
  "formId": "inputCardNumber"
}

{
  "eventType": "screenResize",
  "websiteUrl": "https://stripe.com",
  "sessionId": "123123-123123-123123123",
  "resizeFrom": {
    "width": "1920",
    "height": "1080"
  },
  "resizeTo": {
    "width": "1280",
    "height": "720"
  }
}

{
  "eventType": "timeTaken",
  "websiteUrl": "https://stripe.com",
  "sessionId": "123123-123123-123123123",
  "timeTaken": 72,
}
```
