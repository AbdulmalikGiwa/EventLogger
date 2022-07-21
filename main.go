package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Event struct where all requests will eventually be written to.
type Event struct {
	WebsiteUrl         string          `json:"websiteUrl"`
	SessionId          string          `json:"sessionId"`
	ResizeFrom         Dimension       `json:"resizeFrom"`
	ResizeTo           Dimension       `json:"resizeTo"`
	CopyAndPaste       map[string]bool `json:"copyAndPaste"`       // map[fieldId]true
	FormCompletionTime int             `json:"formCompletionTime"` // Seconds
}

// PostBody struct helps map body of post requests from clients and is used in construction of Event struct
// sync.Mutex is here because the functions that are responsible for writing to the Event struct are receivers
// on PostBody and will make available the Lock method to help with mutual exclusion on Event.
type PostBody struct {
	sync.Mutex           // To avoid race conditions
	WebsiteUrl string    `json:"websiteUrl"`
	SessionId  string    `json:"sessionId"`
	ResizeFrom Dimension `json:"resizeFrom"`
	ResizeTo   Dimension `json:"resizeTo"`
	EventType  string    `json:"eventType"`
	TimeTaken  int       `json:"timeTaken"`
	Pasted     bool      `json:"pasted"`
	FormId     string    `json:"formId"`
}

// SessionStore serves as a session manager to help in fulfilling requirement of multiple requests arriving on the same session at the same time
var SessionStore map[string]Event

// Response represents http response body
type Response struct {
	StatusCode int
	Message    string
}

type Dimension struct {
	Width  string `json:"width"`
	Height string `json:"height"`
}

const port = ":8080"

// 	CopyPasta helps with writing to CopyAndPaste field of Event struct
var CopyPasta map[string]bool

// StartOrContinueSession initiates a session and also helps to map all requests to specific session by session Id
func StartOrContinueSession(p *PostBody) *Event {
	if SessionStore == nil {
		SessionStore = make(map[string]Event)
	}
	event, ok := SessionStore[p.SessionId]
	if !ok {
		fmt.Printf("Starting Session %s .... \n", p.SessionId)
		CopyPasta = make(map[string]bool)
		event = Event{WebsiteUrl: p.WebsiteUrl,
			SessionId:          p.SessionId,
			ResizeFrom:         p.ResizeFrom,
			ResizeTo:           p.ResizeTo,
			CopyAndPaste:       CopyPasta,
			FormCompletionTime: p.TimeTaken,
		}
		SessionStore[p.SessionId] = event
	}
	return &event
}

// HeadersMiddleware sets headers
func HeadersMiddleware(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST")
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	(*w).Header().Set("Content-Type", "application/json")
}

// ServeHTTP is a handler function that processes POST requests to specified path in ListenAndServe Method
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	HeadersMiddleware(&w)
	if r.Method == http.MethodPost {
		// POST Request Body is written to this struct instance
		RequestBody := PostBody{
			WebsiteUrl: "",
			SessionId:  "",
			ResizeFrom: Dimension{
				Width:  "",
				Height: "",
			},
			ResizeTo: Dimension{
				Width:  "",
				Height: "",
			},
			EventType: "",
			TimeTaken: 0,
			Pasted:    false,
			FormId:    "",
		}
		err := json.NewDecoder(r.Body).Decode(&RequestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		RequestBody.Lock()
		if RequestBody.EventType == "copyAndPaste" {
			err = RequestBody.CopyPasteEvent()
		} else if RequestBody.EventType == "screenResize" {
			err = RequestBody.ResizeEvent()
		} else if RequestBody.EventType == "timeTaken" {
			err = RequestBody.TimeTakenEvent()
		}
		RequestBody.Unlock()
		if err != nil {
			_, _ = w.Write(WriteResponse(err.Error(), http.StatusBadRequest))
			return
		}
		_, err = w.Write(WriteResponse("SUCCESS", http.StatusOK))
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(http.StatusOK)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

// LogEvent method prints the struct to console
func (e *Event) LogEvent(s string) {
	fmt.Println(s)
	log.Printf("Event Struct: %#v \n", e)
}

// WriteResponse function instantiates the Response struct for http responses
func WriteResponse(message string, status int) []byte {
	Resp := Response{
		StatusCode: status,
		Message:    message,
	}
	ResponseBody, _ := json.Marshal(Resp)
	return ResponseBody
}

// Validate method is for generic validations to ensure request body makes sense and catch a few errors
func (p *PostBody) Validate() error {
	if p.WebsiteUrl == "" {
		log.Println("Invalid Request")
		return http.ErrBodyNotAllowed
	}
	if p.SessionId == "" {
		log.Println("Invalid Request")
		return http.ErrBodyNotAllowed
	}
	if p.EventType == "" {
		log.Println("Invalid Request")
		return http.ErrBodyNotAllowed
	}
	return nil
}

// ResizeEvent is a method to handle screen resizes, the before and after dimensions.
func (p *PostBody) ResizeEvent() error {
	err := p.Validate()
	if err != nil {
		return err
	}
	EventObject := StartOrContinueSession(p)
	EventObject.ResizeFrom = p.ResizeFrom
	EventObject.ResizeTo = p.ResizeTo
	EventObject.LogEvent("Screen Resize Event:")
	SessionStore[EventObject.SessionId] = *EventObject
	return err
}

// TimeTakenEvent is a method to handle the final stage of construction of the Event struct
func (p *PostBody) TimeTakenEvent() error {
	err := p.Validate()
	if err != nil {
		return err
	}
	EventObject := StartOrContinueSession(p)
	EventObject.FormCompletionTime = p.TimeTaken
	EventObject.LogEvent("Form Submitted.")
	fmt.Printf("Event Fully Constructed, Session %s Ended \n", p.SessionId)
	fmt.Println("........")
	SessionStore[EventObject.SessionId] = *EventObject
	return err
}

// CopyPasteEvent handles copy & paste Event (for each field)
func (p *PostBody) CopyPasteEvent() error {
	err := p.Validate()
	// Custom validation for CopyPasteEvent
	if p.FormId == "" {
		err = http.ErrBodyNotAllowed
	}
	if err != nil {
		return err
	}
	EventObject := StartOrContinueSession(p)
	_, ok := CopyPasta[p.FormId]
	if !ok {
		CopyPasta[p.FormId] = p.Pasted
	}
	EventObject.CopyAndPaste = CopyPasta
	EventObject.LogEvent("Copy And Paste Event:")
	SessionStore[EventObject.SessionId] = *EventObject
	return err
}

func main() {
	http.HandleFunc("/", ServeHTTP)
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(port, nil))
}
