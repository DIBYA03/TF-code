package awsauthcli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

var defaultRegion = "us-west-2"

// CredsFile is the location of the aws credentials file
var CredsFile = "~/.aws/credentials"

// Profile is a aws profile in credentials file
type Profile struct {
	Name            string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Expiration      time.Time
}

// SignInTokenRequest is the session for sign-in token request
type SignInTokenRequest struct {
	ID    string `json:"sessionId"`
	Key   string `json:"sessionKey"`
	Token string `json:"sessionToken"`
}

// SignInTokenResponse is the resposne from sign-in token
type SignInTokenResponse struct {
	SigninToken string `json:"SigninToken"`
}

// awsSAMLAssumeRole assumes a role and returns back the credential details to use
func awsSAMLAssumeRole(principalArn string, roleArn string, duration int64, samlResponse string) (*sts.AssumeRoleWithSAMLOutput, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(defaultRegion),
	}))

	svc := sts.New(sess)
	input := &sts.AssumeRoleWithSAMLInput{
		DurationSeconds: aws.Int64(duration),
		PrincipalArn:    aws.String(principalArn),
		RoleArn:         aws.String(roleArn),
		SAMLAssertion:   aws.String(samlResponse),
	}

	req, resp := svc.AssumeRoleWithSAMLRequest(input)
	err := req.Send()
	if err != nil {
		return &sts.AssumeRoleWithSAMLOutput{}, err
	}

	return resp, nil
}

// HandleSessionCredentials will create the session cli credentials, then parse
// the credentials to create a new aws creds file. Part of this process is the
// AWS assume role call
// If there's an error here, it will only print them, not stop
func HandleSessionCredentials(roles []string, sessDuration int64, b64Response string, credsFile string) ([]Profile, error) {
	var awsCredentials []Profile

	for _, awsRole := range roles {
		r := strings.Split(awsRole, ",")
		role, provider := r[0], r[1]

		assumedRole, err := awsSAMLAssumeRole(provider, role, sessDuration, b64Response)
		if err != nil {
			fmt.Println(role, ":", err)
			continue
		}

		roleParts := strings.Split(role, "/")

		awsCredentials = append(awsCredentials, Profile{
			Name:            roleParts[1],
			AccessKeyID:     *assumedRole.Credentials.AccessKeyId,
			SecretAccessKey: *assumedRole.Credentials.SecretAccessKey,
			SessionToken:    *assumedRole.Credentials.SessionToken,
			Expiration:      *assumedRole.Credentials.Expiration,
		})

	}

	oldCredentials := readCredentials(credsFile)
	newCredentials := awsCredentials

	for _, oldCredential := range oldCredentials {
		credentialsExist := false

		for _, newCredential := range awsCredentials {
			if newCredential.Name == oldCredential.Name {
				credentialsExist = true
			}
		}

		if !credentialsExist {
			newCredentials = append(newCredentials, oldCredential)
		}
	}

	_, err := createNewCredentialsFile(newCredentials, credsFile)
	if err != nil {
		return awsCredentials, err
	}

	return awsCredentials, nil
}

