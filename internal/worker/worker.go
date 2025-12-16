package worker

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Mhbib34/missing-person-service/internal/entity"
	"github.com/Mhbib34/missing-person-service/internal/helper"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ResizeImageJobWorker struct {
	// Fields
	db          *gorm.DB
	workerCount int
}

func NewResizeImageJobWorker(db *gorm.DB, workerCount int) *ResizeImageJobWorker {
	// Set default worker count
	if workerCount <= 0 {
		workerCount = 5 // default 5 concurrent workers
	}

	// Create a new ResizeImageJobWorker
	return &ResizeImageJobWorker{
		db:          db,
		workerCount: workerCount,
	}
}

func (w *ResizeImageJobWorker) Start(ctx context.Context, interval time.Duration) {
	log.Printf("🚀 Starting %d resize image workers", w.workerCount)


	// Create a channel to receive jobs
	jobChan := make(chan entity.MissingPersons, w.workerCount)

	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < w.workerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			w.worker(ctx, id, jobChan)
		}(i + 1)
	}


	// Start a ticker to fetch jobs
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			// Stop the worker if the context is done
			case <-ctx.Done():
				close(jobChan)
				return

			// Fetch jobs
			case <-ticker.C:
				jobs, err := w.claimPendingJobs(ctx, w.workerCount)
				if err != nil {
					log.Println("❌ fetch job error:", err)
					continue
				}

				for _, job := range jobs {
					jobChan <- job
				}
			}
		}
	}()

	// Wait for all workers to finish
	wg.Wait()
	log.Println("🛑 Resize image workers stopped")
}

func (w *ResizeImageJobWorker) claimPendingJobs(
	ctx context.Context,
	limit int,
) ([]entity.MissingPersons, error) {

	// Get pending jobs
	var jobs []entity.MissingPersons

	// Use a transaction to ensure that the jobs are claimed atomically
	tx := w.db.WithContext(ctx).Begin()

	// Use a lock to ensure that the jobs are claimed atomically
	err := tx.
		Clauses(clause.Locking{
			Strength: "UPDATE",
			Options:  "SKIP LOCKED",
		}).
		Where("image_status = ?", "processing").
		Limit(limit).
		Find(&jobs).Error

	// Rollback the transaction if an error occurs
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// If no jobs are found, return nil
	if len(jobs) == 0 {
		tx.Rollback()
		return nil, nil
	}

	// Update the image_status for the claimed jobs
	ids := make([]uuid.UUID, 0, len(jobs))
	for _, j := range jobs {
		ids = append(ids, j.ID)
	}

	if err := tx.Model(&entity.MissingPersons{}).
		Where("id IN ?", ids).
		Update("image_status", "processing").Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return jobs, nil
}


func (w *ResizeImageJobWorker) worker(
	ctx context.Context,
	workerID int,
	jobChan <-chan entity.MissingPersons,
) {
	log.Printf("👷 Worker #%d started", workerID)

	// Process jobs
	for {
		// Check if the context is done
		select {
		case <-ctx.Done():
			return
		// Process a job
		case job, ok := <-jobChan:
			if !ok {
				return
			}
			w.processJob(ctx, workerID, job)
		}
	}
}


func (w *ResizeImageJobWorker) processJob( ctx context.Context, workerID int, job entity.MissingPersons, ) {
	log.Printf("🖼️ Worker #%d processing ID %s", workerID, job.ID) 
	// 1️⃣ Bangun path file lokal 
	localPath := filepath.Join("storage/tmp", job.PhotoID) 

	// 2️⃣ Init cloudinary 
	uploader, err := helper.NewCloudinaryUploader() 
	if err != nil { 
		log.Println("❌ cloudinary init error:", err)
		w.updateImageStatus(ctx, job.ID, "failed") 
		return 
	}

	// 3️⃣ Upload + resize 
	cloudURL, err := uploader.UploadResizedImage( 
		ctx, 
		localPath, 
		job.ID.String(), // public_id = ID missing person 
		)
		if err != nil { 
			log.Println("❌ upload error:", err) 
			w.updateImageStatus(ctx, job.ID, "failed") 
			return 
		}

	// 4️⃣ Update DB: photo_id + status 
	err = w.db.WithContext(ctx). 
	Model(&entity.MissingPersons{}). 
	Where("id = ?", job.ID). 
	Updates(map[string]interface{}{ 
		"photo_id": cloudURL, 
		"image_status": "ready", 
		}).Error 
		if err != nil { 
			log.Println("❌ db update error:", err)
			 return 
			}

			// 5️⃣ (OPSIONAL) hapus file lokal 
			_ = os.Remove(localPath) 
				log.Printf("✅ Worker #%d finished job %s", workerID, job.ID)
}


func (w *ResizeImageJobWorker) updateImageStatus(
	ctx context.Context,
	id uuid.UUID,
	status string,
) error {
	// Update the image_status
	return w.db.WithContext(ctx).
		Model(&entity.MissingPersons{}).
		Where("id = ?", id).
		Update("image_status", status).
		Error
}
