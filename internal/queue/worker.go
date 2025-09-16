package queue

import "log"

func (bq *BlobQueue) worker() {
	for task := range bq.queue {
		_, err := bq.storage.Save(task.Info, task.Bytes)
		if err != nil {
			task.Retries++
			if task.Retries <= 3 {
				log.Printf("retrying to save blob %s (attempt %d)", task.Info.ID, task.Retries)
				bq.Append(task.Info, task.Bytes)
			} else {
				log.Printf("failed to save blob %s after %d attempts: %s", task.Info.ID, task.Retries, err)
			}
		}

		bq.wg.Done()
	}
}
