package repository

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/Alekseizor/ordering-bot/internal/app/config"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
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

func NewsletterWriteUrl(Db *sqlx.DB, attachments []object.MessagesMessageAttachment) {
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
	_, err := Db.Exec("UPDATE newsletter SET (docs_url, docs_title, images_url) = ($1, $2, $3) WHERE newsletter_id IN(SELECT max(newsletter_id) FROM newsletter)", pq.Array(docsURL), pq.Array(docsTitle), pq.Array(imagesURL))
	if err != nil {
		log.WithError(err).Error("can`t record docs")
	}
}

func NewsletterGetDocs(ctx *context.Context, VK *api.VK, urls, titles []string) (string, error) {
	var attachments []string

	for i, val := range urls {
		resp, err := http.Get(val)
		if err != nil {
			log.Println(err)
		}
		upload, _ := VK.DocsGetMessagesUploadServer(api.Params{
			"type":    "doc",
			"peer_id": config.FromContext(*ctx).AdminID,
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

func NewsLetterGetImages(ctx *context.Context, VK *api.VK, urls []string) (string, error) {

	var attachments []string

	for _, val := range urls {
		resp, err := http.Get(val)
		if err != nil {
			log.Println(err)
		}
		upload, _ := VK.PhotosGetMessagesUploadServer(api.Params{
			"peer_id": config.FromContext(*ctx).AdminID,
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

func NewsletterGetAttachments(ctx *context.Context, VK *api.VK, Db *sqlx.DB) (string, error) {
	var attach ds.Newsletter
	//todo: обработка ошибок
	err := Db.QueryRow("SELECT docs_url, docs_title, images_url, attachment FROM newsletter ORDER BY newsletter_id DESC LIMIT 1").Scan(pq.Array(&attach.DocsUrl), pq.Array(&attach.DocsTitle), pq.Array(&attach.ImagesUrl), &attach.Attachment)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Row with VkID unknown")
		} else {
			log.Println("Couldn't find the line with the docs_url")
		}
		log.Error(err)
	}
	if attach.Attachment != nil {
		log.Println("Used quick attachment")
		return *attach.Attachment, nil
	} else {
		var output, outputDocs, outputImages string
		if attach.DocsUrl != nil {
			outputDocs, _ = NewsletterGetDocs(ctx, VK, attach.DocsUrl, attach.DocsTitle)
			output += outputDocs + ","
		}
		if attach.ImagesUrl != nil {
			outputImages, _ = NewsLetterGetImages(ctx, VK, attach.ImagesUrl)
			output += outputImages
		}
		_, err = Db.Exec("UPDATE newsletter SET attachment = $1 WHERE newsletter_id IN(SELECT max(newsletter_id) FROM newsletter)", output)
		if err != nil {
			log.WithError(err).Error("can`t set attachment in newsletter")
		}
		log.Println("Attachment set (newsletter)")
		return output, nil
	}

}

func SetMessage(Db *sqlx.DB, message string) {
	_, err := Db.Exec("INSERT INTO newsletter (text_message) VALUES ($1)", message)
	if err != nil {
		log.WithError(err).Error("can`t record docs")
	}
}

func GetNewsletterMessage(Db *sqlx.DB) string {
	var output string
	err := Db.QueryRow("SELECT text_message FROM newsletter ORDER BY newsletter_id DESC LIMIT 1").Scan(&output)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Row with VkID unknown")
		} else {
			log.Println("Couldn't find the line with the docs_url")
		}
		log.Error(err)
	}
	return output
}

func SetPeerIDsOrders(Db *sqlx.DB) {
	var newsletter ds.Newsletter
	var peer int
	rows, err := Db.Query("SELECT customer_vk_id FROM orders")
	for rows.Next() {
		if err = rows.Scan(&peer); err != nil {
			log.Println(err)
		}
		newsletter.PeerIDs = append(newsletter.PeerIDs, peer)
	}
	_, err = Db.Exec("UPDATE newsletter SET peer_ids = $1 WHERE newsletter_id IN(SELECT max(newsletter_id) FROM newsletter)", pq.Array(newsletter.PeerIDs))
	if err != nil {
		log.WithError(err).Error("can`t set peer_ids in newsletter")
	}
}

func SetPeerIDsExecutors(Db *sqlx.DB) {
	var newsletter ds.Newsletter
	var peer int
	rows, err := Db.Query("SELECT vk_id FROM executors")
	for rows.Next() {
		if err = rows.Scan(&peer); err != nil {
			log.Println(err)
		}
		newsletter.PeerIDs = append(newsletter.PeerIDs, peer)
	}
	_, err = Db.Exec("UPDATE newsletter SET peer_ids = $1 WHERE newsletter_id IN(SELECT max(newsletter_id) FROM newsletter)", pq.Array(newsletter.PeerIDs))
	if err != nil {
		log.WithError(err).Error("can`t set peer_ids in newsletter")
	}
}

func SetPeerIDsAll(Db *sqlx.DB) {
	var newsletter ds.Newsletter
	var peer int
	rows, err := Db.Query("SELECT vk_id FROM users")
	for rows.Next() {
		if err = rows.Scan(&peer); err != nil {
			log.Println(err)
		}
		newsletter.PeerIDs = append(newsletter.PeerIDs, peer)
	}
	_, err = Db.Exec("UPDATE newsletter SET peer_ids = $1 WHERE newsletter_id IN(SELECT max(newsletter_id) FROM newsletter)", pq.Array(newsletter.PeerIDs))
	if err != nil {
		log.WithError(err).Error("can`t set peer_ids in newsletter")
	}
}

func GetNewsletterPeerIDs(Db *sqlx.DB) []int {
	var peerIDs []sql.NullInt64
	err := Db.QueryRow("SELECT peer_ids FROM newsletter ORDER BY newsletter_id DESC LIMIT 1").Scan(pq.Array(&peerIDs))
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Row with newsletter_id unknown")
		} else {
			log.Println("Couldn't find the line with the docs_url")
		}
		log.Error(err)
	}
	log.Println(peerIDs)
	var IDs []int
	for _, val := range peerIDs {
		IDs = append(IDs, int(val.Int64))
	}
	log.Println(IDs)

	return IDs
}

func ClearNewsletter(Db *sqlx.DB) {
	_, err := Db.Exec("TRUNCATE TABLE newsletter")
	if err != nil {
		log.WithError(err).Error("can`t set attachment in docs")
	}
}

//func ClearAttachments(Db *sqlx.DB, VkID int) error {
//	ID, err := GetIDOrder(Db, VkID)
//	_, err = Db.Exec("UPDATE docs SET attachment = $1 WHERE order_id =$2", nil, ID)
//	if err != nil {
//		log.WithError(err).Error("can`t set attachment in docs")
//	}
//	return err
//}
