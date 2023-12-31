package models

import (
	"database/sql"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gastrader/website/errors"
)

type Image struct {
	Path      string
	GalleryUUID string
	Filename  string
}

type Gallery struct {
	ID        int
	UUID      string
	UserID    int
	Title     string
	Is_Public bool
}

type GalleryService struct {
	DB *sql.DB

	//where to store/locate images
	ImagesDir string
}

func (gs *GalleryService) Create(title string, userID int) (*Gallery, error) {
	gallery := Gallery{
		Title:  title,
		UserID: userID,
	}
	row := gs.DB.QueryRow(`
		INSERT INTO galleries (title, user_id)
		VALUES ($1, $2) RETURNING id;`, gallery.Title, gallery.UserID)
	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	return &gallery, nil
}

func (gs *GalleryService) ByID(uuid string) (*Gallery, error) {
	gallery := Gallery{
		UUID: uuid,
	}
	row := gs.DB.QueryRow(`
		SELECT title, user_id, is_public, id
		FROM galleries
		WHERE uuid = $1;`, gallery.UUID)
	err := row.Scan(&gallery.Title, &gallery.UserID, &gallery.Is_Public, &gallery.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query gallery by id: %w", err)
	}
	return &gallery, nil
}

func (gs *GalleryService) ByUserID(userID int) ([]Gallery, error) {
	rows, err := gs.DB.Query(`
		SELECT id, title, uuid
		FROM galleries
		WHERE user_id = $1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	var galleries []Gallery
	for rows.Next() {
		gallery := Gallery{
			UserID: userID,
		}
		err = rows.Scan(&gallery.ID, &gallery.Title, &gallery.UUID)
		if err != nil {
			return nil, fmt.Errorf("query galleries by user: %w", err)
		}
		galleries = append(galleries, gallery)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	return galleries, nil
}

func (gs *GalleryService) Update(gallery *Gallery) error {
	_, err := gs.DB.Exec(`
		UPDATE galleries
		SET title = $2
		WHERE id = $1;`, gallery.ID, gallery.Title)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	return nil
}

func (gs *GalleryService) UpdatePerms(gallery *Gallery) error {
	_, err := gs.DB.Exec(`
		UPDATE galleries
		SET is_public = $2
		WHERE id = $1;`, gallery.ID, gallery.Is_Public)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	return nil
}

func (gs *GalleryService) Delete(uuid string) error {
	_, err := gs.DB.Exec(`
		DELETE FROM galleries
		WHERE uuid = $1;`, uuid)
	if err != nil {
		return fmt.Errorf("delete gallery.: %w", err)
	}
	err = os.RemoveAll(gs.galleryDir(uuid))
	if err != nil {
		return fmt.Errorf("delete gallery images: %w", err)
	}
	return nil
}

func (gs *GalleryService) galleryDir(uuid string) string {
	imagesDir := gs.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%s", uuid))
}

func (gs *GalleryService) Images(galleryUUID string) ([]Image, error) {
	globPattern := filepath.Join(gs.galleryDir(galleryUUID), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving gallery images: %w", err)
	}
	var images []Image
	for _, file := range allFiles {
		if hasExtension(file, gs.extensions()) {
			images = append(images, Image{
				Path:      file,
				Filename:  filepath.Base(file),
				GalleryUUID: galleryUUID,
			})
		}
	}
	return images, nil
}

func hasExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}

func (gs *GalleryService) Image(galleryUUID string, filename string) (Image, error) {
	imagePath := filepath.Join(gs.galleryDir(galleryUUID), filename)
	_, err := os.Stat(imagePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Image{}, ErrNotFound
		}
		return Image{}, fmt.Errorf("querying for image: %w", err)
	}
	return Image{
		Filename:  filename,
		GalleryUUID: galleryUUID,
		Path:      imagePath,
	}, nil
}

func (gs *GalleryService) extensions() []string {
	return []string{".png", ".jpg", ".jpeg", ".gif"}
}

func (gs *GalleryService) imageContentTypes() []string {
	return []string{"image/png", "image/gif", "image/jpeg", "image/jpg"}
}
func (gs *GalleryService) CreateImage(galleryUUID string, filename string, contents io.ReadSeeker) error {
	err := checkContentType(contents, gs.imageContentTypes())
	if err != nil {
		return fmt.Errorf("creating image: %w", err)
	}
	err = checkExtension(filename, gs.extensions())
	if err != nil {
		return fmt.Errorf("creating image: %w", err)
	}

	galleryDir := gs.galleryDir(galleryUUID)
	err = os.MkdirAll(galleryDir, 0755)
	if err != nil {
		return fmt.Errorf("creating gallery image directory: %w", err)
	}
	imagePath := filepath.Join(galleryDir, filename)
	dst, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("creating image file: %w", err)
	}
	defer dst.Close()
	_, err = io.Copy(dst, contents)
	if err != nil {
		return fmt.Errorf("copying contents to image: %w", err)
	}
	return nil
}

func (gs *GalleryService) DeleteImage(galleryUUID string, filename string) error {
	image, err := gs.Image(galleryUUID, filename)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}
	err = os.Remove(image.Path)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}
	return nil
}