// readCredentials reads all the current AWS credentials into a list of structs,
// so they can be recreated when adding a new file. This allows to us to keep
// custom aws creds and rewrite the Wise Company ones
func readCredentials(credsFile string) []Profile {
	var profiles []Profile

	file, err := os.Open(credentialsFile(credsFile))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	profile := Profile{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		configLine := scanner.Text()

		// Ignore line breaks
		if configLine == "" {
			continue
		}

		if strings.HasPrefix(configLine, "[") {
			// Make sure the profile is not empty (1st profile not created yet)
			profiles = append(profiles, profile)

			profileName := strings.Replace(configLine, "[", "", -1)
			profileName = strings.Replace(profileName, "]", "", -1)

			// Create new profile
			profile = Profile{
				Name: profileName,
			}

			continue
		}

		key, val := getConfigLineValue(configLine)

		switch key {
		case "aws_access_key_id":
			profile.AccessKeyID = val
			break

		case "aws_secret_access_key":
			profile.SecretAccessKey = val
			break

		case "aws_session_token":
			profile.SessionToken = val
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	profiles = append(profiles, profile)

	return profiles
}

// createNewCredentialsFile creates the aws credentials file and returns back
// bool on if worked and error if one happens
func createNewCredentialsFile(creds []Profile, credsFile string) (bool, error) {
	f, err := os.Create(credentialsFile(credsFile))
	if err != nil {
		return false, err
	}

	for _, cred := range creds {
		if cred.Name != "" {
			_, err := f.WriteString(fmt.Sprintf("[%s]\n", cred.Name))
			if err != nil {
				f.Close()
				return false, err
			}

			_, err = f.WriteString(fmt.Sprintf("aws_access_key_id = %s\n", cred.AccessKeyID))
			if err != nil {
				f.Close()
				return false, err
			}

			_, err = f.WriteString(fmt.Sprintf("aws_secret_access_key = %s\n", cred.SecretAccessKey))
			if err != nil {
				f.Close()
				return false, err
			}

			_, err = f.WriteString(fmt.Sprintf("aws_session_token = %s\n\n", cred.SessionToken))
			if err != nil {
				f.Close()
				return false, err
			}
		}
	}

	err = f.Close()
	if err != nil {
		return false, err
	}

	return true, nil
}

// getConfigLineValue splits a config by `=` for key/value
func getConfigLineValue(configLine string) (string, string) {
	configParts := strings.Split(configLine, "=")

	configKey := strings.TrimSpace(configParts[0])
	configVal := strings.TrimSpace(configParts[1])

	return configKey, configVal
}

// credentialsFile create the location of the aws creds file
func credentialsFile(credsFile string) string {
	if strings.HasPrefix(credsFile, "~/") {
		u, err := user.Current()
		if err != nil {
			log.Fatalf("err: %s", err)
		}

		homeDir := fmt.Sprintf("%s/", u.HomeDir)
		credsFile = strings.Replace(credsFile, "~/", homeDir, 1)
	}

	return credsFile
}

// HandlePreSignInURLs generates presign in URLs that will be used later to get a proper sign in URL from AWS
func HandlePreSignInURLs(profiles []Profile, issuer string, sessionDuration int64) (map[string]string, error) {
	signInURLs := make(map[string]string)

	// Create the get token URL
	url := url.URL{
		Scheme: "https",
		Host:   "localhost:4433",
		Path:   "aws/signin",
	}

	for _, profile := range profiles {
		q := url.Query()
		q.Set("AccessKeyId", profile.AccessKeyID)
		q.Set("SecretAccessKey", profile.SecretAccessKey)
		q.Set("SessionToken", profile.SessionToken)
		q.Set("SessionDuration", strconv.FormatInt(sessionDuration, 10))
		q.Set("Issuer", issuer)
		url.RawQuery = q.Encode()

		fmt.Println(url.String())
		signInURLs[profile.Name] = url.String()
	}

	return signInURLs, nil
}

// HandleSignInURL helps create the sign-in URLs that will be passed back for
func HandleSignInURL(signInTokenRequest SignInTokenRequest, issuer string, sessionDuration int64) (string, error) {

	token, err := signInToken(signInTokenRequest, sessionDuration)
	if err != nil {
		return "", err
	}

	// Create the get token URL
	url := url.URL{
		Scheme: "https",
		Host:   "signin.aws.amazon.com",
		Path:   "federation",
	}

	q := url.Query()
	q.Set("Action", "login")
	q.Set("Issuer", issuer)
	q.Set("Destination", "https://console.aws.amazon.com/")
	q.Set("SigninToken", token.SigninToken)
	url.RawQuery = q.Encode()

	return url.String(), nil
}

// signInToken makes the AWS call to get a sign-in token from AWS for use when
// creating a sign-in link
func signInToken(fedSession SignInTokenRequest, duration int64) (*SignInTokenResponse, error) {
	fS, err := json.Marshal(fedSession)
	if err != nil {
		return &SignInTokenResponse{}, err
	}

	url := url.URL{
		Scheme: "https",
		Host:   "signin.aws.amazon.com",
		Path:   "federation",
	}

	q := url.Query()
	q.Set("Action", "getSigninToken")
	q.Set("SessionDuration", strconv.FormatInt(duration, 10))
	q.Set("Session", string(fS))
	url.RawQuery = q.Encode()

	fmt.Println(url.String())
	resp, err := http.Get(url.String())
	if err != nil {
		return &SignInTokenResponse{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &SignInTokenResponse{}, err
	}

	signInToken := &SignInTokenResponse{}
	err = json.Unmarshal(body, signInToken)
	if err != nil {
		return &SignInTokenResponse{}, err
	}

	return signInToken, nil
}
