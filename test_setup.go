package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
)

var CopyPaste = make(map[string]bool)

const (
	session1 string = "63efcf5f-c3de-4834-b09f-59d96288c7e3"
	session2 string = "74fgdg6g-d4ef-5945-c18g-60e07399d8f4"
)

// ConstructCopyPasteEvent simulates a copy and paste event posted from
// a client website and is used to set up TestCopyAndPasteEvent test.
func ConstructCopyPasteEvent(formId string, session string) *PostBody {
	CopyPasteBody := PostBody{
		SessionId:  session,
		WebsiteUrl: "https://ravelin.com",
		Pasted:     true,
		FormId:     formId,
		ResizeFrom: Dimension{
			Width:  "",
			Height: "",
		},
		ResizeTo: Dimension{
			Width:  "",
			Height: "",
		},
		EventType: "copyAndPaste",
		TimeTaken: 0,
	}
	CopyPaste[CopyPasteBody.FormId] = CopyPasteBody.Pasted

	return &CopyPasteBody
}

//ConstructScreenResizeEvent helps to mock a ScreenResize Event.
//The arguments for this function are initials just to make it shorter.
// r- Resize, F- From, T- To, H- Height, W- Width **/
func ConstructScreenResizeEvent(rFH string, rFW string, rTH string, rTW string, session string) *PostBody {
	ScreenResizeBody := PostBody{
		SessionId:  session,
		WebsiteUrl: "https://ravelin.com",
		Pasted:     false,
		FormId:     "",
		ResizeFrom: Dimension{
			Width:  rFW,
			Height: rFH,
		},
		ResizeTo: Dimension{
			Width:  rTW,
			Height: rTH,
		},
		EventType: "screenResize",
		TimeTaken: 0,
	}

	return &ScreenResizeBody
}

func ConstructTimeTakenEvent(timeTaken int, session string) *PostBody {
	TimeTakenBody := PostBody{
		SessionId:  session,
		WebsiteUrl: "https://ravelin.com",
		Pasted:     false,
		FormId:     "",
		ResizeFrom: Dimension{
			Width:  "",
			Height: "",
		},
		ResizeTo: Dimension{
			Width:  "",
			Height: "",
		},
		EventType: "copyAndPaste",
		TimeTaken: timeTaken,
	}

	return &TimeTakenBody
}

// PostRequest to test that server can receive POST requests with right body
func PostRequest(Body *PostBody) error {
	body, _ := json.Marshal(Body)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	w := httptest.NewRecorder()
	ServeHTTP(w, req)
	res := w.Result().StatusCode
	if res != 200 {
		return errors.New("status code not 200")
	}
	return nil
}

func (e *Event) Reset() {
	var wipeEvent = &Event{}
	*e = *wipeEvent
}
