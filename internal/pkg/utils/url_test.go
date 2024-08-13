package utils

import (
	"net/url"
	"testing"
)

func TestURLToStringPunycode(t *testing.T) {
	u, err := url.Parse("https://xn----8sbddjhbicfsohgbg1aeo.xn--p1ia/pic/file/map_of_sarlat.pdf")
	if err != nil {
		t.Fatalf("Error parsing URL: %v", err)
	}

	expected := "https://xn----8sbddjhbicfsohgbg1aeo.xn--p1ia/pic/file/map_of_sarlat.pdf"
	actual := URLToString(u)
	if actual != expected {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

func TestURLToStringPunycodeWithPort(t *testing.T) {
	u, err := url.Parse("https://xn----8sbddjhbicfsohgbg1aeo.xn--p1ia:8080/pic/file/map_of_sarlat.pdf")
	if err != nil {
		t.Fatalf("Error parsing URL: %v", err)
	}

	expected := "https://xn----8sbddjhbicfsohgbg1aeo.xn--p1ia:8080/pic/file/map_of_sarlat.pdf"
	actual := URLToString(u)
	if actual != expected {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

func TestURLToStringUnicodetoIDNA(t *testing.T) {
	u, err := url.Parse("https://о-змладйвеклблнозеж.xn--p1ia:8080/pic/file/map_of_sarlat.pdf")
	if err != nil {
		t.Fatalf("Error parsing URL: %v", err)
	}

	expected := "https://xn----8sbddjhbicfsohgbg1aeo.xn--p1ia:8080/pic/file/map_of_sarlat.pdf"
	actual := URLToString(u)
	if actual != expected {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

func TestURLToStringUnicodetoIDNAWithPort(t *testing.T) {
	u, err := url.Parse("https://о-змладйвеклблнозеж.xn--p1ia:8080/pic/file/map_of_sarlat.pdf")
	if err != nil {
		t.Fatalf("Error parsing URL: %v", err)
	}

	expected := "https://xn----8sbddjhbicfsohgbg1aeo.xn--p1ia:8080/pic/file/map_of_sarlat.pdf"
	actual := URLToString(u)
	if actual != expected {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

func TestURLwithIPv6(t *testing.T) {
	u, err := url.Parse("https://[2600:4040:23c7:a620:3642:ebaa:ab23:735e]/test")
	if err != nil {
		t.Fatalf("Error parsing URL: %v", err)
	}

	expected := "https://[2600:4040:23c7:a620:3642:ebaa:ab23:735e]/test"
	actual := URLToString(u)
	if actual != expected {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

func TestURLwithIPv6WithPort(t *testing.T) {
	u, err := url.Parse("https://[2600:4040:23c7:a620:3642:ebaa:ab23:735e]:8080/test")
	if err != nil {
		t.Fatalf("Error parsing URL: %v", err)
	}

	expected := "https://[2600:4040:23c7:a620:3642:ebaa:ab23:735e]:8080/test"
	actual := URLToString(u)
	if actual != expected {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}
