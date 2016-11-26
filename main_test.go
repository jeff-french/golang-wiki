package main

import "testing"

func TestLoadTemplates(t *testing.T) {
	loadTemplates()
	if _, ok := templates["view"]; !ok {
		t.Error("Failed to load view template")
	}
	if _, ok := templates["edit"]; !ok {
		t.Error("Failed to load edit template")
	}
}

func TestSetupConfig(t *testing.T) {
	t.Skip("TODO: test setupConfig")
}

func TestSave(t *testing.T) {
	t.Skip("TODO: test page saving")
}

func TestLoadPage(t *testing.T) {
	t.Skip("TODO: test loadPage")
}

func TestRenderTemplate(t *testing.T) {
	t.Skip("TODO: test renderTemplate")
}

func TestViewHandler(t *testing.T) {
	t.Skip("TODO: test viewHandler")
}

func TestEditHandler(t *testing.T) {
	t.Skip("TODO: test editHandler")
}

func TestSaveHandler(t *testing.T) {
	t.Skip("TODO: test saveHandler")
}

func TestMakeHandler(t *testing.T) {
	t.Skip("TODO: test makeHandler")
}

func TestRootHandler(t *testing.T) {
	t.Skip("TODO: test rootHandler")
}
