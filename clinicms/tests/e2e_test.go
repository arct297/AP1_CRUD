package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func SetupSelenium(t *testing.T) (selenium.WebDriver, string, func()) {
	appURL := "https://3ba3-5-34-127-237.ngrok-free.app"

	service, err := selenium.NewChromeDriverService("chromedriver", 4444)
	if err != nil {
		t.Fatalf("Failed to start ChromeDriver: %v", err)
	}

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{
		Args: []string{
			"--no-sandbox",
			"--disable-gpu",
			"--user-agent=CustomNgrokBypassAgent/1.0",
		},
	})
	wd, err := selenium.NewRemote(caps, "http://localhost:4444/wd/hub")
	if err != nil {
		t.Fatalf("Failed to connect to Selenium WebDriver: %v", err)
	}

	cleanup := func() {
		wd.Quit()
		service.Stop()
	}
	t.Cleanup(cleanup)

	return wd, appURL, cleanup
}

func TestDoctorFiltering(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	wd, appURL, cleanup := SetupSelenium(t)
	defer cleanup()

	doctorsURL := fmt.Sprintf("%s/", appURL)
	if err := wd.Get(doctorsURL); err != nil {
		t.Fatalf("Failed to load doctors page: %v", err)
	}

	title, err := wd.Title()
	if err != nil {
		t.Fatalf("Failed to get page title: %v", err)
	}
	assert.Contains(t, title, "Clinic Management System", "Page title mismatch")

	err = wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		cards, err := wd.FindElements(selenium.ByCSSSelector, ".doctor-card")
		return len(cards) > 0, err
	})
	if err != nil {
		t.Fatalf("Doctors did not load: %v", err)
	}

	oldDoctorsList, err := wd.FindElement(selenium.ByCSSSelector, ".doctors-list")
	if err != nil {
		t.Fatalf("Failed to find doctors list: %v", err)
	}
	oldHTML, _ := oldDoctorsList.Text()

	filterInput, err := wd.FindElement(selenium.ByCSSSelector, "#filter_value")
	if err != nil {
		t.Fatalf("Failed to find filter input: %v", err)
	}
	filterInput.Clear()
	filterInput.SendKeys("Pediatrics")

	filterButton, err := wd.FindElement(selenium.ByCSSSelector, "button")
	if err != nil {
		t.Fatalf("Failed to find filter button: %v", err)
	}
	filterButton.Click()

	err = wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		doctorsList, err := wd.FindElement(selenium.ByCSSSelector, ".doctors-list")
		if err != nil {
			return false, err
		}
		newHTML, _ := doctorsList.Text()
		return newHTML != oldHTML, nil
	})
	if err != nil {
		pageSource, _ := wd.PageSource()
		log.Println("HTML at error moment:", pageSource)

		t.Fatalf("Filter results not updated: %v", err)
	}

	rows, err := wd.FindElements(selenium.ByCSSSelector, ".doctor-card")
	if err != nil {
		t.Fatalf("Failed to find doctor cards: %v", err)
	}
	assert.Equal(t, 2, len(rows), "Expected 2 Pediatric doctors after filtering")

	for _, row := range rows {
		spec, err := row.FindElement(selenium.ByCSSSelector, ".specialization")
		if err != nil {
			t.Fatalf("Failed to find specialization element: %v", err)
		}
		text, _ := spec.Text()
		assert.Equal(t, "Pediatrics", text, "Incorrect specialization found")
	}

	fmt.Println("✅ Test passed: Filtering works correctly!")
	time.Sleep(10 * time.Second)
}
