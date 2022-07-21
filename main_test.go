package main

import (
	"testing"
)

func TestCopyAndPasteEvent(t *testing.T) {
	/** First stage of this test will test for just a single copy and paste event on
	 one field, then subsequently tests for multiple copy and paste events on different
	fields will be tested for**/
	ErrorMessage := "Copy and Paste Event failing"
	CopyBody := ConstructCopyPasteEvent("inputCardNumber", session1)
	err := PostRequest(CopyBody)
	if err != nil {
		t.Error(err.Error())
	}
	err = CopyBody.CopyPasteEvent()

	if err != nil {
		t.Error()
	}
	CopyEvent := SessionStore[session1]
	_, ok := CopyEvent.CopyAndPaste["inputCardNumber"]
	if !ok {
		t.Error(ErrorMessage)
	}
	if CopyEvent.CopyAndPaste["inputCardNumber"] != true {
		t.Error(ErrorMessage)
	}

	// Second phase to test if multiple formIds can be added to the CopyPasteEvent map
	CopyBody2 := ConstructCopyPasteEvent("inputEmail", session1)

	err = CopyBody2.CopyPasteEvent()
	if err != nil {
		t.Error(ErrorMessage)
	}
	CopyEvent2 := SessionStore[session1]
	if len(CopyEvent2.CopyAndPaste) != 2 {
		t.Error(ErrorMessage)
	}

	// Third Phase to test CopyPaste event without formId value, Expected to throw error
	CopyBody3 := ConstructCopyPasteEvent("", session1)
	err = CopyBody3.CopyPasteEvent()
	// Note that this is err==nil and not err!=nil
	if err == nil {
		t.Error(ErrorMessage)
	}
}

func TestScreenResizeEvent(t *testing.T) {
	ErrorMessage := "Resize event failing"
	ResizeBody := ConstructScreenResizeEvent("737", "854", "854", "789", session1)
	err := PostRequest(ResizeBody)
	if err != nil {
		t.Error(err.Error())
	}
	err = ResizeBody.ResizeEvent()
	if err != nil {
		t.Error()
	}
	ResizeEvent := SessionStore[session1]
	if ResizeEvent.ResizeFrom.Height != "737" {
		t.Error(ErrorMessage)
	}
	if ResizeEvent.ResizeFrom.Width != "854" {
		t.Error(ErrorMessage)
	}
	if ResizeEvent.ResizeTo.Height != "854" {
		t.Error(ErrorMessage)
	}
	if ResizeEvent.ResizeTo.Width != "789" {
		t.Error(ErrorMessage)
	}
}

// TestStartOrContinueSession to test if the in memory session storage works
func TestStartOrContinueSession(t *testing.T) {
	_, ok := SessionStore[session1]
	if !ok {
		t.Error("Existing session can't be accessed")
	}
}

func TestTimeTakenEvent(t *testing.T) {
	ErrorMessage := "Time taken event failing"
	TimeTakenBody := ConstructTimeTakenEvent(200, session1)
	err := PostRequest(TimeTakenBody)
	if err != nil {
		t.Error(err.Error())
	}
	err = TimeTakenBody.TimeTakenEvent()
	if err != nil {
		t.Error(ErrorMessage)
	}
	TimeTakenEvent := SessionStore[session1]
	if TimeTakenEvent.FormCompletionTime != 200 {
		t.Error(ErrorMessage)
	}
}
