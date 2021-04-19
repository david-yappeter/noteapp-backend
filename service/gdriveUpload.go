package service

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "./token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getService() (*drive.Service, error) {
	b, err := ioutil.ReadFile("./credentials.json")
	if err != nil {
		fmt.Printf("Unable to read credentials.json file. Err: %v\n", err)
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)

	if err != nil {
		return nil, err
	}

	client := getClient(config)

	service, err := drive.New(client)

	if err != nil {
		fmt.Printf("Cannot create the Google Drive service: %v\n", err)
		return nil, err
	}

	return service, err
}

func createDir(service *drive.Service, name string, parentID string) (*drive.File, error) {
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentID},
	}

	file, err := service.Files.Create(d).Do()

	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}

	return file, nil
}

func createFile(service *drive.Service, name string, mimeType string, content io.Reader, parentID string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentID},
	}
	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}

func GdriveDeleteFile(fileId string) {
	service, err := getService()
	if err != nil {
		panic(err)
	}

	if err = service.Files.Delete(fileId).Do(); err != nil {
		fmt.Println(err)
	}
}

//UploadFile Upload
func UploadFile(ctx context.Context, userUploadFile graphql.Upload) (string, error) {
	// Step 1. Open the file
	f := userUploadFile.File

	service, err := getService()
	if err != nil {
		panic(err)
	}

	// Step 4. Create the file and upload its content
	file, err := createFile(service, userUploadFile.Filename, userUploadFile.ContentType, f, os.Getenv("GDRIVE_FOLDER_ID"))

	if err != nil {
		panic(fmt.Sprintf("Could not create file: %v\n", err))
	}

	return file.Id, nil
}

//UploadFileBatch Upload Batch
func UploadFileBatch(ctx context.Context, userUploadFile []*graphql.Upload) ([]string, error) {
	var uploadPaths []string

	uploadPathChan := make(chan string)

	for _, val := range userUploadFile {

		go func(val *graphql.Upload) {
			// Step 1. Open the file
			f := val.File

			// Step 2. Get the Google Drive service
			service, err := getService()

			file, err := createFile(service, val.Filename, val.ContentType, f, os.Getenv("GDRIVE_FOLDER_ID"))

			if err != nil {
				panic(fmt.Sprintf("Could not create file: %v\n", err))
			}

			uploadPathChan <- file.Id
		}(val)
	}

	lens := len(userUploadFile)

	for i := 0; i < lens; i++ {
		dataUploaded := <-uploadPathChan
		uploadPaths = append(uploadPaths, dataUploaded)
	}

	return uploadPaths, nil
}

//GdriveViewLink View Link
func GdriveViewLink(fileID *string) *string {
	if fileID == nil {
		return fileID
	}
	temp := fmt.Sprintf("https://drive.google.com/uc?export=view&id=%s", *fileID)
	return &temp
}
