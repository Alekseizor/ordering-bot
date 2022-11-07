package repository

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type docsPhoto struct {
	Server int    `json:"server"`
	Photo  string `json:"photo"`
	Hash   string `json:"hash"`
}

type docsDoc struct {
	File string `json:"file"`
}

func WriteUrl(Db *sqlx.DB, VkID int, attachments []object.MessagesMessageAttachment) {
	var docsURL, docsTitle, imagesURL []string
	var num int
	for _, val := range attachments {
		switch val.Type {
		case "doc":
			docsURL = append(docsURL, val.Doc.URL)
			docsTitle = append(docsTitle, val.Doc.Title)
		case "photo":
			for i, a := range val.Photo.Sizes {
				if a.Type == "z" {
					num = i
					break
				}
				if a.Type == "x" {
					num = i
				}
			}
			imagesURL = append(imagesURL, val.Photo.Sizes[num].URL)
		default:
			break
		}
	}
	ID, err := GetIDOrder(Db, VkID)
	_, err = Db.Exec("INSERT INTO docs (docs_url, docs_title, images_url, order_id) VALUES ($1, $2, $3, $4) ON CONFLICT (order_id) DO UPDATE SET docs_url = $1, docs_title = $2, images_url = $3", pq.Array(docsURL), pq.Array(docsTitle), pq.Array(imagesURL), ID)
	if err != nil {
		log.WithError(err).Error("can`t record docs")
	}
}

func GetDocs(VK *api.VK, Db *sqlx.DB, VkID int) (string, error) {
	var urls, titles []string
	var attachments []string
	ID, err := GetIDOrder(Db, VkID)
	err = Db.QueryRow("SELECT docs_url, docs_title FROM docs WHERE order_id =$1", ID).Scan(pq.Array(&urls), pq.Array(&titles))
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Row with VkID unknown")
		} else {
			log.Println("Couldn't find the line with the docs_url")
		}
		log.Error(err)
	}
	log.Println(urls)
	log.Println(titles)

	for i, val := range urls {
		resp, err := http.Get(val)
		if err != nil {
			log.Fatal(err)
		}
		upload, _ := VK.DocsGetMessagesUploadServer(api.Params{
			"type":    "doc",
			"peer_id": VkID,
		})
		file, err := io.ReadAll(resp.Body)
		fileBody := bytes.NewReader(file)
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", titles[i])
		io.Copy(part, fileBody)
		writer.Close()
		r, _ := http.NewRequest("POST", upload.UploadURL, bytes.NewReader(body.Bytes()))
		r.Header.Set("Content-Type", writer.FormDataContentType())
		client := &http.Client{}
		response, _ := client.Do(r)
		docs := &docsDoc{}
		json.NewDecoder(response.Body).Decode(docs)
		log.Println(docs.File)
		savedDoc, _ := VK.DocsSave(api.Params{
			"file":  docs.File,
			"title": titles[i],
		})
		attachments = append(attachments, "doc"+strconv.Itoa(savedDoc.Doc.OwnerID)+"_"+strconv.Itoa(savedDoc.Doc.ID)+"_"+savedDoc.Doc.AccessKey)
	}
	return strings.Join(attachments[:], ","), nil
}

func GetImages(VK *api.VK, Db *sqlx.DB, VkID int) (string, error) {
	var urls []string
	var attachments []string
	ID, err := GetIDOrder(Db, VkID)
	err = Db.QueryRow("SELECT images_url FROM docs WHERE order_id =$1", ID).Scan(pq.Array(&urls))
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Row with VkID unknown")
		} else {
			log.Println("Couldn't find the line with the docs_url")
		}
		log.Error(err)
	}

	for _, val := range urls {
		resp, err := http.Get(val)
		if err != nil {
			log.Fatal(err)
		}
		upload, _ := VK.PhotosGetMessagesUploadServer(api.Params{
			"peer_id": VkID,
		})

		substr := strings.Split(strings.Split(val, "?")[0], "/")
		photoTitle := substr[len(substr)-1]

		file, err := io.ReadAll(resp.Body)
		fileBody := bytes.NewReader(file)
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", photoTitle)
		io.Copy(part, fileBody)
		writer.Close()

		r, _ := http.NewRequest("POST", upload.UploadURL, bytes.NewReader(body.Bytes()))
		r.Header.Set("Content-Type", writer.FormDataContentType())
		client := &http.Client{}
		response, _ := client.Do(r)
		docs := &docsPhoto{}
		json.NewDecoder(response.Body).Decode(docs)
		savedPhoto, _ := VK.PhotosSaveMessagesPhoto(api.Params{
			"photo":  docs.Photo,
			"server": docs.Server,
			"hash":   docs.Hash,
		})

		attachments = append(attachments, "photo"+strconv.Itoa(savedPhoto[0].OwnerID)+"_"+strconv.Itoa(savedPhoto[0].ID)+"_"+savedPhoto[0].AccessKey)
	}
	return strings.Join(attachments[:], ","), nil
}

func GetAttachments(VK *api.VK, Db *sqlx.DB, VkID int) (string, error) {
	var output string
	//todo: перенести запросы в базу данных сюда из функций get docs и get images +
	// + в них передевать соответственно только нужные поля из базы данных
	return output, nil
}
