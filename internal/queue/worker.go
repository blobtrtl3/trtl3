package queue

import "log"

func (bq *BlobQueue) worker() {
	for task := range bq.queue {
		_, err := bq.storage.Save(task.Info, task.Bytes)
		if err != nil {
			log.Println("error while saving blob: ", err)
			bq.Append(task.Info, task.Bytes)
		}

		bq.wg.Done()
	}
}
