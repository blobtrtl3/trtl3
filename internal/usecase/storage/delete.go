package storage

func (bs *BlobStorage) Delete(id string) (bool, error) {
	_, err := bs.db.Exec(
		"DELETE FROM blobsinfo WHERE id=?",
		id,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}
