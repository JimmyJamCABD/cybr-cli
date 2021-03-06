package conjur

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var netrcTemplate string = `machine {{ APPLIANCE_URL }}/authn
  login {{ USERNAME }}
  password {{ PASSWORD }}
`

// CreateNetRc create a conjur netrc file
func CreateNetRc(username string, password string) error {
	// creatr ~/.netrc pas
	homeDir, err := GetHomeDirectory()
	if err != nil {
		return err
	}

	conjurrcFileName := fmt.Sprintf("%s/.conjurrc", homeDir)
	url := GetURLFromConjurRc(conjurrcFileName)
	if url == "" {
		return fmt.Errorf("Failed to get appliance url from '%s'. Run 'cam init' to set this file", conjurrcFileName)
	}

	// create the ~/.netrc file
	netrcFileName := fmt.Sprintf("%s/.netrc", homeDir)
	fmt.Print("Replace ~/.netrc file [y]: ")
	// prompt user
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	answer := strings.Replace(text, "\n", "", -1)
	if answer == "" || answer == "y" {
		// create the ~/.netrc file
		netrcContent := strings.Replace(netrcTemplate, "{{ USERNAME }}", username, 1)
		netrcContent = strings.Replace(netrcContent, "{{ PASSWORD }}", password, 1)
		netrcContent = strings.Replace(netrcContent, "{{ APPLIANCE_URL }}", url, 1)

		os.Remove(netrcFileName)
		err = ioutil.WriteFile(netrcFileName, []byte(netrcContent), 0400)
		if err != nil {
			return fmt.Errorf("Failed to write file '%s'. %s", netrcFileName, err)
		}
	}

	return err
}
