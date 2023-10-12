package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gastrader/website/context"
	"github.com/gastrader/website/errors"
	"github.com/gastrader/website/models"
	"github.com/go-chi/chi/v5"
)

type Galleries struct {
	Templates struct {
		New   Template
		Show  Template
		Edit  Template
		Index Template
		Error Template
	}
	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Templates.New.Execute(w, r, data)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID int
		Title  string
	}
	data.UserID = context.User(r.Context()).ID
	data.Title = r.FormValue("title")

	gallery, err := g.GalleryService.Create(data.Title, data.UserID)
	if err != nil {
		g.Templates.New.Execute(w, r, data, err)
		return
	}
	editPath := fmt.Sprintf("/galleries/%s/edit", gallery.UUID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	type Image struct {
		GalleryID       int
		GalleryUUID     string
		Filename        string
		FilenameEscaped string
	}
	var data struct {
		ID        int
		Title     string
		Is_Public bool
		Images    []Image
		UUID string
	}
	data.Is_Public = gallery.Is_Public
	data.ID = gallery.ID
	data.Title = gallery.Title
	data.UUID = gallery.UUID
	images, err := g.GalleryService.Images(gallery.UUID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryUUID:     image.GalleryUUID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}
	g.Templates.Edit.Execute(w, r, data)
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	//does gallery exist?
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	gallery.Title = r.FormValue("title")
	err = g.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/galleries/%s/edit", gallery.UUID), http.StatusFound)
}

func (g Galleries) Index(w http.ResponseWriter, r *http.Request) {
	type Gallery struct {
		ID    int
		UUID  string
		Title string
	}
	var data struct {
		Galleries []Gallery
	}
	user := context.User(r.Context())
	galleries, err := g.GalleryService.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			ID:    gallery.ID,
			UUID: gallery.UUID,
			Title: gallery.Title,
		})
	}
	g.Templates.Index.Execute(w, r, data)
}

func (g Galleries) Show(w http.ResponseWriter, r *http.Request) {
	type Image struct {
		GalleryID       int
		GalleryUUID     string
		Filename        string
		FilenameEscaped string
	}
	var data struct {
		ID        int
		UUID      string
		Title     string
		Images    []Image
		Is_Public bool
		Owner string
	}
	gallery, err := g.galleryByID(w, r, galleryPublic)
	if err != nil {
		if errors.Is(err, models.ErrNotAuthorized) {
			err = errors.Public(err, "You are not authorized to view this resource.")
		}
		g.Templates.Error.Execute(w, r, data, err)
		return
	}
	data.ID = gallery.ID
	data.UUID = gallery.UUID
	data.Title = gallery.Title
	images, err := g.GalleryService.Images(gallery.UUID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryUUID:       image.GalleryUUID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}
	data.Is_Public = gallery.Is_Public
	g.Templates.Show.Execute(w, r, data)
}

func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = g.GalleryService.Delete(gallery.UUID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g Galleries) Image(w http.ResponseWriter, r *http.Request) {
	var data struct{}
	filename := g.filename(w, r)
	galleryUUID := chi.URLParam(r, "id")

	_, err := g.galleryByID(w, r, galleryPublic)
	if err != nil {
		if errors.Is(err, models.ErrNotAuthorized) {
			err = errors.Public(err, "You are not authorized to view this resource.")
		}
		g.Templates.Error.Execute(w, r, data, err)

		return
	}
	image, err := g.GalleryService.Image(galleryUUID, filename)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, image.Path)
}

func (g Galleries) UpdatePerms(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	if r.FormValue("is_public") == "on" {
		gallery.Is_Public = true
	} else {
		gallery.Is_Public = false
	}
	err = g.GalleryService.UpdatePerms(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/galleries/%s/edit", gallery.UUID), http.StatusFound)
}

func (g Galleries) UploadImage(w http.ResponseWriter, r *http.Request) {
	//gallery exists and is owned by user
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = r.ParseMultipartForm(5 << 20) //5mb
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}
	fileHeaders := r.MultipartForm.File["images"]
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		err = g.GalleryService.CreateImage(gallery.UUID, fileHeader.Filename, file)
		if err != nil {
			var fileErr models.FileError
			if errors.As(err, &fileErr) {
				msg := fmt.Sprintf("%v has an invalid content type. Only png, gif and jpg can be uploaded.", fileHeader.Filename)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}

			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
			return
		}
	}
	editPath := fmt.Sprintf("/galleries/%s/edit", gallery.UUID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) DeleteImage(w http.ResponseWriter, r *http.Request) {
	filename := g.filename(w, r)
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = g.GalleryService.DeleteImage(gallery.UUID, filename)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
	editPath := fmt.Sprintf("/galleries/%s/edit", gallery.UUID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) filename(w http.ResponseWriter, r *http.Request) string {
	filename := chi.URLParam(r, "filename")
	filename = filepath.Base(filename)
	return filename
}

type galleryOpt func(http.ResponseWriter, *http.Request, *models.Gallery) error

func (g Galleries) galleryByID(w http.ResponseWriter, r *http.Request, opts ...galleryOpt) (*models.Gallery, error) {
	uuid := chi.URLParam(r, "id")
	gallery, err := g.GalleryService.ByID(uuid)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return nil, err
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return nil, err
	}
	for _, opt := range opts {
		err = opt(w, r, gallery)
		if err != nil {
			return nil, err
		}
	}
	return gallery, nil
}

func userMustOwnGallery(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You are not authorized to edit this gallery", http.StatusForbidden)
		return models.ErrNotAuthorized
	}
	return nil
}

func galleryPublic(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	user := context.User(r.Context())
	if gallery.Is_Public || user != nil && gallery.UserID == user.ID {
		return nil
	}
	// http.Error(w, "You are not authorized to view this gallery", http.StatusForbidden)
	return models.ErrNotAuthorized
}
