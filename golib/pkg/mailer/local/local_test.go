package local

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vldem/go-code-example/golib/pkg/mailer"
	mock_utils "github.com/vldem/go-code-example/golib/pkg/utils/mocks"
)

type TestCommander struct{}

func (c TestCommander) CombinedOutput(command string, args ...string) ([]byte, error) {
	return []byte{}, nil
}

type TestStdin struct{}

var testOutput = []byte{}

func (s TestStdin) Write(p []byte) (int, error) {
	testOutput = append(testOutput, p...)
	return len(p), nil
}

func (s TestStdin) Close() error {
	return nil
}

func TestOutput(*testing.T) {
	if os.Getenv("GO_WANT_TEST_OUTPUT") != "1" {
		return
	}

	defer os.Exit(0)
	fmt.Printf("sendmail output")
}

func TestLocalMailer(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	f := setUp(t)

	cs := []string{"-test.run=TestOutput", "--"}
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_TEST_OUTPUT=1"}
	stdin := TestStdin{}

	mockUtils := mock_utils.NewMockUtilsInterface(ctrl)
	mockUtils.EXPECT().GetRandIdString(12).Return(f.boundaryId)

	testMailer := &localMailerImplementation{
		Commander: commander{
			Cmd:   cmd,
			stdin: stdin,
		},
		Utils: mockUtils,
	}

	mail := mailer.Mail{
		From:           f.from,
		To:             f.to,
		Cc:             f.cc,
		Bcc:            f.bcc,
		Subject:        f.subject,
		ReplayTo:       f.replayTo,
		Organization:   f.organization,
		CustomHeader:   f.customHeader,
		UnsubscribeUrl: f.unsubscribeUrl,
	}
	mail.TextBody.Write([]byte(f.mailText))
	mail.HtmlBody.Write([]byte(f.mailHTML))

	// act
	out, err := testMailer.SendMail(mail)

	// assert
	require.NoError(t, err)
	assert.Equal(t, "sendmail output", string(out))
	assert.Equal(t, f.expectedResult, string(testOutput))

}
