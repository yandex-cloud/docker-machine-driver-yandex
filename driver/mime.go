package driver

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"strings"
)

const dockerMachineMIMEBoundary = `DOCKERMACHINEMIMEBOUNDARY`

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func combineTwoCloudConfigs(userInput, dockerMachineInput string) (string, error) {
	var buffer bytes.Buffer
	w := multipart.NewWriter(&buffer)
	err := w.SetBoundary(dockerMachineMIMEBoundary)
	if err != nil {
		return "", err
	}

	// add specialized mime headers to be correct processed by cloudinit/handlers/__init__.py
	addMixedHeader(&buffer, w.Boundary())

	wp, err := w.CreatePart(createMimeHeader("custom-user-data.yaml", ""))
	if err != nil {
		return "", err
	}
	_, err = wp.Write(convertCRtoCRLF(userInput))
	if err != nil {
		return "", err
	}

	wp, err = w.CreatePart(createMimeHeader("docker-machine-yandex-driver.yaml", ""))
	if err != nil {
		return "", err
	}
	_, err = wp.Write(convertCRtoCRLF(dockerMachineInput))
	if err != nil {
		return "", err
	}

	w.Close()

	return buffer.String(), nil
}

func convertCRtoCRLF(s string) []byte {
	return []byte(strings.Replace(s, "\n", "\r\n", -1))
}

//nolint:errcheck
func addMixedHeader(writer io.Writer, boundary string) {
	writer.Write([]byte(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", boundary)))
	writer.Write([]byte("MIME-Version: 1.0\r\n\r\n"))
}

// CreateFormFile is a convenience wrapper around CreatePart. It creates
// a new form-data header with the provided field name and file name.
func createMimeHeader(filename, contentType string) textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, escapeQuotes(filename)))
	if contentType == "" {
		h.Set("Content-Type", "text/cloud-config")
	} else {
		h.Set("Content-Type", contentType)
	}
	return h
}
